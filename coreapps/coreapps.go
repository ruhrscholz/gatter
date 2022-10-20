package coreapps

import "net/http"

type App struct {
	Slug   string
	Routes *http.ServeMux
}
