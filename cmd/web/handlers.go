package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	// Initialize a slice containing the paths of the two files. Note that the
	// home.page.gohtml file must be the first file in the slice.
	files := []string{
		"./ui/html/home.page.gohtml",
		"./ui/html/base.layout.gohtml",
		"./ui/html/footer.partial.gohtml",
	}
	
	// Use the template.ParseFiles() function to read the files and store the templates in a
	// template set. Note the ... here works very much like the spread operator in JavaScript as it
	// unpacks the contents of the slice. If there is an error, log the detailed error message
	// and use the http.Error() function to send a generic 500 Internal Server Error response to the user.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	
	// Then use the Execute() method from the html/template library on the template set to write the
	// template content as the response body. The last parameter to Execute() represents any dynamic
	// data that is passed in, which is currently nil.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func navError (w http.ResponseWriter, r *http.Request) {
_, err := w.Write([]byte("There's no page here, Madison."))
if err != nil {
log.Fatal("There was a problem with the NavError route", err)
}
}

func showSnippet (w http.ResponseWriter, r *http.Request) {
id, err := strconv.Atoi(r.URL.Query().Get("id"))
if err != nil || id < 1 {
http.NotFound(w, r)
return
}
_, e := fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
if e != nil {
log.Fatal("There was a problem with the ShowSnippet route", e)
}
}

func createSnippet (w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
// If the method is not POST use the w.WriteHeader() method to send a 405
// status code and the w.Write() method to write a "Method Not Allowed"
// response body. We then return from the function so that the
// subsequent code is not executed
// Use the Header().Set() method to provide additional information to the user
// by setting a header as Allow:POST
// Need to ensure headers are set BEFORE WriteHeader() or Write() are called
// Remember the difference between WriteHeader() and Header().Set() is that
// WriteHeader() only sets a status code and that can't be changed once set or set again
// where Header().Set() sends the headers in the standard key:value format
w.Header().Set("Allow", http.MethodPost)
//w.WriteHeader(405)
//_, err := w.Write([]byte("Method not allowed \n"))
//if err != nil {
//	return
//}
// Can combine the WriteHeader() and Write() into using the http.Error() method
// which takes the ResponseWriter, an error message string, and the http code to be returned
// We can avoid having to do error handling on the Write() method using this method
http.Error(w, "Method not allowed", 405)
return
}
w.Header().Set("Content-Type", "application/json")
w.Header().Add("Cache-Control", "public")
_, err := w.Write([]byte(`{"name":"Madison"}`))
if err != nil {
http.Error(w, "There was a problem with the CreateSnippet route", 400)
}
}