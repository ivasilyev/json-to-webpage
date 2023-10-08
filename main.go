package json_to_webpage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Data struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	http.HandleFunc("/", jsonDataHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func jsonDataHandler(w http.ResponseWriter, r *http.Request) {
	// Read the JSON file from disk
	filePath := "data.json"
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read JSON file", http.StatusInternalServerError)
		return
	}

	// Parse the JSON data
	var jsonData []Data
	err = json.Unmarshal(fileData, &jsonData)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusInternalServerError)
		return
	}

	// Generate HTML table
	html := `
 <html>
 <head>
  <link rel="stylesheet" type="text/css" href="/static/style.css">
  <script src="/static/script.js"></script>
 </head>
 <body>
  <table>
   <tr>
    <th>Name</th>
    <th>Value</th>
   </tr>
 `
	for _, data := range jsonData {
		html += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", data.Name, data.Value)
	}
	html += `
  </table>
 </body>
 </html>
 `

	// Write the HTML response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}
