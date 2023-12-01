package main

import (
	"CloudScan/pkg/plugin"
	pb "CloudScan/pkg/proto"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"strings"
	"time"
	"bytes"
	"runtime"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

type GWorkspace struct {
	pb.ConnectorServer  `json:"-"`
	pb.PluginServer     `json:"-"`
	pb.NegotiatorServer `json:"-"`
	host                string
	grpcPort            string
	httpPort            string
	//cache               *plugin.CacheEngine
	token    *oauth2.Token
	service  *drive.Service
	key      chan string
	nextPage string
	backOff  int
	cache    *plugin.CacheEngine

	Auth json.RawMessage `json:"auth"`
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func (g *GWorkspace) Negotiate(server pb.Negotiator_NegotiateServer) error {
	//TODO implement me
	request, err := server.Recv()
	if err != nil {
		return err
	}
	googleCredentials, err := json.Marshal(g.Auth)
	if err != nil {
		return err
	}
	fmt.Println("Google Credentials")
	googleCredentials = []byte(strings.ReplaceAll(string(googleCredentials), "{httpPort}", g.httpPort))

	config, err := google.ConfigFromJSON(googleCredentials, drive.DriveMetadataReadonlyScope, drive.DriveReadonlyScope)
	if err != nil {
		return err
	}
	var client *http.Client
	var token *oauth2.Token
	getRemoteCreds := true
	if getRemoteCreds {
		fmt.Println("Requesting Credentials")
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		authReq := pb.NegotiateResponse{
			Seq: request.GetSeq() + 1,
			Message: fmt.Sprintf("Go to the following link in your browser then type the "+
				"authorization code: \n%v\n", authURL),
			Html: fmt.Sprintf("Go to the following link in your browser then type the "+
				"authorization code: \n%v\n", authURL),
			Type: pb.DataType_NULL,
		}
		err = server.Send(&authReq)
		if err != nil {
			return err
		}

		fmt.Println("Waiting for Response")
		_, err := server.Recv()
		if err != nil {
			return err
		}
		fmt.Println("Response received")

		authCode := <-g.key
		fmt.Println("Auth Code: ", authCode)

		tok, err := config.Exchange(context.TODO(), authCode)
		if err != nil {
			log.Fatalf("Unable to retrieve token from web %v", err)
		}
		token = tok
		client = config.Client(context.Background(), token)
	}

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	g.service = srv

	res := pb.NegotiateResponse{
		Seq:     request.Seq + 1,
		Message: "success",
		Html:    "success",
		Type:    pb.DataType_END,
	}
	err = server.Send(&res)

	if err != nil {
		return err
	}
	return nil

}

func (g *GWorkspace) Cat(ctx context.Context, request *pb.ReadFileRequest) (*pb.File, error) {

	path := request.GetPath()
	name := request.GetName()
	mime := request.GetMime()
	fmt.Println("Cat: ", path, name, mime)

	if !g.cache.Contains(path) {
		fmt.Println()
		body, err := getFile(g, path, mime)
		if err != nil {
			fmt.Println("ERROR GETTING FILE IN WORKSPACE/CAT: ", err)
			body := []byte("ERROR GETTING FILE IN WORKSPACE/CAT: " + err.Error())
			s256 := sha256.Sum256(body)
			m5 := md5.Sum(body)

			f := pb.File{
				Url:        fmt.Sprintf("http://localhost:%s/dl?path=%s&name=%s&mime=%s", g.httpPort, url.QueryEscape(request.GetPath()), url.QueryEscape(name), url.QueryEscape(mime)),
				Name:       name,
				ProviderId: "gworkspace",
				MimeType:   mime,
				Size:       uint64(len(body)),
				Sha256:     base64.StdEncoding.EncodeToString(s256[:]),
				Md5:        base64.StdEncoding.EncodeToString(m5[:]),
			}

			return &f, nil
		}

		cached := g.cache.Cache(path, body)
		if !cached {
			fmt.Println("ERROR CACHING FILE IN WORKSPACE/CAT: ", err)
		}
		s256 := sha256.Sum256(body)
		m5 := md5.Sum(body)

		f := pb.File{
			Url:        fmt.Sprintf("http://localhost:%s/dl?path=%s&name=%s&mime=%s", g.httpPort, url.QueryEscape(request.GetPath()), url.QueryEscape(name), url.QueryEscape(mime)),
			Name:       name,
			ProviderId: "gworkspace",
			MimeType:   mime,
			Size:       uint64(len(body)),
			Sha256:     base64.StdEncoding.EncodeToString(s256[:]),
			Md5:        base64.StdEncoding.EncodeToString(m5[:]),
		}

		return &f, nil

	} else {
		body := g.cache.GetCachedValue(path)
		s256 := sha256.Sum256(body)
		m5 := md5.Sum(body)

		f := pb.File{
			Url:        fmt.Sprintf("http://localhost:%s/dl?path=%s&name=%s&mime=%s", g.httpPort, url.QueryEscape(request.GetPath()), url.QueryEscape(name), url.QueryEscape(mime)),
			Name:       name,
			ProviderId: "gworkspace",
			MimeType:   mime,
			Size:       uint64(len(body)),
			Sha256:     base64.StdEncoding.EncodeToString(s256[:]),
			Md5:        base64.StdEncoding.EncodeToString(m5[:]),
		}

		return &f, nil

	}
	//g.cache.Cache(path, body, 1*time.Minute)
	//fmt.Println("Cache: ", path, name, mime, "done")

}

func (g *GWorkspace) Terminate(ctx context.Context, request *pb.TerminateRequest) (*pb.TerminateResponse, error) {
	//TODO implement me
	go awaitTerminate()
	return &pb.TerminateResponse{}, nil
}
func awaitTerminate() {
	time.Sleep(5 * time.Second)
	os.Exit(0)
}

func (g *GWorkspace) Load(ctx context.Context, msg *pb.LoadRequest) (*pb.LoadResponse, error) {
	//fmt.Println("Load")
	params := msg.GetLaunchParams()
	cfg := GWorkspace{}
	err := json.Unmarshal([]byte(params), &cfg)
	if err != nil {
		return &pb.LoadResponse{}, err
	}
	g.Auth = cfg.Auth

	capabilities := []pb.PluginCapability{
		pb.PluginCapability_VIRTUAL_FILESYSTEM,
	}
	res := pb.LoadResponse{Status: 1, Capabilities: capabilities, ShouldNegotiate: true}
	return &res, nil

}

func (g *GWorkspace) Ls(stream pb.Connector_LsServer) error {
	//fmt.Println("GWorkspace LS")

	firstRun := true
	for firstRun || g.nextPage != "" {
		req, err := stream.Recv()
		if err != nil {
			fmt.Println("*****Error in Getting Files: ", err)
			time.Sleep(5 * time.Second)
			break
		}
		fmt.Println("Next Page: ", g.nextPage)
		firstRun = false
		files, err := g.service.Files.List().Corpora("allDrives").SupportsAllDrives(true).IncludeItemsFromAllDrives(true).PageToken(g.nextPage).PageSize(req.GetRequestSize()).Do()
		if err == io.EOF {
			fmt.Println("*****Error in Getting Files (Workspace) EOF: ", err)
			break
		} else if err != nil {
			fmt.Println("*****Error in Getting Files: ", err)
			time.Sleep(5 * time.Second)
			break
		} else if files == nil || files.Files == nil {
			fmt.Println("*****No Files Found, Next Page", g.nextPage)
			break
		}

		g.nextPage = files.NextPageToken
		for i, f := range files.Files {
			mime := f.MimeType
			if f.Size > 20000000 {
				fmt.Println("File too large: ", f.Size, " ", f.Name)
				res := pb.DirectoryEntry{
					Path:      f.Id,
					EntryType: pb.EntryType_NO_INDEX,
					Name:      f.Name,
					Mime:      mime,
					Provider:  "workspace",
					Final:     g.nextPage == "" && i == len(files.Files)-1,
				}
				err := stream.Send(&res)
				if err != nil {
					fmt.Println("Error in Sending File: ", err)
					return err
				}
			} else if mime == "application/vnd.google-apps.folder" {
				res := pb.DirectoryEntry{
					Path:      "",
					EntryType: pb.EntryType_DIRECTORY,
					Name:      "",
					Mime:      mime,
					Provider:  "workspace",
					Final:     g.nextPage == "" && i == len(files.Files)-1,
				}
				err := stream.Send(&res)
				if err != nil {
					fmt.Println("Error in Sending File: ", err)
					return err
				}

			} else {
				res := pb.DirectoryEntry{
					Path:      f.Id,
					EntryType: pb.EntryType_FILE,
					Name:      f.Name,
					Mime:      mime,
					Provider:  "workspace",
					Final:     g.nextPage == "" && i == len(files.Files)-1,
				}
				err := stream.Send(&res)
				if err != nil {
					fmt.Println("Error in Sending File: ", err)
					return err
				}

			}

			if err != nil {
				fmt.Println("Error in Sending File: ", err)
				return err
			}

		}

		fmt.Println("Looping for Next Page: ", g.nextPage)
	}
	fmt.Println("Done Getting Google Drive")

	return nil

}

func getFile(g *GWorkspace, id string, mime string) ([]byte, error) {
	if mime == "application/vnd.google-apps.document" {
		httpRes, err := g.service.Files.Export(id, "application/vnd.openxmlformats-officedocument.wordprocessingml.document").Download()
		if err != nil {
			fmt.Println("Error in Getting File Workspace: ", err)
			return nil, err
		}

		var buf bytes.Buffer
		_, err = io.Copy(&buf, httpRes.Body)
		
		if err != nil{
			return nil, err
		}
		return buf.Bytes(), nil
		
	} else if mime == "application/vnd.google-apps.spreadsheet" {
		httpRes, err := g.service.Files.Export(id, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet").Download()
		if err != nil {
			fmt.Println("Error in Getting File Workspace3: ", err)

			return []byte("ERROR GETTING FILE"), err
		}
		var buf bytes.Buffer
		_, err = io.Copy(&buf, httpRes.Body)
		if err != nil{
			return nil, err
		}
		return buf.Bytes(), nil
	} else if strings.Contains(mime, "application/vnd.google-apps.") {
		fmt.Println("Google Docs Mime Skipping", mime, id)
		return []byte(mime), nil
	} else {
		httpRes, err := g.service.Files.Get(id).Download()
		if err != nil {
			fmt.Printf("Error in Getting File Workspace5: %s, %s, %s", err, id, mime)

			return []byte("ERROR GETTING FILE"), err
		}
		var buf bytes.Buffer
		
		_, err = io.Copy(&buf, httpRes.Body)
		if err != nil{
			return nil, err
		}
		return buf.Bytes(), nil
	}

}

func serveDl(g *GWorkspace) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Dl")
		path := r.URL.Query().Get("path")
		mime, err := url.PathUnescape(r.URL.Query().Get("mime"))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if g.cache.Contains(path) {
			//do something here
			body := g.cache.GetCachedValue(path)
			w.WriteHeader(200)
			w.Write(body)
		} else {
			body, err := getFile(g, path, string(mime))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			w.WriteHeader(200)
			w.Write(body)

		}

	}
}
func serveAuthorize(g *GWorkspace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Println(code)
		io.WriteString(w, fmt.Sprintf("OK. "))
		g.key <- code
	}
}
func serve(g *GWorkspace) {
	r := mux.NewRouter()
	r.HandleFunc("/dl", serveDl(g))
	r.HandleFunc("/authorize", serveAuthorize(g))
	r.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	r.Handle("/debug/pprof/{cmd}", http.HandlerFunc(pprof.Index))
	http.ListenAndServe(fmt.Sprintf(":%s", g.httpPort), r)
}

func main() {
	port := os.Args[1]
	httpPort := os.Args[2]
	server := GWorkspace{backOff: 0}

	//server.cache = plugin.NewCache()
	server.httpPort = httpPort
	server.grpcPort = port
	server.key = make(chan string)
	c := plugin.NewCache(5 * time.Second)
	server.cache = c

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterPluginServer(grpcServer, &server)
	pb.RegisterConnectorServer(grpcServer, &server)
	pb.RegisterNegotiatorServer(grpcServer, &server)
	//fmt.Println("Listening on port", port)
	go serve(&server)
	go func(){
		for{
			time.Sleep(30*time.Second)
			runtime.GC()
		}
	
	}()

	err = grpcServer.Serve(lis)
	if err != nil {
		fmt.Println("##############Error in Serving: ", err)
		os.WriteFile("error.txt", []byte(err.Error()), 0644)
	}
}
