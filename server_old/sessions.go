package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	ID        string
	ExpiresAt time.Time
	Data      map[string]interface{}
}

var (
	sessionMutex sync.Mutex
	sessions     = make(map[string]Session)
)

const (
	sessionCookieName      = "gorlami_id"
	sessionDuration        = 24 * time.Hour
	sessionIDLength        = 32
	sessionCleanupInterval = 1 * time.Hour
)

func createSession() Session {
	sessionID := generateSessionID()
	expiresAt := time.Now().Add(sessionDuration)
	session := Session{
		ID:        sessionID,
		ExpiresAt: expiresAt,
		Data:      make(map[string]interface{}),
	}

	log.Println("Session created")
	return session
}

func generateSessionID() string {
	b := make([]byte, sessionIDLength)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}

func getSession(r *http.Request) *Session {
	log.Println("Getting sessions")
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil
	}

	sessionID := cookie.Value
	session, ok := sessions[sessionID]
	if !ok || session.ExpiresAt.Before(time.Now()) {
		if ok {
			deleteSession(sessionID)
		}

		return nil
	}

	session.ExpiresAt = time.Now().Add(sessionDuration)
	sessions[sessionID] = session
	return &session
}

func saveSession(session Session, w http.ResponseWriter) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	session.ExpiresAt = time.Now().Add(sessionDuration)
	sessions[session.ID] = session

	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    session.ID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
	log.Println("Session saved")
}

func deleteSession(sessionID string) {
	log.Println("Deleting sessions")
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	delete(sessions, sessionID)
}

func clearSessionCookie(w http.ResponseWriter) {
	log.Println("Clearing session cookies")
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
}

func cleanupSessions() {
	for {
		time.Sleep(sessionCleanupInterval)
		sessionMutex.Lock()
		for id, session := range sessions {
			if session.ExpiresAt.Before(time.Now()) {
				delete(sessions, id)
			}
		}

		sessionMutex.Unlock()
	}
}
