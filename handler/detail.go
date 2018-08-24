package handler

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/stevensu1977/burnq/service"
	toolbox "github.com/stevensu1977/toolbox/net"
)

func Detail(w http.ResponseWriter, r *http.Request) {

	tenant := r.URL.Query().Get("tenant")
	if tenant == "" {
		tenant = "xjbank"
	}
	account, err := service.FetchAccount(tenant, true)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Debug("%+v", *account)

	if service.GlobalCache.IsExist(account.AccountID + "|detail") {
		toolbox.ServerJSON(w, service.GlobalCache.Get(account.GetCacheKey("detail")))
	} else {
		service.GlobalCache.Put(account.GetCacheKey("detail"), struct{}{}, DetailMaxCache)
		go func() {
			client := service.NewClient(*account)
			detail, err := client.Detail()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			service.GlobalCache.Put(account.GetCacheKey("detail"), detail, DetailMaxCache)

		}()
		//	ServerJSON(w, detail)

		toolbox.ServerJSON(w, struct{}{})
	}

}
