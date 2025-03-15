package cookies

import (
	"GoStreamRecord/internal/db"
	dbapi "GoStreamRecord/internal/db/api"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var UserStore map[string]string

var Session *session

type session struct {
	subs_mutex *sync.Mutex
	cookies    *sessions.CookieStore
	apiKeys    []string // preloaded valid API keys
}

func New() *session {
	var key *sessions.CookieStore
	if db.Config.Settings.App.UseEnvKey {
		// TODO: Implement check if .env is available. fallback to random key if not.
		key = sessions.NewCookieStore([]byte("KKp5q43GBoaSgd6QmQyehpSqW/08KRhIJpO90URviFs="))
	} else {
		key = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

	}
	session := session{
		subs_mutex: &sync.Mutex{},
		cookies:    key,
	}
	session.cookies.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to false for HTTP
		MaxAge:   3600,
	}
	return &session
}

func (s *session) Store() *sessions.CookieStore {
	return s.cookies
}

// IsLoggedIn checks if the user is logged in either by session or via a valid API key.
func (s *session) IsLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	// Retrieve the session
	session, err := s.cookies.Get(r, "session")
	if err != nil {
		http.RedirectHandler("/login", http.StatusInternalServerError)
		return false
	}

	// If session is authenticated, return true immediately.
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		return true
	}

	// Attempt to get API key from header (fallback to query param if needed)
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		apiKey = r.URL.Query().Get("api_key")
	}

	// If a valid API key is provided, consider the user authenticated.
	if apiKey != "" && s.isValidAPIKey(apiKey) {
		// Optionally, mark the session as authenticated for subsequent requests.
		session.Values["authenticated"] = true
		if err := s.cookies.Save(r, w, session); err != nil {
			http.Error(w, "Session save error", http.StatusInternalServerError)
			return false
		}
		return true
	}

	// Not authenticated by either method.
	return false
}

// isValidAPIKey compares the provided API key against the preloaded valid keys.
func (s *session) isValidAPIKey(providedKey string) bool {
	if len(s.apiKeys) == 0 {

		err := db.Read("api", "api.json", &db.Config.APIKeys)
		if err != nil {
			log.Println("Error getting existing keys..", err)
			return false
		}

		for _, k := range db.Config.APIKeys.Keys {

			exist := false
			for _, existingKey := range s.apiKeys {
				if existingKey == k.Key {
					exist = true
					break
				}

			}
			if !exist {
				s.apiKeys = append(s.apiKeys, k.Key)
			}
		}
	}
	for _, key := range s.apiKeys {
		if dbapi.VerifyAPIKey(key, providedKey) {
			return true
		}
	}
	return false
}
