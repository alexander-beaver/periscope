package main

import (
	"CloudScan/pkg/plugin"
	pb "CloudScan/pkg/proto"
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gorilla/mux"
	azure "github.com/microsoft/kiota-authentication-azure-go"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Graph365 struct {
	pb.ConnectorServer  `json:"-"`
	pb.NegotiatorServer `json:"-"`
	pb.PluginServer     `json:"-"`
	host                string
	grpcPort            string
	httpPort            string
	cache               *plugin.CacheEngine

	DeviceCodeCredential *azidentity.DeviceCodeCredential `json:"-"`
	graphUserScopes      []string

	ClientID string                                     `json:"client"`
	TenantID string                                     `json:"tenant"`
	Secret   string                                     `json:"secret"`
	Scopes   []string                                   `json:"-"`
	Client   *azure.AzureIdentityAuthenticationProvider `json:"-"`
}

func (g *Graph365) Negotiate(server pb.Negotiator_NegotiateServer) error {
	//TODO implement me
	allowedHosts := []string{"graph.microsoft.com"}
	//fmt.Println("Negotiating with client")
	var request *pb.NegotiateRequest
	request, err := server.Recv()
	//fmt.Println(request)
	if request != nil {
		//fmt.Println("Got request")
	} else if err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	//fmt.Println("Creating device code credential")

	creds, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		ClientID: g.ClientID,
		TenantID: g.TenantID,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			//fmt.Println("Loading device code prompt")

			//fmt.Println(message)
			prompt := pb.NegotiateResponse{
				Seq:     request.Seq + 1,
				Message: message.Message,
				Html:    message.Message,
				Type:    pb.DataType_NULL,
			}
			err := server.Send(&prompt)
			if err == nil {
				return err
			}

			_, err = server.Recv()
			if err != nil {
				return err
			}
			//fmt.Println(res.Seq)

			return nil
		},
	})
	if err != nil {
		//fmt.Println(err)
		return err
	}
	//fmt.Println("Created device code credential")

	authProvider, err := azure.NewAzureIdentityAuthenticationProviderWithScopesAndValidHosts(creds, g.Scopes, allowedHosts)
	if err != nil {
		fmt.Printf("Error creating auth provider: %v\n", err)
	}

	g.Client = authProvider
	meUrl, err := url.ParseRequestURI("https://graph.microsoft.com/v1.0/me")

	token, err := authProvider.GetAuthorizationTokenProvider().GetAuthorizationToken(context.Background(), meUrl, nil)
	if err != nil {
		return err
	}

	req := http.Request{
		Method: "GET",
		URL:    meUrl,
		Header: map[string][]string{"Authorization": {fmt.Sprintf("Bearer %s", token)}},
	}
	client := &http.Client{}

	me, err := client.Do(&req)
	if err != nil {
		return err
	}
	_, err = io.ReadAll(me.Body)
	if err != nil {
		return err
	}
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

func (g *Graph365) getFile(client http.Client, path string) ([]byte, error) {
	reqUrl := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/items/%s/content", path)
	reqAsUrl, err := url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}
	token, err := g.Client.GetAuthorizationTokenProvider().GetAuthorizationToken(context.Background(), reqAsUrl, nil)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	req.Header.Add("Accept", "application/json")
	fmt.Println("Sending Command")
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
func (g *Graph365) Cat(ctx context.Context, request *pb.ReadFileRequest) (*pb.File, error) {

	path := request.GetPath()
	name := request.GetName()
	fmt.Println("Getting file " + path)
	client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// Print the URL of the redirected request
		// Allow a maximum of 5 redirects
		if len(via) >= 5 {
			return errors.New("too many redirects")
		}
		return nil
	}}

	body, err := g.getFile(client, path)
	if err != nil {
		return nil, err
	}

	s256 := sha256.Sum256(body)

	m5 := md5.Sum(body)
	//g.cache[path] = body
	//fmt.Println("Cached file")
	//go g.clearFromCache(path)
	g.cache.Cache(path, body)
	f := pb.File{
		Url:        fmt.Sprintf("http://localhost:%s/dl?path=%s&name=%s", g.httpPort, request.GetPath(), base64.URLEncoding.EncodeToString([]byte(name))),
		Name:       name,
		ProviderId: "graph365",
		MimeType:   plugin.DetermineMimeType(name, body),
		Size:       uint64(len(body)),
		Sha256:     base64.StdEncoding.EncodeToString(s256[:]),
		Md5:        base64.StdEncoding.EncodeToString(m5[:]),
	}

	return &f, nil
}

func (g *Graph365) Terminate(ctx context.Context, request *pb.TerminateRequest) (*pb.TerminateResponse, error) {
	//TODO implement me
	go awaitTerminate()
	return &pb.TerminateResponse{}, nil
}
func awaitTerminate() {
	time.Sleep(5 * time.Second)
	os.Exit(0)
}

func (g *Graph365) Load(ctx context.Context, msg *pb.LoadRequest) (*pb.LoadResponse, error) {
	//fmt.Println("Load")
	params := msg.GetLaunchParams()
	cfg := Graph365{}
	err := json.Unmarshal([]byte(params), &cfg)
	if err != nil {
		return &pb.LoadResponse{}, err
	}
	g.ClientID = cfg.ClientID
	g.TenantID = cfg.TenantID
	g.Secret = cfg.Secret
	g.Scopes = []string{"User.Read", "User.Read.All", "Sites.Read.All", "Files.Read.All"}

	capabilities := []pb.PluginCapability{
		pb.PluginCapability_VIRTUAL_FILESYSTEM,
	}
	res := pb.LoadResponse{Status: 1, Capabilities: capabilities, ShouldNegotiate: true}
	return &res, nil

}

func (g *Graph365) Ls(stream pb.Connector_LsServer) error {
	//fmt.Println("Graph365 LS")

	reqUrl := "https://graph.microsoft.com/v1.0/search/query"
	reqAsUrl, err := url.Parse(reqUrl)
	if err != nil {
		return err
	}
	//fmt.Println("Request URL Built")

	i := 0
	max := 9999

	client := &http.Client{}
	//fmt.Println("HTTP Client Created")

	for i < max {
		request, err := stream.Recv()
		if err != nil {
			return err
		}
		//fmt.Println(i)

		query := DriveSearchRequest{
			Requests: []DriveSearchRequestItem{
				{
					EntityTypes: []string{"driveItem"},
					Query: struct {
						QueryString string `json:"queryString"`
					}{QueryString: "*"},
					From: i,
					Size: int(request.GetRequestSize()),
				},
			},
		}
		i += int(request.GetRequestSize())

		if err != nil {
			return err
		}
		//fmt.Println("Getting Authorization Token")
		token, err := g.Client.GetAuthorizationTokenProvider().GetAuthorizationToken(context.Background(), reqAsUrl, nil)
		//fmt.Println("Got Auth Token")
		if err != nil {
			return err
		}
		queryAsJson, err := json.Marshal(query)
		if err != nil {
			return err
		}
		queryReader := bytes.NewReader(queryAsJson)

		req, err := http.NewRequest("POST", reqUrl, queryReader)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
		//fmt.Println("Making Request")
		response, err := client.Do(req)
		if err != nil {
			//fmt.Println("error recevied")
			return err
		}
		var res DriveSearchResponse
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			//fmt.Println(err.Error())
			return err
		}

		err = json.Unmarshal(body, &res)
		if err != nil {
			return err
		}

		for _, entry := range res.Value {
			for _, hitC := range entry.HitsContainers {
				max = hitC.Total

				for j, h := range hitC.Hits {
					i += 1
					fmt.Printf("Total: %d, i: %d\n", max, i)

					r := h.Resource
					res := pb.DirectoryEntry{
						Path:      r.ID,
						EntryType: pb.EntryType_FILE,
						Name:      r.Name,
						Provider:  "graph365",
						Final:     i >= max-1 && j == len(hitC.Hits)-1,
					}
					err = stream.Send(&res)
					if err != nil {
						//fmt.Println(err.Error())
					}

				}
			}
		}

	}
	fmt.Println("Finished Getting Files From Graph365")
	return nil

}

func serveDl(g *Graph365) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		fmt.Println("Getting file " + path)
		if g.cache.Contains(path) {
			//do something here
			body := g.cache.GetCachedValue(path)
			w.Write(body)
			return
		}

		client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Print the URL of the redirected request
			// Allow a maximum of 5 redirects
			if len(via) >= 5 {
				return errors.New("too many redirects")
			}
			return nil
		}}

		body, err := g.getFile(client, path)
		if err != nil {
			w.WriteHeader(400)
			io.WriteString(w, err.Error())
			return
		}

		fmt.Println("Read file")
		fmt.Print("!")
		fmt.Print("Body Length: " + strconv.Itoa(len(body)))

		w.WriteHeader(200)
		w.Write(body)

	}
}
func serve(g *Graph365) {
	r := mux.NewRouter()
	r.HandleFunc("/dl", serveDl(g))

	http.ListenAndServe(fmt.Sprintf(":%s", g.httpPort), r)
}

func main() {
	port := os.Args[1]
	httpPort := os.Args[2]
	server := Graph365{}
	c := plugin.NewCache(10 * time.Second)
	server.cache = c
	server.httpPort = httpPort
	server.grpcPort = port

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

	err = grpcServer.Serve(lis)
	if err != nil {
		return
	}
}
