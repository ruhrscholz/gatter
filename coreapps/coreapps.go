package coreapps

import (
	"log"
	"net/http"
)

var (
	apps = make(map[string]CoreApp)
)

type CoreApp interface {
	GetRoutes() *http.ServeMux
	GetSlug() string
}

func Register(name string, coreapp CoreApp) {
	if coreapp == nil {
		log.Panicf("coreapps: Attempted coreapps registration with nil CoreApp")
	}
	if coreapp.GetSlug() == "" {
		log.Panicf("coreapps: Attempted registration of %s with empty slug", name)
	}
	if coreapp.GetRoutes() == nil {
		log.Panicf("coreapps: Attempted registration of %s with nil routes", name)
	}
	if _, dup := apps[name]; dup {
		log.Panicf("coreapps: Attempted double registration of %s CoreApp", name)
	}
	apps[name] = coreapp
}

func Apps() []string {
	keys := make([]string, len(apps))
	i := 0
	for k := range apps {
		keys[i] = k
		i++
	}
	return keys
}

func GetRoutes(app string) *http.ServeMux {
	return apps[app].GetRoutes()
}

func GetSlug(app string) string {
	return apps[app].GetSlug()
}
