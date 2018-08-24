package handler

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/stevensu1977/burnq/model"
	"github.com/stevensu1977/burnq/service"
	"github.com/stevensu1977/burnq/utils"
	toolbox "github.com/stevensu1977/toolbox/net"
)

func LoginGet(w http.ResponseWriter, r *http.Request) {

	template, err := utils.LoadTemplate("login")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = template.Execute(w, "")
	if err != nil {
		log.Error(err)
		http.Error(w, "Please contact service@flowq.io", 500)
		return
	}
}
func Account(w http.ResponseWriter, r *http.Request) {

	session, err := service.SessionStore().Get(r, SESSION_KEY)

	if err != nil {
		log.Error(err)
		http.Error(w, "NeedLogin", 500)
		return
	}

	if login, ok := session.Values[SESSION_KEY]; ok {

		toolbox.ServerJSON(w, map[string]string{"Email": login.(model.Admin).Email})
		return
	} else {
		http.Error(w, "NeedLogin", 500)
	}

}

func Login(w http.ResponseWriter, r *http.Request) {

	session, err := service.SessionStore().Get(r, SESSION_KEY)

	if err != nil {
		log.Error(err)
	}

	if login, ok := session.Values[SESSION_KEY]; ok {
		log.Printf("username %v already login", login)
		http.Redirect(w, r, "/", 302)
		return
	}

	r.ParseForm()
	log.Debug(r.Form)
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" {
		w.Write([]byte("User or Password not math, try again"))
		return
	}

	user, err := service.FetchAdmin(username, password)
	if err != nil {
		http.Redirect(w, r, LOGIN_URL, 302)
		log.Error(err)
		return
	}

	if user == nil {
		http.Redirect(w, r, LOGIN_URL, 302)
		return
	}

	// Set some session values.
	session.Values[SESSION_KEY] = *user

	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("User: %v login", user)

	http.Redirect(w, r, "/static/", 302)

}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := service.SessionStore().Get(r, SESSION_KEY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(session.Values, SESSION_KEY)
	session.Save(r, w)
	http.Redirect(w, r, LOGIN_URL, 302)
}
