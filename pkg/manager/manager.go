package manager

import (
	pb "CloudScan/pkg/proto"
	"CloudScan/pkg/reporting"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type PlugInType string
type PlugInMap map[string]ActivePlugin
type PlugInHandlerMap map[string]ActivePlugin

const (
	PlugInTypeLocal PlugInType = "local"
)

type PlugIn struct {
	Name     string
	Type     PlugInType
	Command  string
	Args     []string
	Parallel bool
	Nodes    int
	Config   json.RawMessage
}

type Config struct {
	PlugIns []PlugIn
}

type ActivePlugin struct {
	Name         string
	Capabilities []pb.PluginCapability
	URLs         []string
	Handlers     []string
	Nodes        int
	Conn         []*grpc.ClientConn
}

func StartPlugin(plugin PlugIn) ActivePlugin {
	return StartLocalPlugin(plugin)

}
func StartLocalPlugin(plugin PlugIn) ActivePlugin {
	var conns []*grpc.ClientConn
	var capabilities []pb.PluginCapability
	var urls []string
	var handlers []string
	for i := 0; i < plugin.Nodes; i++ {

		port := 1024 + rand.Intn(65535-1024)
		httpPort := 1024 + rand.Intn(65535-1024)

		url := "localhost:" + strconv.Itoa(port)
		status := make(chan int)
		flag := rand.Int()
		go RunPlugin(plugin, port, httpPort, status, flag)
		for (<-status) != flag {
		}
		var opts []grpc.DialOption
		callOps := []grpc_retry.CallOption{grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond))}

		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDisableHealthCheck(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                60 * time.Minute,
			Timeout:             60 * time.Minute,
			PermitWithoutStream: false,
		}), grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoff.Config{
			BaseDelay:  1 * time.Second,
			Multiplier: 1.6,
			Jitter:     1,
			MaxDelay:   30 * time.Second,
		}}),
			grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(callOps...)),
			grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(callOps...)))
		conn, err := grpc.Dial(url, opts...)

		if err != nil {
			panic(err)
		}
		pluginClient := pb.NewPluginClient(conn)

		cfg, err := plugin.Config.MarshalJSON()
		if err != nil {
			panic(err)
		}
		loadMsg, err := pluginClient.Load(context.Background(), &pb.LoadRequest{LaunchParams: string(cfg)})

		if err != nil {
			panic(err)
		}
		//fmt.Println("Loaded plugin", plugin.Name, "with status", loadMsg.GetStatus())
		//fmt.Println(loadMsg.GetShouldNegotiate())
		capabilities = loadMsg.GetCapabilities()
		handlers = loadMsg.GetHandlers()
		urls = append(urls, url)
		if loadMsg.GetShouldNegotiate() {
			negConnect := pb.NewNegotiatorClient(conn)
			negClient, err := negConnect.Negotiate(context.Background())
			if err != nil {
				panic(err)
			}

			initial := pb.NegotiateRequest{
				Seq:   1,
				Type:  pb.DataType_NULL,
				Input: "HELLO",
			}
			err = negClient.Send(&initial)
			if err != nil {
				return ActivePlugin{}
			}

			run := true
			for run == true {
				//fmt.Println("Waiting for response")
				res, err := negClient.Recv()
				//fmt.Println("Response received")
				//fmt.Println(res)

				if err == io.EOF {
					//fmt.Println("EOF")
					run = false

				} else if err != nil {
					//fmt.Println("Error in StartLocalPlugin", err)
					/*if strings.Contains(err.Error(), "EOF") {
						run = false
					}*/
					/*if strings.Contains(err.Error(), "read: connection reset by peer") {
						//fmt.Println("Connection reset by peer")
						os.Exit(98)
					}*/
				} else {
					//fmt.Println("Got response", res.GetSeq(), res.GetType(), res.GetMessage())
					if res.GetType() == pb.DataType_END {
						run = false
					} else if res.GetType() == pb.DataType_STR {
						fmt.Println(res.GetMessage())

						reader := bufio.NewReader(os.Stdin)
						line, err := reader.ReadString('\n')
						if err != nil {
							log.Fatal(err)
						}
						query := pb.NegotiateRequest{
							Seq:   res.GetSeq() + 1,
							Type:  pb.DataType_STR,
							Input: line,
						}
						negClient.Send(&query)
					} else if res.GetType() == pb.DataType_NULL {
						fmt.Println(res.GetMessage())
						query := pb.NegotiateRequest{
							Seq:   res.GetSeq() + 1,
							Type:  pb.DataType_NULL,
							Input: "OK",
						}
						err = negClient.Send(&query)
						if err != nil {
							return ActivePlugin{}
						}
					}
				}

			}
		}
		conns = append(conns, conn)
	}
	//fmt.Println(loadMsg.GetCapabilities())
	fmt.Println("Loaded plugin",plugin.Name,"with",plugin.Nodes,"nodes")

	ap := ActivePlugin{
		Name:         plugin.Name,
		Capabilities: capabilities,
		URLs:         urls,
		Handlers:     handlers,
		Conn:         conns,
		Nodes:        plugin.Nodes,
	}

	return ap

}
func RunPlugin(plugin PlugIn, port int, httpPort int, status chan int, flag int) {
	//fmt.Println("Starting Plugin")
	var args []string
	if plugin.Args != nil {
		args = plugin.Args
	} else {
		args = []string{}
	}
	args = append(args, []string{strconv.Itoa(port), strconv.Itoa(httpPort)}...)

	//fmt.Println("Starting plugin", plugin.Name, "with args", args)
	cmd := exec.Command(plugin.Command, args...)
	//cmd.Stdout = os.Stdout

	err := cmd.Start()
	//fmt.Println("Command started")

	if err != nil {
		panic(err)
	}

	// Look for "Listening on port" in stdout
	// If found, set status to 2
	time.Sleep(1000 * time.Millisecond)
	//fmt.Println("Returning Flag")
	status <- flag
	defer cmd.Wait()

}

func GetFileFromPlugin(pluginConn *grpc.ClientConn, conn pb.ConnectorClient, output chan pb.File, f *pb.DirectoryEntry) {
	file, err := conn.Cat(context.Background(), &pb.ReadFileRequest{Path: f.GetPath(), Name: f.GetName(), Mime: f.GetMime()}, grpc_retry.WithMax(5)) //TODO Change to 10
	if err != nil {
		//fmt.Println("Error in GET FILE", err)
		if strings.Contains(err.Error(), "connection refused") {
			pluginConn.Connect()

		}

	}
	if err == nil && file != nil {
		output <- *file
	}

}

func GetFileFromPluginWorker(pluginConn *grpc.ClientConn, conn pb.ConnectorClient, output chan pb.File, input chan pb.DirectoryEntry) {
	for f := range input {
		GetFileFromPlugin(pluginConn, conn, output, &f)
	}
}
func GetFilesFromPlugin(plugin ActivePlugin, output chan pb.File, plugInHandlerMap PlugInHandlerMap) {
	//fmt.Println("Getting Files")
	conn := pb.NewConnectorClient(plugin.Conn[0])
	files := make(chan pb.DirectoryEntry)

	go func() {
		stream, err := conn.Ls(context.Background())

		if err != nil {
			//fmt.Println("Error in GetFilesFromPlugin", err)

		} else {
			run := true
			for run {
				reqSize := int64(100)
				req := pb.ListChildrenRequest{
					RequestSize: reqSize,
				}
				stream.Send(&req)
				for i := 0; int64(i) < reqSize; i++ {
					file, err := stream.Recv()
					if err == io.EOF {
						//fmt.Println("EOF ERROR")
						run = false
						break
					} else if err != nil {
						if strings.Contains(err.Error(), "EOF") {
							//fmt.Println("ERROR CONTAINS EOF IN GETFILESFROMPLUGIN")
							time.Sleep(10 * time.Second)
						} else {
							//fmt.Println("Error in GetFilesFromPlugin", err)
							plugin.Conn[0].Connect()
						}
						continue
					}
					files <- *file
				}
				//time.Sleep(100 * time.Millisecond)
			}
			close(files)
		}

	}()

	run := true
	var wg sync.WaitGroup
	for i := int64(0); i < 40; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for run {

				file, ok := <-files

				if !ok{
					fmt.Println("Receive channel closed")
					run = false
					break;
				}

				if file.GetEntryType() == pb.EntryType_FILE {
					if _, allowed := plugInHandlerMap[file.Mime]; allowed {
						fmt.Println("Getting file", file.GetName())
						GetFileFromPlugin(plugin.Conn[0], conn, output, &file)
					}else{
						fmt.Println("MIME",file.Mime,"cannot be processed")
					}
				}
			}
		}()
	}
	wg.Wait()
	//fmt.Println("Waiting for files to finish")
	//fmt.Printf("Finished getting files from %s\n", plugin.Name)
}

func Transformer(plugin ActivePlugin, input *pb.File, output chan pb.TransformResponse) {
	randId := rand.Intn(plugin.Nodes)
	fmt.Println("Transforming file", input.GetName(), "with plugin", plugin.Name, "on node", randId, "of", plugin.Nodes)
	conn := pb.NewTransformerClient(plugin.Conn[randId])
	f := pb.TransformRequest{File: input}

	transformed, err := conn.Transform(context.Background(), &f)
	if err != nil {
		//fmt.Println(err)
		return
	}

	output <- *transformed

	//fmt.Println("Transformed file")

}

func GetAllFiles(pluginMap PlugInMap, plugins []string, output chan pb.File, plugInHandlerMap PlugInHandlerMap) {
	var wg sync.WaitGroup
	for _, plugin := range plugins {
		wg.Add(1)
		plugin := plugin
		go func() {
			defer wg.Done()
			GetFilesFromPlugin(pluginMap[plugin], output, plugInHandlerMap)
		}()
	}
	wg.Wait()
	close(output)
	//fmt.Println("Finished getting files")
}

func Analyzer(p ActivePlugin, t pb.TransformResponse, c chan pb.AnalyzeResponse) {
	randId := rand.Intn(p.Nodes)
	conn := pb.NewAnalyzerClient(p.Conn[randId])
	transformRequest := pb.AnalyzeRequest{File: t.File, Url: t.Url}
	res, err := conn.Analyze(context.Background(), &transformRequest)
	if err != nil {
		//fmt.Println(err)
		return
	}
	//fmt.Println(res)
	c <- *res
}

func AnalyzerRunner(pluginMap PlugInMap, analyzers []string, transformations chan pb.TransformResponse, c chan pb.AnalyzeResponse) {
	//fmt.Println("Starting analyzer runner")
	var wg sync.WaitGroup
	for {
		transformation, ok := <-transformations
		//fmt.Println("Got transformation", transformation.File.GetName())

		if ok {

			for _, plugin := range analyzers {
				p := pluginMap[plugin]
				wg.Add(1)
				//fmt.Printf("Starting analyzer %s for file %s\n", p.Name, transformation.File.GetName())
				go func() {
					defer wg.Done()
					Analyzer(p, transformation, c)
					fmt.Println("Analyzed file", transformation.File.GetName())
				}()

			}
		} else {
			break
		}
	}
	wg.Wait()
	//fmt.Println("Finished analyzing")
	close(c)
}

func TransformerRunner(pluginMap PlugInMap, transformers []string, plugInHandlerMap PlugInHandlerMap, files chan pb.File, transformations chan pb.TransformResponse) {
	var wg sync.WaitGroup
	for {
		file, ok := <-files
		if ok {

			mime := file.MimeType

			if p, a := plugInHandlerMap[mime]; a {

				wg.Add(1)
				go func() {
					defer wg.Done()
					Transformer(p, &file, transformations)
					//fmt.Println("Transformed file", file.GetName())
				}()
			} else {
				//fmt.Println("No transformer for", mime)
			}

		} else {
			break
		}
	}
	wg.Wait()
	//fmt.Println("Finished transforming")
	close(transformations)
}
func TerminateAllConnections(pluginMap PlugInMap) {
	for _, plugin := range pluginMap {
		for _, pconn := range plugin.Conn {
			conn := pb.NewPluginClient(pconn)
			conn.Terminate(context.Background(), &pb.TerminateRequest{})
			pconn.Close()
		}
	}
}

func Run(config Config) {
	pluginMap := make(PlugInMap)
	pluginHandlerMap := make(PlugInHandlerMap)
	pluginTypes := make(map[pb.PluginCapability][]string)
	//fmt.Println("Starting plugins")
	for _, plugin := range config.PlugIns {
		fmt.Println("Loading plugin", plugin.Name)

		fmt.Println("Starting plugin", plugin.Name)
		ap := StartPlugin(plugin)
		pluginMap[plugin.Name] = ap
		for _, capability := range ap.Capabilities {
			fmt.Println("Plugin", plugin.Name, "has capability", capability)
			pluginTypes[capability] = append(pluginTypes[capability], plugin.Name)
			if capability == pb.PluginCapability_TRANSFORMER {
				for _, handler := range ap.Handlers {
					pluginHandlerMap[handler] = ap
				}
			}
		}

	}
	//fmt.Println("Started plugins")
	//fmt.Println(pluginTypes)

	fileResponse := make(chan pb.File)
	transformations := make(chan pb.TransformResponse)
	analysis := make(chan pb.AnalyzeResponse)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		GetAllFiles(pluginMap, pluginTypes[pb.PluginCapability_VIRTUAL_FILESYSTEM], fileResponse, pluginHandlerMap)
	}()
	go func() {
		defer wg.Done()
		TransformerRunner(pluginMap, pluginTypes[pb.PluginCapability_TRANSFORMER], pluginHandlerMap, fileResponse, transformations)
	}()
	go func() {
		defer wg.Done()
		AnalyzerRunner(pluginMap, pluginTypes[pb.PluginCapability_ANALYZER], transformations, analysis)
	}()
	// Populate analyses with the channel analysis

	response := make(chan []pb.AnalyzeResponse)
	go func(response chan []pb.AnalyzeResponse) {
		analyses := make([]pb.AnalyzeResponse, 0)

		for {
			a, ok := <-analysis
			//fmt.Println("Got analysis", a)
			if ok {
				//fmt.Println("Appending analysis", a)
				if len(a.GetFindings()) > 0{
					analyses = append(analyses, a)
				}
			} else {
				response <- analyses
				break
			}
		}
	}(response)
	wg.Wait()

	analyses := <-response
	//fmt.Println("Finished")
	TerminateAllConnections(pluginMap)
	//fmt.Println("All connections closed")
	report := reporting.Report{
		Timestamp: 0,
		Title:     "File Results",
		Results:   analyses,
	}
	//fmt.Println("Generating report")

	//fmt.Println(string(reportJson))
	reporting.GenerateHTMLReport(report)
	//fmt.Println("Finished generating report")

}
