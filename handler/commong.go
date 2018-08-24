package handler

import (
	"net/http"
	"time"

	toolbox "github.com/stevensu1977/toolbox/net"
)

const (
	LOGIN_URL  = "/login"
	LOGOUT_URL = "/logout"

	//USERNAME , PASSWORD is static form key
	USERNAME    = "username"
	PASSWORD    = "password"
	SESSION_KEY = "user"

	UsageMaxCache  = 45 * time.Minute
	DetailMaxCache = 45 * time.Minute
)

func APIVersion(w http.ResponseWriter, r *http.Request) {
	toolbox.ServerJSON(w, map[string]string{
		"api":    "v1",
		"server": "BrunQ v0.0.1",
	})
}
