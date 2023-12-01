package reporting

import (
	_ "embed"
	"html/template"
	"os"
)

//go:embed report.gohtml
var reportTemplate []byte

func GenerateHTMLReport(report Report) {
	tmpl, err := template.New("report").Parse(string(reportTemplate))
	if err != nil {
		panic(err)
	}
	var f *os.File
	f, err = os.Create("report.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.ExecuteTemplate(f, "report", report)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}

}
