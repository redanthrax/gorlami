package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/websocket"
)

func startWebServer() {
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.Handle("/ws", websocket.Handler(handleWebSocket))

	//handle sessions
	go cleanupSessions()

	//initialize agents
	agents = []Agent{}

	//mock agents
	go mockAddRemoveAgents()

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
		//make sure we have a good session
		session := getSession(r)
		if session == nil {
			clearSessionCookie(w)
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

	handleTemplateData(w, tmpl, r.URL.Path)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	//mock login
	username := r.FormValue("username")
	password := r.FormValue("password")
	//normally you'd do db stuff here, not right now
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
	session := getSession(r)
	if session != nil {
		deleteSession(session.ID)
		clearSessionCookie(w)
	}

	w.Header().Set("HX-Redirect", "/")
}

func handleTemplateData(w http.ResponseWriter,
	tmpl *template.Template, path string) {
	data := make(map[string]interface{})
	//setup switch case for data here
	switch path {
	case "/agents.html":
		log.Println("Godamnit")
		log.Printf("%#v\n", agents)
		data["Agents"] = agents
	default:
		data = nil
	}

	log.Printf("%#v\n", data)

	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
