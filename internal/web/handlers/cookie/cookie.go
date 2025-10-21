package cookie

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/sessions"
)

type Session struct {
	Name      string
	ExpiresAt time.Time
	cookies   *sessions.CookieStore
}

var (
	mu           sync.Mutex
	UserSessions = make(map[string]Session)
)

const (
	SessionCookieName = "session_cookie"
	CsrfCookieName    = "session_csrf"
)

// NewLogin creates a new session and sets a cookie in the HTTP response
func NewLogin(w http.ResponseWriter, name string, sessionKey string, maxAgeSec int) string {

	expires := time.Now().Add(time.Duration(maxAgeSec) * time.Second)

	// Store session in memory
	mu.Lock()
	UserSessions[sessionKey] = Session{
		Name:      name,
		ExpiresAt: expires,
	}
	mu.Unlock()

	// Set the cookie for the browser
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionKey,
		Path:     "/",       // available to all paths
		MaxAge:   maxAgeSec, // in seconds
		HttpOnly: true,      // not accessible via JS
		Secure:   false,     // true if using HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return sessionKey
}

// ClearLogin invalidates a session
func ClearLogin(w http.ResponseWriter, key string) {
	mu.Lock()
	delete(UserSessions, key)
	mu.Unlock()

	// Delete cookie in browser
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
	})
}
func CurrentUser(r *http.Request) string {
	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		return ""
	}
	return UserSessions[c.Value].Name
}

// ValidateSession returns the username if the session is valid
func ValidateSession(r *http.Request) (string, bool) {
	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		fmt.Println(err)

		var s Session
		session, err := s.cookies.Get(r, "session")
		if auth, ok := session.Values["authenticated"].(bool); ok && auth {
			fmt.Println("Ok")
		}
		fmt.Println(err)
		return "", false
	}
	mu.Lock()
	defer mu.Unlock()
	sess, ok := UserSessions[c.Value]
	if !ok || time.Now().After(sess.ExpiresAt) {
		return "", false
	}

	return sess.Name, true
}

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; form-action 'self'; base-uri 'none'")
		next.ServeHTTP(w, r)
	})
}

// --- Sessions & CSRF ---

type SessionData struct {
	UID  int64
	Name string
	Exp  int64
}

var Sc = NewSecureCookie(NewAEADKey(NewSecret(32)))

func SetSession(w http.ResponseWriter, session SessionData) error {

	encoded, err := Sc.Encode(SessionCookieName, session)
	if err != nil {
		return err
	}

	c := http.Cookie{
		Name:     SessionCookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true, // readable by form
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	maxAge := int(time.Until(time.Unix(session.Exp, 0)).Seconds())
	NewLogin(w, session.Name, encoded, maxAge)
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	})
	csrf := make([]byte, 32)
	if _, err := rand.Read(csrf); err != nil {
		return err
	}

	c = http.Cookie{
		Name:     CsrfCookieName,
		Value:    base64.RawStdEncoding.EncodeToString(csrf),
		Path:     "/",
		HttpOnly: false, // readable by form
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &c)

	return nil
}
