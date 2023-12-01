package main

import (
	pb "CloudScan/pkg/proto"
	"encoding/xml"
	"fmt"
)

func ParseDOCX(fileMap FileMap) *pb.TransformEntry {
	//TODO implement me
	contents := fileMap["word/document.xml"]
	var n Node
	err := xml.Unmarshal(contents, &n)
	if err != nil {
		fmt.Println("Error unmarshalling xml", err.Error())
		panic(err)
	}
	_, entry := traverse(n, func(n Node) (bool, string) {
		if n.XMLName.Local == "t" {
			return true, string(n.Content)
		}
		return false, ""
	})

	return &entry

}
