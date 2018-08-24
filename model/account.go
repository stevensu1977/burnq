package model

import (
	"encoding/json"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type CloudAccount struct {
	ID        int `storm:"unique,increment=10000"`
	Email     string
	Password  string `json:"Password,omitempty"`
	Tenant    string `storm:"unique"`
	TenantID  string `storm:"unique"`
	AccountID string `storm:"unique"`
}

type CloudDetail struct {
	Link  string `json:"canonicalLink,omitempty"`
	Items []struct {
		PurchasedResources []struct {
			ID        string         `json:"id,omitempty"`
			CreatedOn time.Time      `json:"createdOn,omitempty"`
			Purchases []PurchaseItem `json:"purchasedResources,omitempty"`
		} `json:"purchase,omitempy"`
	} `json:"items,omitempty"`
}

type PurchaseItem struct {
	Name     string    `json:"name,omitempty"`
	Value    float64   `json:"value,omitempty"`
	Unit     string    `json:"unit,omitempty"`
	Start    time.Time `json:"startDate,omitempty"`
	End      time.Time `json:"endDate,omitempty"`
	Purchase string    `json:"purchaseType,omitempty"`
	Using    float64   `json:"using"`
}

type UsageCost struct {
	AccountID string `json:"accountId,omitempty"`
	Items     []struct {
		Currency     string     `json:"currency,omitempty"`
		ProductID    string     `json:"gsiProductId,omitempty"`
		ResourceName string     `json:"resourceName,omitempty"`
		ServiceName  string     `json:"serviceName,omitempty"`
		DataCenter   string     `json:"dataCenterId,omitempty"`
		CostItems    []CostItem `json:"costs,omitemtpy"`
		EndTime      string     `json:"endTimeUtc,omitempty"`
		// CostItems    []struct {
		// 	Amount    float64 `json:"computedAmount,omitempty"`
		// 	Quantity  float64 `json:"computedQuantity,omitempty"`
		// 	UnitPrice float64 `json:"unitPrice,omitempty"`
		// } `json:"costs,omitemtpy"`
	} `json:"items,omitempty"`
}

type CostItem struct {
	Amount    float64 `json:"computedAmount"`
	Quantity  float64 `json:"computedQuantity,omitempty"`
	UnitPrice float64 `json:"unitPrice,omitempty"`
}

type Service struct {
	Amount float64
	Items  map[string]*CostItem
}

type ServiceLine struct {
	EndTime      string `json:"endTimeUtc,omitempty"`
	ServiceName  string
	ResourceName string
	Amount       float64
	Quantity     float64
	UnitPrice    float64
}

func (acc *CloudAccount) GetCacheKey(parameter ...string) string {
	if len(parameter) == 0 {
		return acc.AccountID
	}
	return acc.AccountID + "|" + strings.Join(parameter, "|")

}

func (cost *UsageCost) AggByService() map[string]*Service {
	service := make(map[string]*Service)

	for idx, _ := range cost.Items {
		serviceName := cost.Items[idx].ServiceName
		resourceName := cost.Items[idx].ResourceName
		if _, ok := service[serviceName]; !ok {
			service[serviceName] = &Service{Amount: 0, Items: make(map[string]*CostItem)}
		}
		if _, ok := service[serviceName].Items[resourceName]; !ok {
			service[serviceName].Items[resourceName] = &CostItem{Amount: 0.0, Quantity: 0.0, UnitPrice: 0.0}
		}

		for idx1, _ := range cost.Items[idx].CostItems {
			service[serviceName].Items[resourceName].Amount += cost.Items[idx].CostItems[idx1].Amount
			service[serviceName].Items[resourceName].Quantity += cost.Items[idx].CostItems[idx1].Quantity
			service[serviceName].Items[resourceName].UnitPrice = cost.Items[idx].CostItems[idx1].UnitPrice
			service[serviceName].Amount += cost.Items[idx].CostItems[idx1].Amount
		}

	}

	log.Debugf("%+v\n", service)

	return service
}

func (cost *UsageCost) AggByServiceLine() map[string][]ServiceLine {
	service := make(map[string][]ServiceLine)
	timeRecord := make(map[string]bool)
	for idx, _ := range cost.Items {
		if _, ok := service[cost.Items[idx].ServiceName+"."+cost.Items[idx].ResourceName]; !ok {
			service[cost.Items[idx].ServiceName+"."+cost.Items[idx].ResourceName] = []ServiceLine{}
		}

		if _, ok := timeRecord[cost.Items[idx].ServiceName+"."+cost.Items[idx].ResourceName+"."+cost.Items[idx].EndTime]; !ok {

			service[cost.Items[idx].ServiceName+"."+cost.Items[idx].ResourceName] = append(service[cost.Items[idx].ServiceName+"."+cost.Items[idx].ResourceName], ServiceLine{
				ServiceName:  cost.Items[idx].ServiceName,
				ResourceName: cost.Items[idx].ResourceName,
				EndTime:      cost.Items[idx].EndTime,
				Amount:       cost.Items[idx].CostItems[0].Amount,
				Quantity:     cost.Items[idx].CostItems[0].Quantity,
				UnitPrice:    cost.Items[idx].CostItems[0].UnitPrice,
			})
			timeRecord[cost.Items[idx].ServiceName+"."+cost.Items[idx].ResourceName+"."+cost.Items[idx].EndTime] = true

		} else {
			for subIdx, _ := range service[cost.Items[idx].ServiceName+"."+cost.Items[idx].ResourceName] {
				if service[cost.Items[idx].ServiceName+"."+cost.Items[idx].ResourceName][subIdx].EndTime == cost.Items[idx].EndTime {
					service[cost.Items[idx].ServiceName+"."+cost.Items[idx].ResourceName][subIdx].Amount += cost.Items[idx].CostItems[0].Amount
					service[cost.Items[idx].ServiceName+"."+cost.Items[idx].ResourceName][subIdx].Quantity += cost.Items[idx].CostItems[0].Quantity
				}
			}
		}

	}

	data, _ := json.Marshal(service)
	log.Debug(string(data))

	return service

}
