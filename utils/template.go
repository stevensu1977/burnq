package utils

import (
	asset "app/views"
	"github.com/elazarl/go-bindata-assetfs"
	"html/template"
	"time"
)

const (
	TEMPLATE_PATH   = "./views"
	TEMPLATE_PREFIX = ".html"
)

var funcMap = template.FuncMap{
	"dateFormat": dateFormat,
	"safe":       safe,
}

var fs = assetfs.AssetFS{
	Asset:     asset.Asset,
	AssetDir:  asset.AssetDir,
	AssetInfo: asset.AssetInfo,
}

func dateFormat(t time.Time) string {
	return t.Format("2006/01/02 15:04:05")
}

func safe(source string) template.HTML {
	return template.HTML(source)
}

//AddFunc Privode helper func help add custom template func
func AddFunc(funcname string, handler interface{}) {
	funcMap[funcname] = handler
}

func AllFunc() []string {
	keys := make([]string, len(funcMap))
	for k := range funcMap {
		keys = append(keys, k)
	}
	return keys
}

//LoadTemplate is helper funcs load templatepath
func LoadTemplate(templateName string) (*template.Template, error) {

	//dev model use local static path
	//fmt.Println(filepath.Abs(templatePath))

	//dev model use local static path
	//return template.New(templateName+TEMPLATE_PREFIX).Delims("<<%", "%>>").Funcs(funcMap).ParseFiles(TEMPLATE_PATH + "/" + templateName + TEMPLATE_PREFIX)

	//release model use Asset
	data, err := fs.Asset("views/" + templateName + ".html")
	if err != nil {
		return nil, err
	}
	return template.New(templateName+TEMPLATE_PREFIX).Delims("<<%", "%>>").Funcs(funcMap).Parse(string(data))

}
