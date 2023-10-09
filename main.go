package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

const PORT = 12000

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
	filePath := "./data.json"
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read JSON file", http.StatusInternalServerError)
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
	http.HandleFunc("/", jsonDataHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))

	log.Println(fmt.Sprintf("Server started on http://localhost:%d", PORT))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(": %d", PORT), nil))
}
