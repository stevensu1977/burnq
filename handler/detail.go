package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

func Detail(w http.ResponseWriter, r *http.Request) {

	tenant := r.URL.Query().Get("tenant")
	if tenant == "" {
		http.Error(w, "Not tenant", http.StatusBadRequest)
		return
	}
	account, err := service.FetchAccount(tenant, true)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Debug("%+v", *account)

	if service.GlobalCache.IsExist(account.AccountID + "|detail") {
		toolbox.ServerJSON(w, service.GlobalCache.Get(account.GetCacheKey("detail")))
		return
	} else {
		service.GlobalCache.Put(account.GetCacheKey("detail"), struct{}{}, DetailMaxCache)
		go func() {
			client := service.NewClient(*account)
			detail, err := client.Detail()
			if err != nil {
				service.GlobalCache.Put(account.GetCacheKey("detail"), "500", DetailMaxCache)
				http.Error(w, err.Error(), 500)
				return
			} else {
				service.GlobalCache.Put(account.GetCacheKey("detail"), detail, DetailMaxCache)
			}

		}()
		//	ServerJSON(w, detail)

		toolbox.ServerJSON(w, struct{}{})
	}

}
