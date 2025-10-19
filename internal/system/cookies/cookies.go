package cookies

import (
	"log"
	"net/http"
	"remoteCtrl/internal/db/jsondb"
	"remoteCtrl/internal/system/settings"
	"sync"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var UserStore map[string]string

var Session *session

type session struct {
	subs_mutex *sync.Mutex
	cookies    *sessions.CookieStore
	apiKeys    []string
}

func New(s settings.Settings) *session { 
	session := session{
		subs_mutex: &sync.Mutex{},
		cookies:    sessions.NewCookieStore(securecookie.GenerateRandomKey(32)),
	}
	session.cookies.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   3600,
	}
	return &session
}

func (s *session) Store() *sessions.CookieStore {
	return s.cookies
}

// IsLoggedIn checks if the user is logged in either by session or via a valid API key.
func (s *session) IsLoggedIn(apiKeys settings.API_secrets, w http.ResponseWriter, r *http.Request) bool {

	session, err := s.cookies.Get(r, "session")
	if err != nil {
		http.RedirectHandler("/login", http.StatusInternalServerError)
		//http.Error(w, "Session error. Try clearing your cookies.", http.StatusInternalServerError)
		return false
	}

	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		return true
	}

	// Attempt to get API key from header
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		apiKey = r.URL.Query().Get("api_key")
	}

	// If a valid API key is provided, consider the user authenticated.
	if apiKey != "" && s.isValidAPIKey(apiKeys, apiKey) {
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
func (s *session) isValidAPIKey(apiKeys settings.API_secrets, providedKey string) bool {
	if len(s.apiKeys) == 0 {

		err := jsondb.Load(settings.CONFIG_API_PATH, &apiKeys)
		if err != nil {
			log.Println("Error getting existing keys..", err)
			return false
		}

		for _, k := range apiKeys.Keys {

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
		if settings.VerifyAPIKey(key, providedKey) {
			return true
		}
	}
	return false
}
