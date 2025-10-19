package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/system/cookies"
	"remoteCtrl/internal/web/handlers/cookie"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	storedHash, ok := cookies.UserStore[username]
	if !ok {

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Invalid credentials")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Invalid credentials")
		return
	}

	session, err := cookies.Session.Store().Get(r, "session")
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Session error. Try clearing your cookies.")
		return
	}
	session.Values["authenticated"] = true
	session.Values["user"] = username
	if err := session.Save(r, w); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Could not save session")
		//http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// Verifies that username only contain letters, numbers, and or underscores
func ValidUsername(username string) bool {
	for _, c := range username {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' {
			return false
		}
	}
	return true
}

type RequestData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func IsNotValid(reqData RequestData, w http.ResponseWriter) (outErr error) {
	switch true {
	case len(reqData.Username) == 0 || len(reqData.Password) == 0:
		outErr = fmt.Errorf("username and password cannot be empty")
	case len(reqData.Username) < 3 || len(reqData.Password) < 3:
		outErr = fmt.Errorf("Username and password must be at least 3 characters long!")
	case len(reqData.Username) > 20 || len(reqData.Password) > 20:
		outErr = fmt.Errorf("Username and password must be at most 20 characters long!")
	case !ValidUsername(reqData.Username):
		outErr = fmt.Errorf("Username can only contain letters, numbers, and underscores!")
	}
	return outErr
}

// RequireAuth wraps a handler and redirects to /login if not logged in
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := cookie.ValidateSession(r); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

type User struct {
	ID       int64
	Username string
	Role     string
	PassHash []byte
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// a.renderLogin(w, r, map[string]any{"Error": ""})
	case http.MethodPost:
		user := strings.TrimSpace(r.FormValue("username"))
		pass := r.FormValue("password")
		if user == "" || pass == "" {

			http.Redirect(w, r, "/login", http.StatusBadRequest)
			//  a.renderLogin(w, r, map[string]any{"Error": "Missing credentials"})
			return
		}
		var u User
		err := db.DataBase.SQL.QueryRow(`SELECT id, username, password_hash FROM users WHERE username = ?`, user).
			Scan(&u.ID, &u.Username, &u.PassHash)
		if err != nil || cookie.CheckHash(u.PassHash, pass) != nil {
			http.Redirect(w, r, "/login", http.StatusBadRequest)
			return
		}
		exp := time.Now().Add(8 * time.Hour).Unix()

		if err := cookie.SetSession(w, cookie.SessionData{UID: u.ID, Name: u.Username, Exp: exp}); err != nil {
			http.Error(w, "session error", http.StatusInternalServerError)
			return
		}
		fmt.Printf("'%s' logged in", user)
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	c, err := r.Cookie(cookie.SessionCookieName)
	if err != nil {
		return
	}
	cookie.ClearLogin(w, c.Value)
	http.SetCookie(w, &http.Cookie{Name: cookie.SessionCookieName, Path: "/", HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode, MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: cookie.CsrfCookieName, Path: "/", HttpOnly: false, Secure: true, SameSite: http.SameSiteLaxMode, MaxAge: -1})

	http.Redirect(w, r, "/login", http.StatusFound)
}
