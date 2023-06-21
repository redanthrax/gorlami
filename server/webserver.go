package main

import (
	"html/template"
	"log"
	"net/http"
)

func startWebserver() {
	//static files
	static := http.StripPrefix("/static", http.FileServer(http.Dir("frontend/static")))
	http.Handle("/static/", static)
	//setup root view
	http.HandleFunc("/", handleView)
	log.Println("Listening on port 3000.")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleView(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"frontend/templates/base.html",
		"frontend/login.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("Error parsing template files: %s", err)
		http.Error(w, "Server error.", http.StatusInternalServerError)
	}

	tmpl.ExecuteTemplate(w, "base", nil)
}
