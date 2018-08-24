package handler

import (
	"net/http"

	"github.com/stevensu1977/burnq/service"

	log "github.com/sirupsen/logrus"
)

func Index(w http.ResponseWriter, r *http.Request) {

	session, err := service.SessionStore().Get(r, SESSION_KEY)

	if err != nil {
		log.Error(err)
		http.Redirect(w, r, LOGIN_URL, 302)
		return
	}
	if _, ok := session.Values[SESSION_KEY]; !ok {
		log.Debugf("user need login %s", r.RemoteAddr)
		http.Redirect(w, r, LOGIN_URL, 302)
		return
	}
	http.Redirect(w, r, "/static/", 302)
	return

}
