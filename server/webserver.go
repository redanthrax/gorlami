package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func startWebServer() {
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)

	//handle sessions
	go cleanupSessions()

	//start server
	log.Println("Listening on :3000...")
	err := http.ListenAndServe("127.0.0.1:3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/favicon.ico")
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("web/templates", "layout.html")
	fp := ""
	if r.URL.Path == "/" {
		fp = "web/templates/index.html"
	} else {
		session := getSession(w, r)
		log.Printf("%#v\n", session)
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		fp = filepath.Join("web/templates", filepath.Clean(r.URL.Path))
	}

	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) || info.IsDir() {
			http.NotFound(w, r)
			return
		}
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	//mock login
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "admin" && password == "admin" {
		log.Println("Login action performed...")
		session := createSession()
		session.Data["username"] = username
		saveSession(session, w)
		w.Header().Set("HX-Redirect", "/agents.html")
	} else {
		w.Header().Set("HX-Redirect", "/")
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session := getSession(w, r)
	if session != nil {
		deleteSession(session.ID)
		clearSessionCookie(w)
	}

	w.Header().Set("HX-Redirect", "/")
}
