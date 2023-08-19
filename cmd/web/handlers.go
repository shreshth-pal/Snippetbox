package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)



func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	files := []string {
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}
	ts,err := template.ParseFiles(files...)
	if err != nil {
		fmt.Fprint(w,"Internal Server Error",http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w,"base",nil)
	if err != nil {
		log.Print(err,"error")
		http.Error(w,"Internal Server Error",http.StatusInternalServerError)
		return
	}

	
}



func snippetView(w http.ResponseWriter, r *http.Request) {
id, err := strconv.Atoi(r.URL.Query().Get("id"))
if err != nil || id < 1 {
http.NotFound(w, r)
return
}
fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}




func snippetCreate(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
w.Header().Set("Allow", http.MethodPost)
http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
return
}
w.Write([]byte("Create a new snippet..."))
}