package router

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"app/asset"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stevensu1977/burnq/handler"
)

const API_ROOT = "/api/v1"

func init() {

	log.Debug("Init Router succesful...")
	http.Handle("/", InitRouter())

}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func routerPathDebug(r *mux.Router, path string, f func(w http.ResponseWriter, r *http.Request)) {
	log.Debugf("Path: %s, Func: %s ", path, getFuncName(f))
	r.HandleFunc(path, f)
}

func APIPath(path string) string {
	return fmt.Sprintf("%s/%s", API_ROOT, path)
}

//InitRouter provide single router access
func InitRouter() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/", handler.Index).Methods("GET")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/login", handler.LoginGet).Methods("GET")
	r.HandleFunc("/logout", handler.Logout).Methods("GET")
	r.HandleFunc("/acc", handler.Account).Methods("GET")

	r.HandleFunc("/pref/passwd", handler.UpdateAdminPassword).Methods("POST")

	//Static Path load from AssetFS
	fs := assetfs.AssetFS{
		Asset:     asset.Asset,
		AssetDir:  asset.AssetDir,
		AssetInfo: asset.AssetInfo,
	}
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(&fs)))

	//dev model load static file from static directory
	//r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	//API
	r.HandleFunc("/api/v1", handler.APIVersion)

	r.HandleFunc(APIPath("auth"), handler.Auth)

	r.HandleFunc(APIPath("detail"), handler.Detail)

	r.HandleFunc(APIPath("usage"), handler.Usage)
	r.HandleFunc(APIPath("usage/daliy"), handler.UsageDaliy)

	r.HandleFunc(APIPath("account"), handler.FetchAllAccount).Methods("GET")
	r.HandleFunc(APIPath("account"), handler.CreateCloudAccount).Methods("POST")
	r.HandleFunc(APIPath("account"), handler.RemoveCloudAccount).Methods("DELETE")

	r.HandleFunc(APIPath("account"), handler.UpdateCloudAccountPasswd).Methods("PUT")

	r.HandleFunc(APIPath("tenant"), handler.FetchAllTenant).Methods("GET")

	r.Use(handler.AuthCheck)
	return r

}
