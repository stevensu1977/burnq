package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/stevensu1977/burnq/model"
	"github.com/stevensu1977/burnq/service"

	toolbox "github.com/stevensu1977/toolbox/net"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	var account model.CloudAccount

	json.Unmarshal(data, &account)

	client := service.NewClient(account)
	if client.CheckAuth() {
		w.Write([]byte("ok"))
	} else {
		w.WriteHeader(401)
		http.Error(w, "401", http.StatusUnauthorized)
	}

}

func UpdateCloudAccountPasswd(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	var account model.CloudAccount

	json.Unmarshal(data, &account)

	client := service.NewClient(account)
	if client.CheckAuth() {

		if err := service.UpdateAccountPasswd(&account); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("ok"))

	} else {
		w.WriteHeader(401)
		http.Error(w, "401", http.StatusUnauthorized)
	}

}

func CreateCloudAccount(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	var account model.CloudAccount

	json.Unmarshal(data, &account)

	err = service.SaveAccount(&account)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	toolbox.ServerJSON(w, account)
}

func FetchAllTenant(w http.ResponseWriter, r *http.Request) {
	tenants, err := service.FetchAllTenant()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if len(tenants) == 0 {
		toolbox.ServerJSON(w, struct{}{})
		return
	}
	toolbox.ServerJSON(w, tenants)
}

func FetchAllAccount(w http.ResponseWriter, r *http.Request) {
	tenants, err := service.FetchAllAccount()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for idx, _ := range tenants {
		tenants[idx].Password = ""
	}
	toolbox.ServerJSON(w, tenants)
}

func RemoveCloudAccount(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(500)
		return
	}
	_id, err := strconv.Atoi(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	err = service.RemoveAccount(_id)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	toolbox.ServerJSON(w, map[string]string{"status": "ok"})
}

func UpdateAdminPassword(w http.ResponseWriter, r *http.Request) {

	session, err := service.SessionStore().Get(r, SESSION_KEY)

	if err != nil {
		log.Error(err)
	}

	r.ParseForm()
	log.Debug(r.Form)
	newPassword := r.FormValue("newpassword")
	password := r.FormValue("password")

	if password == "" || newPassword == "" {
		w.WriteHeader(500)
		return
	}

	user, ok := session.Values[SESSION_KEY]

	if !ok {
		w.WriteHeader(500)
		return
	}

	admin, err := service.FetchAdmin(user.(model.Admin).Email, password)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	admin.Password = newPassword

	err = service.UpdateAdminPasswd(admin)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	toolbox.ServerJSON(w, map[string]string{"status": "ok"})
}
