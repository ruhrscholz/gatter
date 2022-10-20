package main

import (
	"dos2/coreapps"
	_ "dos2/coreapps/macro"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type configOptions struct {
	SingleUser bool   `json:"singleUser"`
	Language   string `json:"language"`
}

func main() {
	configFile, err := os.ReadFile("./configs/config.development.json")
	if err != nil {
		log.Fatalf("Could not open config file: %s", err)
	}

	log.Println(string(configFile))

	var config configOptions
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Could not parse config file: %s", err)
	} else {
		log.Printf("Successfully read config file")
		log.Printf("Single user mode: %b", config.SingleUser)
		log.Printf("Language Code: %s", config.Language)
	}

	mux := http.NewServeMux()

	for _, k := range coreapps.Apps() {
		slug := coreapps.GetSlug(k)
		mux.Handle(
			fmt.Sprintf("/%s/", slug),
			http.StripPrefix(
				fmt.Sprintf("/%s", slug),
				coreapps.GetRoutes(k)))
	}

	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatalf("Could not start http server: %s", err)
	}
}
