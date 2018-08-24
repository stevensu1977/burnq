package service

import (
	log "github.com/sirupsen/logrus"

	"github.com/astaxie/beego/cache"
)

var GlobalCache cache.Cache

func init() {
	var err error
	GlobalCache, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		log.Fatal(err)
	}

}
