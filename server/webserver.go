package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/pion/webrtc/v3"
	"golang.org/x/net/websocket"
)

func startWebserver() {
	//static files
	static := http.StripPrefix("/static", http.FileServer(http.Dir("frontend/static")))
	http.Handle("/static/", static)
	//ws
	http.Handle("/ws", websocket.Handler(handleWebsocket))

	//setup root view
	http.HandleFunc("/", handleView)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/agent/register", handleAgentRegister)
	log.Println("Listening on port 9001.")
	log.Fatal(http.ListenAndServe(":9001", nil))
}

func handleView(w http.ResponseWriter, r *http.Request) {
	tmplData := make(map[string]interface{})
	files := []string{
		"frontend/templates/base.html",
	}

	log.Println(r.URL.Path)
	switch r.URL.Path {
	case "/":
		files = append(files, "frontend/login.html")
	case "/agents.html":
		files = append(files,
			"frontend/templates/dash_base.html",
			"frontend/agents.html")
		tmplData["Agents"] = agents
	case "/connect.html":
		files = append(files,
			"frontend/templates/dash_base.html",
			"frontend/connect.html")
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

	tmpl.ExecuteTemplate(w, "base", tmplData)
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

func handleAgentRegister(w http.ResponseWriter, r *http.Request) {
	//handle agent registration
}

func handleWebsocket(ws *websocket.Conn) {
	log.Println("Websockets connection started")
	log.Printf("%#v\n", ws)
	//this will move to using nats and get setup with the client
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})

	if err != nil {
		panic(err)
	}

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		if err := websocket.JSON.Send(ws, c.ToJSON()); err != nil {
			panic(err)
		}
	})

	for {
		msg := &webrtc.SessionDescription{}
		if err := websocket.JSON.Receive(ws, msg); err != nil {
			return
		}

		if err := peerConnection.SetRemoteDescription(*msg); err != nil {
			panic(err)
		}

		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			panic(err)
		}

		if err := peerConnection.SetLocalDescription(answer); err != nil {
			panic(err)
		}

		if err := websocket.JSON.Send(ws, answer); err != nil {
			panic(err)
		}
	}
}
