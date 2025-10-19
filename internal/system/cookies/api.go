package cookies

import (
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/system/settings"

	"encoding/json"
	"log"
	"net/http"
)

type api_response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Key     string `json:"key"`
}

func GenAPIKeyHandler(apiKeys settings.API_secrets, w http.ResponseWriter, r *http.Request) {
	// if !Session.IsLoggedIn(apiKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }

	err := db.LoadConfig(settings.CONFIG_API_PATH, &apiKeys)
	if err != nil {
		log.Println("Error getting existing keys..", err)
		http.Error(w, "Error getting existing keys..", http.StatusBadRequest)
		return
	}

	session, err := Session.Store().Get(r, "session")
	new_api_config := apiKeys.NewKey()
	new_api_config.User = session.Values["user"].(string)

	if new_api_config.User == "" {
		http.Error(w, "Unable generate api keys..", http.StatusForbidden)
		return
	}

	new_api_config.Name = r.URL.Query().Get("name")

	for _, k := range apiKeys.Keys {
		if k.Name == new_api_config.Name {
			if err != nil {
				http.Error(w, "Named key already exists!", http.StatusConflict)
				return
			}
		}
	}

	key, err := new_api_config.GenerateAPIKey(32)
	if err != nil {
		http.Error(w, "Unable generate api keys..", http.StatusBadRequest)
		return
	}
	hashedKey, err := new_api_config.HashAPIKey(key)

	new_api_config.Key = hashedKey
	apiKeys.Keys = append(apiKeys.Keys, new_api_config)
	err = db.Write(settings.CONFIG_API_PATH, apiKeys)
	if err != nil {
		http.Error(w, "error saving new key..", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := api_response{Status: true, Message: "Generated api key.", Key: key}
	json.NewEncoder(w).Encode(response)
}

func GetAPIkeys(apiKeys settings.API_secrets, w http.ResponseWriter, r *http.Request) {
	// if !Session.IsLoggedIn(apiKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }

	err := db.LoadConfig(settings.CONFIG_API_PATH, &apiKeys)
	if err != nil {
		log.Println("Error getting existing keys..", err)
		http.Error(w, "Error getting existing keys..", http.StatusBadRequest)
		return
	}

	type data struct {
		Name string `json:"name"`
	}
	var apiList []data
	for _, k := range apiKeys.Keys {
		apiList = append(apiList, data{Name: k.Name})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiList)
}

func DeleteAPIKeyHandler(apiKeys settings.API_secrets, w http.ResponseWriter, r *http.Request) {
	// if !Session.IsLoggedIn(apiKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }

	var tmp_secrets settings.API_secrets

	type data struct {
		Name string `json:"new"`
	}
	type request struct {
		Data data `json:"data"`
	}
	var reqData request

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		log.Println("Invalid request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := db.LoadConfig(settings.CONFIG_API_PATH, &apiKeys)
	if err != nil {
		log.Println("Error getting existing keys..", err)
		http.Error(w, "Error getting existing keys..", http.StatusBadRequest)
		return
	}

	session, err := Session.Store().Get(r, "session")
	if err != nil {
		log.Println("Error getting session..", err)
		http.Error(w, "Error getting session..", http.StatusBadRequest)
		return
	}

	username := session.Values["user"].(string)

	for _, k := range apiKeys.Keys {
		if k.Name == reqData.Data.Name && k.User == username {
			continue
		}
		tmp_secrets.Keys = append(tmp_secrets.Keys, k)
	}

	apiKeys.Keys = tmp_secrets.Keys
	err = db.Write(settings.CONFIG_API_PATH, apiKeys)
	if err != nil {
		log.Println("Error saving new key..", err)
		http.Error(w, "error saving new key..", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := api_response{Status: true, Message: "Deleted api key.", Key: "nil"}
	json.NewEncoder(w).Encode(response)
}
