package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/web/handlers/cookie"
	"remoteCtrl/internal/web/handlers/login"
	"remoteCtrl/internal/web/handlers/status"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	username := cookie.CurrentUser(r)
	user_id := db.DataBase.Users.NameToID(username)
	is_admin, _ := db.DataBase.Users.IsAdmin(username)

	db.DataBase.Users.GetUserByID(user_id)
	avalable_tabs, err := db.DataBase.Tabs.GetAvailableTabsForUser(user_id)
	if err != nil {
		fmt.Println(err)
	}
	type resp struct {
		IsAdmin bool              `json:"is_admin"`
		Tabs    map[string]db.Tab `json:"tabs"`
	}
	var response resp = resp{
		IsAdmin: is_admin,
		Tabs:    avalable_tabs,
	} 
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }
	// if r.Method != http.MethodGet {
	// 	http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
}

func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	// if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }
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

	//modified := system.System.DB.Users.Modify(reqData.OldUsername, reqData.NewUsername, string(cookies.HashedPassword(reqData.NewPassword)))
	//if modified {
	//	db.Update(settings.CONFIG_USERS_PATH, system.System.DB.Users)
	//}
	//
	resp := status.Response{
		Message: "User modified!",
	}
	//for _, u := range system.System.DB.Users.Users {
	//	cookies.UserStore[u.Name] = u.Key
	//}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func AddUser(w http.ResponseWriter, r *http.Request) {
	// if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData login.RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if login.IsNotValid(reqData, w) != nil {
		return
	}
	// if system.System.DB.Users.Exists(reqData.Username) {
	// 	resp := status.Response{
	// 		Message: "User already exists!",
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(resp)
	// 	return
	// }

	// system.System.DB.Users.Add(reqData.Username, string(cookies.HashedPassword(reqData.Password)))
	// db.Update(settings.CONFIG_USERS_PATH, &system.System.DB.Users)

	resp := status.Response{
		Message: reqData.Username + " added!",
	}
	// for _, u := range system.System.DB.Users.Users {
	// 	cookies.UserStore[u.Name] = u.Key
	// }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
