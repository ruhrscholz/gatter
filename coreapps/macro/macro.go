package macro

import (
	"gatter/coreapps"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var m Macro
var db *gorm.DB

type Macro struct {
	Routes *http.ServeMux
}

func init() {
	m = Macro{
		Routes: http.NewServeMux(),
	}

	m.Routes.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	coreapps.Register("macro", &m)
}

func (m Macro) InitDb(db_ *gorm.DB) {
	db = db_
	err := db.AutoMigrate(&Post{})
	if err != nil {
		log.Panicf("macro: Database migration failed for Post model")
	}
	err = db.AutoMigrate(&PostRevision{})
	if err != nil {
		log.Panicf("macro: Database migration failed for PostRevision model")
	}
	err = db.AutoMigrate(&Tag{})
	if err != nil {
		log.Panicf("macro: Database migration failed for Tag model")
	}
}

func (m Macro) GetRoutes() *http.ServeMux {
	return m.Routes
}

func (m Macro) GetSlug() string {
	return "blog"
}
