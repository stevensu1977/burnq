package handler

import (
	"fmt"
	"net/http"
	"time"

	"encoding/csv"

	log "github.com/sirupsen/logrus"
	toolbox "github.com/stevensu1977/toolbox/net"

	"github.com/stevensu1977/burnq/model"
	"github.com/stevensu1977/burnq/service"
)

func Usage(w http.ResponseWriter, r *http.Request) {

	//client.UsageCostCurrentMonth()
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")
	tenant := r.URL.Query().Get("tenant")
	if tenant == "" {
		tenant = "xjbank"
	}
	account, err := service.FetchAccount(tenant, true)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	client := service.NewClient(*account)

	log.Debug(startTime, endTime)

	if service.GlobalCache.IsExist(account.GetCacheKey(startTime, endTime)) {
		toolbox.ServerJSON(w, service.GlobalCache.Get(account.GetCacheKey(startTime, endTime)))
	} else {
		if startTime != "" && endTime != "" {
			times := service.BuildStart2End(startTime, endTime)
			usage, err := client.UsageCost(times[0], times[1], service.Query_DAILY)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			service.GlobalCache.Put(account.GetCacheKey(startTime, endTime), usage, UsageMaxCache)
			toolbox.ServerJSON(w, usage)
		} else {
			usage, err := client.UsageCostCurrentMonth()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			service.GlobalCache.Put(account.GetCacheKey(startTime, endTime), usage, UsageMaxCache)
			toolbox.ServerJSON(w, usage)
		}

	}

}

func UsageDaliy(w http.ResponseWriter, r *http.Request) {

	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")
	export := r.URL.Query().Get("export")

	tenant := r.URL.Query().Get("tenant")
	if tenant == "" {
		tenant = "xjbank"
	}
	account, err := service.FetchAccount(tenant, true)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	client := service.NewClient(*account)

	log.Debug(startTime, endTime)

	if service.GlobalCache.IsExist(account.AccountID + "|daliy|" + startTime + "|" + endTime) {
		if export == "" {
			toolbox.ServerJSON(w, service.GlobalCache.Get(account.AccountID+"|daliy|"+startTime+"|"+endTime))
		} else {
			csvOut(w, service.GlobalCache.Get(account.AccountID+"|daliy|"+startTime+"|"+endTime).(map[string][]model.ServiceLine))

		}

	} else {
		if startTime != "" && endTime != "" {
			times := service.BuildStart2End(startTime, endTime)
			usage, err := client.UsageCostDetail(times[0], times[1], service.Query_DAILY)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			log.Printf("%s  From %s to %s Detail Usage load complete.", account.Tenant, startTime, endTime)
			service.GlobalCache.Put(account.AccountID+"|daliy|"+startTime+"|"+endTime, usage, 50*time.Minute)
			if export == "" {
				toolbox.ServerJSON(w, usage)
			} else {
				csvOut(w, usage)

			}

		} else {
			usage, err := client.UsageCostCurrentMonthDetail()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			service.GlobalCache.Put(account.AccountID+"|daliy|"+startTime+"|"+endTime, usage, 30*time.Minute)
			if export == "" {
				toolbox.ServerJSON(w, usage)
			} else {
				csvOut(w, usage)

			}
		}

	}

}

func csvOut(w http.ResponseWriter, daliy map[string][]model.ServiceLine) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "daliy.csv"))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	w.Header().Set("Cache-Control", "must-revalidate")
	csvWriter := csv.NewWriter(w)
	csvWriter.Write([]string{"EndTime", "ServiceName", "ResourceName", "Amount", "Quantity", "UnitPrice"})
	for _, v := range daliy {
		for _, v1 := range v {
			log.Debug(v1.EndTime, v1.ServiceName, v1.ResourceName, v1.Amount, v1.Quantity, v1.UnitPrice)
			csvWriter.Write([]string{v1.EndTime, v1.ServiceName, v1.ResourceName, fmt.Sprintf("%.2f", v1.Amount), fmt.Sprintf("%.2f", v1.Quantity), fmt.Sprintf("%.2f", v1.UnitPrice)})
		}
	}

	csvWriter.Flush()
}
