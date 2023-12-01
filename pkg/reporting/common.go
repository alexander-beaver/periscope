package reporting

import pb "CloudScan/pkg/proto"

type Report struct {
	Timestamp int64
	Title     string
	Results   []pb.AnalyzeResponse
}
