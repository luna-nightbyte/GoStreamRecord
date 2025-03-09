package users

import (
	"GoStreamRecord/internal/db"
	"GoStreamRecord/internal/file"
	"GoStreamRecord/internal/handlers/cookies"
	"GoStreamRecord/internal/handlers/login"
	"GoStreamRecord/internal/handlers/status"
	"encoding/json"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(db.Config.Users.Users)
}

func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		OldUsername string `json:"oldUsername"`
		NewUsername string `json:"newUsername"`
		NewPassword string `json:"newPassword"`
	}
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	modified := db.Config.Users.Modify(reqData.OldUsername, reqData.NewUsername, string(login.HashedPassword(reqData.NewPassword)))
	if modified {
		db.Config.Update("users", file.Users_json, db.Config.Users)
	}

	resp := status.Response{
		Message: "User modified!",
	}
	for _, u := range db.Config.Users.Users {
		cookies.UserStore[u.Name] = u.Key
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func AddUser(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData login.RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if login.IsNotValid(reqData, w) {
		return
	}
	if db.Config.Users.Exists(reqData.Username) {
		resp := status.Response{
			Message: "User already exists!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	db.Config.Users.Add(reqData.Username, string(login.HashedPassword(reqData.Password)))
	db.Config.Update("users", file.Users_json, &db.Config.Users)

	resp := status.Response{
		Message: reqData.Username + " added!",
	}
	for _, u := range db.Config.Users.Users {
		cookies.UserStore[u.Name] = u.Key
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
