package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/stevensu1977/toolbox/crypto"

	log "github.com/sirupsen/logrus"

	"github.com/stevensu1977/burnq/model"
)

const UCMeterAPI_Root = "https://itra.oraclecloud.com//metering/api/v1"
const UCMeterAPI_Detail = "cloudbucks"
const UCMeterAPI_UsageCost = "usagecost"

const Query_HOURLY = "HOURLY"
const Query_DAILY = "DAILY"
const Query_MONTHLY = "MONTHLY"

var wg sync.WaitGroup

type UCMeterClient struct {
	Account model.CloudAccount
	Request *http.Request
}

type DateWrapp struct {
	Year  int
	Month int
	Days  int
	Start time.Time
	End   time.Time
}

func (uc *UCMeterClient) Detail() (*model.CloudDetail, error) {
	var result bytes.Buffer
	var err error
	var endpoint *url.URL

	endpoint, err = url.Parse(fmt.Sprintf("%s/%s/%s", UCMeterAPI_Root, UCMeterAPI_Detail, uc.Account.AccountID))

	if err != nil {
		return nil, err
	}
	uc.Request, err = http.NewRequest("GET", endpoint.String(), &result)

	if err != nil {
		return nil, err
	}
	uc.Auth()

	log.Debug(uc.Request)

	client := &http.Client{}
	resp, err := client.Do(uc.Request)
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	detail := &model.CloudDetail{}
	err = json.Unmarshal(data, detail)
	if err != nil {
		panic(err)
	}
	if len(detail.Items) == 0 {
		fmt.Println(string(data))
		return nil, fmt.Errorf("Not found items")
	}
	sort.Slice(detail.Items[0].PurchasedResources, func(i, j int) bool {
		return detail.Items[0].PurchasedResources[i].Purchases[0].Start.Unix() < detail.Items[0].PurchasedResources[j].Purchases[0].Start.Unix()
	})

	data, _ = json.MarshalIndent(detail, "", " ")

	log.Debug(string(data))

	for idx, _ := range detail.Items[0].PurchasedResources {
		for idx1, _ := range detail.Items[0].PurchasedResources[idx].Purchases {

			service, err := uc.UsageCost(detail.Items[0].PurchasedResources[idx].Purchases[idx1].Start, detail.Items[0].PurchasedResources[idx].Purchases[idx1].End, Query_MONTHLY)
			if err != nil {
				panic(err)
			}
			amount := 0.0
			for _, v := range service {
				amount += v.Amount
			}
			detail.Items[0].PurchasedResources[idx].Purchases[idx1].Using = amount
			log.Debugf("%s, %f, %f\n", detail.Items[0].PurchasedResources[idx].ID, detail.Items[0].PurchasedResources[idx].Purchases[idx1].Value, amount)

		}

	}
	log.Infof("------------%s load Purchases Complete------------", uc.Account.Tenant)

	return detail, nil
}

func (uc *UCMeterClient) UsageCostCurrentMonth() (map[string]*model.Service, error) {
	now := time.Now()
	times := GetMonth(now.Year(), int(now.Month()))
	return uc.UsageCost(times[0], times[1], Query_MONTHLY)

}

func (uc *UCMeterClient) UsageCostMonth(year, month int) (map[string]*model.Service, error) {
	times := GetMonth(year, month)
	return uc.UsageCost(times[0], times[1], Query_MONTHLY)

}

func (uc *UCMeterClient) UsageCostCurrentMonthDetail() (map[string][]model.ServiceLine, error) {
	now := time.Now()
	times := GetMonth(now.Year(), int(now.Month()))
	return uc.UsageCostDetail(times[0], times[1], Query_MONTHLY)

}

func (uc *UCMeterClient) UsageCostMonthDetail(year, month int) (map[string][]model.ServiceLine, error) {
	times := GetMonth(year, month)
	return uc.UsageCostDetail(times[0], times[1], Query_MONTHLY)

}

func (uc *UCMeterClient) UsageCost(startTime, endTime time.Time, query string) (map[string]*model.Service, error) {

	var result bytes.Buffer
	var err error
	var endpoint *url.URL
	//HOURLY,DAILY, MONTHLY
	endpoint, err = url.Parse(fmt.Sprintf("%s/%s/%s?startTime=%sT00:00:00.000Z&endTime=%sT24:00:00.000Z&computeTypeEnabled=Y&dcAggEnabled=Y&timeZone=America/Los_Angeles&usageType=%s", UCMeterAPI_Root, UCMeterAPI_UsageCost, uc.Account.AccountID, Time2Date(startTime), Time2Date(endTime), query))

	log.Debug(endpoint)
	if err != nil {
		return nil, err
	}
	uc.Request, err = http.NewRequest("GET", endpoint.String(), &result)

	if err != nil {
		return nil, err
	}
	uc.Auth()

	client := &http.Client{}
	resp, err := client.Do(uc.Request)
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	cost := &model.UsageCost{}
	json.Unmarshal(data, cost)

	aggReuslt := cost.AggByService()
	allCost, err := json.Marshal(aggReuslt)

	log.Debug(string(data))
	log.Debug(string(allCost))

	log.Infof("------------%s load Usage Complete------------", uc.Account.Tenant)

	return aggReuslt, nil
}

func (uc *UCMeterClient) Auth() error {
	if uc.Request == nil {
		return fmt.Errorf("wrong auth , you need first init request")
	}
	uc.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", crypto.BasicAuthEncode(uc.Account.Email, uc.Account.Password)))
	uc.Request.Header.Set("X-ID-TENANT-NAME", uc.Account.TenantID)
	return nil
}

func Time2Date(t time.Time) string {
	return t.Format("2006-01-02")
}

func GetDays(year, month int) int {
	if month > 12 {
		return -1
	}
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}

func GetMonth(year, month int) []time.Time {
	if month > 12 {
		return nil
	}
	times := []time.Time{}
	startTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	times = append(times, startTime)
	times = append(times, time.Date(year, time.Month(month), GetDays(year, month)-1, 24, 0, 0, 0, time.UTC))

	log.Debugf("GetMonth %v", times)
	return times
}

func BuildStart2End(start, end string) []time.Time {

	times := []time.Time{}
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		panic(err)
	}
	endTime, err := time.Parse("2006-01-02  15:04:05", end+"  23:59:59")

	if err != nil {
		panic(err)
	}

	if endTime.Before(startTime) {
		return nil
	}

	times = append(times, startTime)
	times = append(times, endTime)

	log.Debugf("%s, %s,  %v", start, end, times)

	return times

}

func (uc *UCMeterClient) UsageCostDetail(startTime, endTime time.Time, query string) (map[string][]model.ServiceLine, error) {

	var result bytes.Buffer
	var err error
	var endpoint *url.URL
	//HOURLY,DAILY, MONTHLY
	endpoint, err = url.Parse(fmt.Sprintf("%s/%s/%s?startTime=%sT00:00:00.000Z&endTime=%sT24:00:00.000Z&computeTypeEnabled=Y&dcAggEnabled=Y&timeZone=America/Los_Angeles&usageType=DAILY", UCMeterAPI_Root, UCMeterAPI_UsageCost, uc.Account.AccountID, Time2Date(startTime), Time2Date(endTime)))

	log.Debug(endpoint)

	if err != nil {
		return nil, err
	}
	uc.Request, err = http.NewRequest("GET", endpoint.String(), &result)

	if err != nil {
		return nil, err
	}
	uc.Auth()

	client := &http.Client{}
	resp, err := client.Do(uc.Request)
	data, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	cost := &model.UsageCost{}

	log.Debug(string(data))

	json.Unmarshal(data, cost)

	return cost.AggByServiceLine(), nil
}

func NewClient(account model.CloudAccount) *UCMeterClient {
	return &UCMeterClient{Account: account}
}
