package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"remoteCtrl/internal/system/cookies"

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
