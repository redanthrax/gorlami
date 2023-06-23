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
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	log.Println("Listening on port 9001.")
	log.Fatal(http.ListenAndServe(":9001", nil))
}

func handleView(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	files := []string{
		"frontend/templates/base.html",
	}

	switch r.URL.Path {
	case "/":
		files = append(files, "frontend/login.html")
	case "/agents.html":
		files = append(files,
			"frontend/templates/dash_base.html",
			"frontend/agents.html")
	case "/settings.html":
		files = append(files,
			"frontend/templates/dash_base.html",
			"frontend/settings.html")
	default:
		http.NotFound(w, r)
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("Error parsing template files: %s", err)
		http.Error(w, "Server error.", http.StatusInternalServerError)
	}

	tmpl.ExecuteTemplate(w, "base", nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	user := r.FormValue("username")
	pass := r.FormValue("password")
	if user == "admin" && pass == "admin" {
		//handle session stuff
		w.Header().Add("HX-Redirect", "/agents.html")
	} else {
		//send denied fragment
		tmpl, err := template.ParseFiles("frontend/templates/fragments/warning.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		err = tmpl.ExecuteTemplate(w, "warning", "Access Denied")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	//clear session data
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
