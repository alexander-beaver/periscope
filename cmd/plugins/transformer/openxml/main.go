package main

import (
	"CloudScan/pkg/plugin"
	pb "CloudScan/pkg/proto"
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
	"time"
)

type FileMap map[string][]byte
type OpenXML struct {
	pb.TransformerServer
	pb.PluginServer
	cache    *plugin.CacheEngine
	grpcPort string
	httpPort string
	backOff  int
}

func (o OpenXML) Transform(ctx context.Context, request *pb.TransformRequest) (*pb.TransformResponse, error) {
	//TODO implement me

	req, err := http.Get(request.File.GetUrl())
	if err != nil {
		return nil, err
	}
	combined, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return nil, err
	}
	fmt.Println("Transforming file", request.File.GetName())
	fmt.Println("Request File Length: ", len(combined))
	//Unzip rawBytes

	//Loop through files in zip
	if request.File.GetMimeType() == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" || request.File.GetMimeType() == "application/vnd.google-apps.document" {

		reader, err := zip.NewReader(bytes.NewReader(combined), int64(len(combined)))
		if err != nil {
			return nil, err
		}
		fmt.Println("DOCX")
		tree := BuildFileTree(reader)
		fmt.Println("PARSING DOCX")

		doc := ParseDOCX(tree)
		fmt.Println("CACHING DOCX")
		fmt.Printf("DOCX: %s", doc)
		cacheUrl := cache(&o, request.File, doc)
		fmt.Printf("RETURNING URL: %s \n", cacheUrl)

		res := pb.TransformResponse{Url: cacheUrl, File: request.File}
		return &res, nil
	} else if request.File.GetMimeType() == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" || request.File.GetMimeType() == "application/vnd.google-apps.spreadsheet" {
		fmt.Println("SPREADSHEET")
		reader, err := zip.NewReader(bytes.NewReader(combined), int64(len(combined)))
		if err != nil {
			fmt.Printf("ERROR READING XLSX: %s", err.Error())
			return nil, err
		}
		tree := BuildFileTree(reader)
		fmt.Println("PARSING XLSX")
		xlsx := ParseXLSX(tree)
		fmt.Println("CACHING XLSX")
		cacheUrl := cache(&o, request.File, xlsx)
		fmt.Printf("RETURNING URL: %s \n", cacheUrl)
		res := pb.TransformResponse{File: request.File, Url: cacheUrl}
		return &res, nil
	}

	//Encode contents
	entry := pb.TransformEntry{
		Type:        0,
		Uid:         0,
		Contents:    "Office",
		Children:    nil,
		Correlation: 0,
	}

	cacheUrl := cache(&o, request.File, &entry)
	res := pb.TransformResponse{File: request.File, Url: cacheUrl}
	return &res, nil
}

func cache(o *OpenXML, file *pb.File, entry *pb.TransformEntry) string {
	m := jsonpb.Marshaler{
		OrigName:     false,
		EnumsAsInts:  false,
		EmitDefaults: false,
		Indent:       "",
		AnyResolver:  nil,
	}
	marshalBuffer := new(bytes.Buffer)
	err := m.Marshal(marshalBuffer, entry)
	marshaled := marshalBuffer.Bytes()
	fmt.Printf("Caching File: %s, entry %s \n", file.GetName(), string(marshaled))
	if err != nil {
		fmt.Printf("Error Cacheing: %s", err.Error())
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	id := strconv.Itoa(rand.Int())
	res := o.cache.Cache(id, marshaled)
	if res == false {
		fmt.Printf("Error with cache in OpenXML")
		return ""
	}

	cacheUrl := fmt.Sprintf("http://localhost:%s/dl?path=%s", o.httpPort, id)
	return cacheUrl
}
func (o OpenXML) Load(ctx context.Context, request *pb.LoadRequest) (*pb.LoadResponse, error) {
	//TODO implement me
	//handlers := []string{"application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.openxmlformats-officedocument.presentationml.presentation", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}
	handlers := []string{"application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.google-apps.document", "application/vnd.google-apps.spreadsheet"}

	capabilities := []pb.PluginCapability{
		pb.PluginCapability_TRANSFORMER,
	}
	res := pb.LoadResponse{Status: 1, Capabilities: capabilities, Handlers: handlers}
	return &res, nil

}

func (o OpenXML) Terminate(ctx context.Context, request *pb.TerminateRequest) (*pb.TerminateResponse, error) {
	//TODO implement me
	go awaitTerminate()
	return &pb.TerminateResponse{}, nil
}
func awaitTerminate() {
	time.Sleep(5 * time.Second)
	os.Exit(0)
}

func serveDl(g *OpenXML) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")

		if g.cache.Contains(path) {
			//do something here
			body := g.cache.GetCachedValue(path)
			fmt.Println("Serving Cached File")
			_, err := w.Write(body)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			return
		}

		fmt.Println("Serving Unknown File")
		fmt.Println("Failed to find file in cache", path)
		w.WriteHeader(200)
		res := pb.TransformEntry{Contents: path, Type: pb.TransformEntryType_NOINDEX}
		m := jsonpb.Marshaler{}
		buffer := new(bytes.Buffer)

		err := m.Marshal(buffer, &res)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		encode := buffer.Bytes()
		w.Write(encode)

	}
}
func serve(g *OpenXML) {
	r := mux.NewRouter()
	r.HandleFunc("/dl", serveDl(g))
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
	server := OpenXML{
		grpcPort: port,
		httpPort: httpPort,
		cache:    plugin.NewCache(30 * time.Second),
		backOff:  100,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPluginServer(grpcServer, &server)
	pb.RegisterTransformerServer(grpcServer, &server)

	fmt.Println("Listening on port", port)
	go serve(&server)
	err = grpcServer.Serve(lis)
	if err != nil {
		return
	}
}
