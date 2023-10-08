package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

const PORT = 12000

func templateRenderer(w http.ResponseWriter, m map[string]string) {
	tmpl, err := template.ParseFiles("./html/index.html")
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
		o, _ := json.Marshal(v)
		jsonStringData[k] = string(o)
	}

	for k, v := range jsonStringData {
		fmt.Println("%s, %s", k, v)
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
