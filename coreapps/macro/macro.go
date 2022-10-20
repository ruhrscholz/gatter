package macro

import (
	"dos2/coreapps"
	"net/http"
)

type Macro struct {
}

func init() {
	coreapps.Register("macro", &Macro{})
}

func (m Macro) GetRoutes() *http.ServeMux {
	return nil
}

func (m Macro) GetSlug() string {
	return "blog"
}
