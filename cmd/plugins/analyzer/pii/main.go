package main

import (
	pb "CloudScan/pkg/proto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type PiiScan struct {
	pb.AnalyzerServer
	pb.PluginServer
	httpPort string
	grpcPort string
}

func (l *PiiScan) Load(ctx context.Context, request *pb.LoadRequest) (*pb.LoadResponse, error) {
	//TODO implement me
	capabilities := []pb.PluginCapability{
		pb.PluginCapability_ANALYZER,
	}
	return &pb.LoadResponse{
		Status:          1,
		Capabilities:    capabilities,
		ShouldNegotiate: false,
	}, nil

}

func (l *PiiScan) Terminate(ctx context.Context, request *pb.TerminateRequest) (*pb.TerminateResponse, error) {
	//TODO implement me
	go awaitTerminate()
	return &pb.TerminateResponse{}, nil
}
func awaitTerminate() {
	time.Sleep(5 * time.Second)
	os.Exit(0)
}
func GetTransformEntryFromURL(url string) (*pb.TransformEntry, error) {
	//TODO implement me
	var t pb.TransformEntry

	req, err := http.Get(url)
	if err != nil {

		return nil, err
	}
	combined, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	err = jsonpb.Unmarshal(bytes.NewReader(combined), &t)
	if err != nil {
		fmt.Println("Error unmarshalling transform entry")
		fmt.Println(combined)
		return nil, err
	}
	return &t, nil
}

func (l *PiiScan) Analyze(ctx context.Context, r *pb.AnalyzeRequest) (*pb.AnalyzeResponse, error) {
	rUrl := r.Url
	transformEntry, err := GetTransformEntryFromURL(rUrl)
	fmt.Println("Analyzing: ", rUrl, "With body", transformEntry)

	if err != nil {
		return nil, err
	}
	c, err := json.Marshal(transformEntry)
	contents := string(c)
	if err != nil {
		return &pb.AnalyzeResponse{}, err
	}
	//TODO implement me
	ssnRegex := regexp.MustCompile(`\d{3}-\d{2}-\d{4}`)

	findings := []*pb.AnalyzeFinding{}
	ssn := ssnRegex.FindAllString(contents, -1)
	score := int32(0)
	if len(ssn) > 0 {
		for _, s := range ssn {
			score += 100
			findings = append(findings, &pb.AnalyzeFinding{
				Score:       100,
				Location:    "",
				Contents:    s,
				Description: "PII-SSN",
			})
		}
	}

	visaRegex := regexp.MustCompile(`4[0-9]{12}(?:[0-9]{3})?`)

	cards := visaRegex.FindAllString(contents, -1)
	if len(cards) > 0 {
		for _, c := range cards {
			if Luhns(c) {
				fmt.Println("Found a visa", c)
				score += 100
				findings = append(findings, &pb.AnalyzeFinding{
					Score:       100,
					Location:    "",
					Contents:    c,
					Description: "PCI-Visa",
				})
			}
		}
	}

	mcRegex := regexp.MustCompile(`(?:5[1-5][0-9]{2}|222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}`)
	mcCards := mcRegex.FindAllString(contents, -1)
	if len(mcCards) > 0 {
		for _, c := range mcCards {
			if Luhns(c) {
				fmt.Println("Found a mastercard", c)
				score += 100
				findings = append(findings, &pb.AnalyzeFinding{
					Score:       100,
					Location:    "",
					Contents:    c,
					Description: "PCI-Mastercard",
				})
			}
		}
	}
	amexRegex := regexp.MustCompile(`3[47][0-9]{13}`)
	amexCards := amexRegex.FindAllString(contents, -1)
	if len(amexCards) > 0 {
		for _, c := range amexCards {
			if Luhns(c) {
				fmt.Println("Found a amex", c)
				score += 100
				findings = append(findings, &pb.AnalyzeFinding{
					Score:       100,
					Location:    "",
					Contents:    c,
					Description: "PCI-Amex",
				})
			}
		}
	}
	dinersRegex := regexp.MustCompile(`3(?:0[0-5]|[68][0-9])[0-9]{11}`)
	dinersCards := dinersRegex.FindAllString(contents, -1)
	if len(dinersCards) > 0 {
		for _, c := range dinersCards {
			if Luhns(c) {
				fmt.Println("Found a diners", c)
				score += 100
				findings = append(findings, &pb.AnalyzeFinding{
					Score:       100,
					Location:    "",
					Contents:    c,
					Description: "PCI-Diners",
				})
			}
		}
	}
	jcbRegex := regexp.MustCompile(`(?:2131|1800|35\d{3})\d{11}`)

	jcbCards := jcbRegex.FindAllString(contents, -1)
	if len(jcbCards) > 0 {
		for _, c := range jcbCards {
			if Luhns(c) {
				fmt.Println("Found a jcb", c)
				score += 100
				findings = append(findings, &pb.AnalyzeFinding{
					Score:       100,
					Location:    "",
					Contents:    c,
					Description: "PCI-JCB",
				})
			}
		}
	}

	unionRegex := regexp.MustCompile(`(?:62[0-9]{14,17})`)
	unionCards := unionRegex.FindAllString(contents, -1)
	if len(unionCards) > 0 {
		for _, c := range unionCards {
			if Luhns(c) {
				fmt.Println("Found a unionpay", c)
				score += 100
				findings = append(findings, &pb.AnalyzeFinding{
					Score:       100,
					Location:    "",
					Contents:    c,
					Description: "PCI-UnionPay",
				})
			}
		}
	}
	eloRegex := regexp.MustCompile(`(?:636[0-9]{13,16})`)
	eloCards := eloRegex.FindAllString(contents, -1)
	if len(eloCards) > 0 {
		for _, c := range eloCards {
			if Luhns(c) {
				fmt.Println("Found a elo", c)
				score += 100
				findings = append(findings, &pb.AnalyzeFinding{
					Score:       100,
					Location:    "",
					Contents:    c,
					Description: "PCI-Elo",
				})
			}
		}
	}

	discoverRegex := regexp.MustCompile(`(?:6011|65\d{2}|64[4-9]\d)\d{12}|(?:62\d{14,17})`)
	discoverCards := discoverRegex.FindAllString(contents, -1)
	if len(discoverCards) > 0 {
		for _, c := range discoverCards {
			if Luhns(c) {
				fmt.Println("Found a discover", c)
				score += 100
				findings = append(findings, &pb.AnalyzeFinding{
					Score:       100,
					Location:    "",
					Contents:    c,
					Description: "PCI-Discover",
				})
			}
		}
	}

	triggers := []string{"CVV", "CVC", "CVV2", "CID", "CSC", "CCV", "CVC2", "CVVC", "CVD", "CVN", "SSN", "Social Security"}
	for _, t := range triggers {
		if strings.Contains(contents, t) {
			fmt.Println("Found a PII identifier", t)
			score += 80
			findings = append(findings, &pb.AnalyzeFinding{
				Score:       80,
				Location:    "",
				Contents:    t,
				Description: "PII-Potential",
			})
		}
	}
	if len(findings) > 0 {
		fmt.Println("Found PII risks")
	} else {
		fmt.Println("No PII risks found")
	}
	res := pb.AnalyzeResponse{
		Findings: findings,
		Score:    score,
		File:     r.GetFile(),
	}
	return &res, nil
}

type ConfigParams struct {
	AddIns []string `json:"enable"`
}

func Luhns(card string) bool {
	var sum int
	var alternate bool
	for i := len(card) - 1; i >= 0; i-- {
		digit := int(card[i] - '0')
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}
	return sum%10 == 0
}

func main() {
	port := os.Args[1]
	httpPort := os.Args[2]
	server := PiiScan{}
	server.httpPort = httpPort
	server.grpcPort = port

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterPluginServer(grpcServer, &server)
	pb.RegisterAnalyzerServer(grpcServer, &server)
	//fmt.Println("Listening on port", port)

	err = grpcServer.Serve(lis)
	if err != nil {
		return
	}
}
