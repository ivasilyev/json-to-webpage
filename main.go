package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

var serverPort uint64
var filePath string

func templateRenderer(w http.ResponseWriter, m map[string]string) {
	var s = ""
	for k, v := range m {
		if strings.HasPrefix(v, "\"http") {
			s += fmt.Sprintf("<tr><td>%s</td><td><a href=%s>%s</a></td></tr>\n", k, v, v)
		} else {
			s += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", k, v)
		}
	}

	tmpl, err := template.ParseFiles("./html/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func jsonDataHandler(w http.ResponseWriter, r *http.Request) {
	// Read the JSON file from disk
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read JSON file: '%s'", filePath), http.StatusInternalServerError)
		return
	}

	// Parse the JSON data
	jsonRawData := make(map[string]json.RawMessage)
	jsonStringData := make(map[string]string)

	err = json.Unmarshal(fileData, &jsonRawData)
	if err != nil {
		panic(err)
	}

	for k, v := range jsonRawData {
		b, _ := json.Marshal(v)
		var s = string(b)
		jsonStringData[k] = s
	}

	// Generate HTML table
	templateRenderer(w, jsonStringData)
}

func main() {
	var j = flag.String("j", "./test/data.json", "JSON file")
	var p = flag.Uint64("p", 12000, "Server port")
	flag.Parse()
	filePath = *j
	serverPort = *p

	http.HandleFunc("/", jsonDataHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))

	log.Println(fmt.Sprintf("Server started on http://localhost:%d", serverPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(": %d", serverPort), nil))
}
