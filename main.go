package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const PORT = 12000

func templateRenderer(w http.ResponseWriter, m map[string]string) {
	funcMap := template.FuncMap{
		"isUrl": func(value string) bool {
			fmt.Println(value, strings.HasPrefix(value, "\"http"))
			return strings.HasPrefix(value, "\"http")
		},
	}

	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("./html/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, m)
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
		/*
			if strings.HasPrefix(s, "\"http") {
				s = fmt.Sprintf("<a href=\"%s\">%s</a>", s, s)
				fmt.Println("%s, %s", s)
			}
		*/
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
