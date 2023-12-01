package main

import (
	pb "CloudScan/pkg/proto"
	"archive/zip"
	"encoding/xml"
	"io"
	"strconv"
)

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`
}

func BuildFileTree(z *zip.Reader) FileMap {
	fileMap := make(FileMap)

	for _, file := range z.File {
		f, err := file.Open()
		if err != nil {
			panic(err)
		}
		contents, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		fileMap[file.Name] = contents
	}
	return fileMap

}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node Node

	return d.DecodeElement((*node)(n), &start)
}

func traverse(n Node, f func(Node) (bool, string)) (bool, pb.TransformEntry) {
	shouldIndex, contents := f(n)
	e := pb.TransformEntry{
		Correlation: 0.1,
		Contents:    "",
	}
	if shouldIndex {
		e.Type = pb.TransformEntryType_STRING
		e.Contents = contents
	}
	var children []*pb.TransformEntry
	for _, node := range n.Nodes {
		e.Type = pb.TransformEntryType_GROUP
		childShouldIndex, child := traverse(node, f)
		if childShouldIndex {
			children = append(children, &child)
		}
	}
	e.Children = children
	if len(children) > 1 || shouldIndex {
		return true, e
	}
	if len(children) == 1 {
		c := *children[0]
		c.Correlation += e.Correlation
		return true, c
	}
	return false, pb.TransformEntry{}
}
func walk(nodes []Node, path string, f func(string, Node) (bool, pb.TransformEntry)) {
	for i, n := range nodes {
		newPath := path + "/" + strconv.Itoa(i)
		found, _ := f(newPath, n)
		if found {
			walk(n.Nodes, newPath, f)
		}
	}
}
