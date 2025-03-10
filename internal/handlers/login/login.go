package login

import (
	"GoStreamRecord/internal/handlers/cookies"
	"GoStreamRecord/internal/handlers/status"
	web_status "GoStreamRecord/internal/handlers/status"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	resp := web_status.Response{}

	storedHash, ok := cookies.UserStore[username]
	if !ok {

		resp.Message = "Invalid credentials"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		fmt.Println(err)
		resp.Message = "Invalid credentials"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	session, err := cookies.Session.Store().Get(r, "session")
	if err != nil {
		resp.Message = "Session error. Try clearing your cookies."

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	session.Values["authenticated"] = true
	session.Values["user"] = username
	if err := session.Save(r, w); err != nil {
		resp.Message = "Could not save session"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		fmt.Println(err)
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

func IsNotValid(reqData RequestData, w http.ResponseWriter) bool {
	if len(reqData.Username) == 0 || len(reqData.Password) == 0 {
		resp := status.Response{
			Message: "Username and password cannot be empty!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return true
	}

	if len(reqData.Username) < 3 || len(reqData.Password) < 3 {
		resp := status.Response{
			Message: "Username and password must be at least 3 characters long!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return true
	}

	if len(reqData.Username) > 20 || len(reqData.Password) > 20 {
		resp := status.Response{
			Message: "Username and password must be at most 20 characters long!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return true
	}

	if !ValidUsername(reqData.Username) {
		resp := status.Response{
			Message: "Username can only contain letters, numbers, and underscores!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return true
	}

	return false
}
