package main

import (
	pb "CloudScan/pkg/proto"
	"encoding/xml"
	"fmt"
	"github.com/beevik/etree"
	"strconv"
)

type Sst struct {
	XMLName     xml.Name `xml:"sst"`
	Text        string   `xml:",chardata"`
	Xmlns       string   `xml:"xmlns,attr"`
	Count       string   `xml:"count,attr"`
	UniqueCount string   `xml:"uniqueCount,attr"`
	Si          []Si     `xml:"si"`
}
type Si struct {
	Text string `xml:",chardata"`
	T    string `xml:"t"`
}

func ParseXLSX(f FileMap) *pb.TransformEntry {
	t := pb.TransformEntry{
		Type:        pb.TransformEntryType_GROUP,
		Uid:         0,
		Contents:    "",
		Children:    make([]*pb.TransformEntry, 0),
		Correlation: 0,
	}
	sharedStrings := f["xl/sharedStrings.xml"]
	//fmt.Println("Opened Shared Strings")
	sst := make(map[string]string)
	sstXML := etree.NewDocument()
	if err := sstXML.ReadFromBytes(sharedStrings); err != nil {
		//fmt.Println("Error Reading Shared Strings XML", err.Error())

		return nil
	}
	fmt.Println("Read Shared Strings XML", string(sharedStrings))
	root := sstXML.SelectElement("sst")
	if root != nil {

		for i, entry := range root.FindElements("./[t]") {
			fmt.Println("Adding to sst ", entry.SelectElement("t").Text())
			sst[strconv.Itoa(i)] = entry.SelectElement("t").Text()
		}
		/*for i, entry := range root.FindElements("./[v]") {
			fmt.Println("Adding to sst ", entry.SelectElement("v").Text())
			sst[strconv.Itoa(i)] = entry.SelectElement("v").Text()
		}*/

	}

	//fmt.Println("Processing Workbook")
	fmt.Printf("SST: %+v\n", sst)
	rawWorkbook := f["xl/workbook.xml"]
	workbook := etree.NewDocument()
	if err := workbook.ReadFromBytes(rawWorkbook); err != nil {
		//fmt.Println(err)
		return nil
	}
	//fmt.Println("Opened Workbook")
	worksheets := make(map[string]string)
	wbRoot := workbook.SelectElement("workbook")
	if wbRoot != nil {

		sheets := wbRoot.SelectElement("sheets")
		for i, sheet := range sheets.SelectElements("sheet") {
			worksheets[sheet.SelectAttrValue("sheetId", strconv.Itoa(i))] = sheet.SelectAttrValue("name", strconv.Itoa(i))
			//fmt.Println("Adding to worksheets ", sheet.SelectAttrValue("name", strconv.Itoa(i)))
		}

		for i, name := range worksheets {
			//fmt.Println("Processing Worksheet: ", name)
			tWorkSheet := pb.TransformEntry{
				Type:        pb.TransformEntryType_GROUP,
				Contents:    name,
				Children:    make([]*pb.TransformEntry, 0),
				Correlation: 0.1,
			}
			sheet := f[fmt.Sprintf("xl/worksheets/sheet%s.xml", i)]
			if sheet != nil {
				s := etree.NewDocument()
				if err := s.ReadFromBytes(sheet); err != nil {
					//fmt.Println("Error reading worksheet")
					//fmt.Println(err)
					return nil
				}
				//fmt.Println("Read from worksheet")
				worksheet := s.SelectElement("worksheet")
				//fmt.Println("Selected worksheet")
				if worksheet != nil {
					data := worksheet.SelectElement("sheetData")

					//fmt.Println("Selected sheet data")
					if data != nil {
						for i, row := range data.ChildElements() {
							fmt.Println("Processing row: ", i)
							tRow := pb.TransformEntry{
								Type:        pb.TransformEntryType_GROUP,
								Contents:    strconv.Itoa(i),
								Children:    make([]*pb.TransformEntry, 0),
								Correlation: 1,
							}
							for _, col := range row.SelectElements("c") {
								if col != nil {
									val := col.SelectElement("v")
									if val != nil {
										txt := val.Text()
										if value, ok := sst[txt]; ok {
											fmt.Println("Found value in sst: ", value)
											txt = value
											//do something here
										}
										tVal := pb.TransformEntry{
											Type:        pb.TransformEntryType_STRING,
											Contents:    txt,
											Children:    nil,
											Correlation: 1,
										}
										tRow.Children = append(tRow.Children, &tVal)
									}
								}

							}
							tWorkSheet.Children = append(tWorkSheet.Children, &tRow)

						}
						t.Children = append(t.Children, &tWorkSheet)
					}
				}

			}
		}
	}

	//fmt.Println("Sheet assembled")

	return &t
}
