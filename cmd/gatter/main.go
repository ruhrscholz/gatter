package main

import (
	"gatter/internal/endpoints/activitypub/users"
	"gatter/internal/endpoints/auth"
	"gatter/internal/endpoints/client"
	"gatter/internal/endpoints/oauth"
	"gatter/internal/endpoints/web"
	"gatter/internal/endpoints/wellknown"
	"gatter/internal/middleware"
	"gatter/internal/util/config"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	env, err := config.ReadConfig("./configs/config.json")
	if err != nil {
		log.Fatalf("Error while parsing config file: %s", err.Error())
		return
	}

	mux := http.NewServeMux()

	// .well-known
	mux.Handle("/.well-known/webfinger", wellknown.Webfinger(env))
	mux.Handle("/.well-known/nodeinfo", wellknown.Nodeinfo(env))

	// Auth
	mux.Handle("/auth/sign_in", middleware.Auth(env, auth.HandleSignIn(env)))

	// Oauth
	mux.Handle("/oauth/authorize", http.StripPrefix("/oauth/authorize", oauth.HandleAuthorize(env)))
	mux.Handle("/oauth/token", http.StripPrefix("/oauth/token", oauth.HandleToken(env)))

	// ActivityPub
	mux.Handle("/users/", http.StripPrefix("/users/", users.HandleUsers(env)))

	// Client
	client.Init(env)
	//mux.Handle("/api/v1/accounts/", http.StripPrefix("/api/v1/accounts/", client.HandleAccounts()))
	//mux.Handle("/api/v1/apps/", http.StripPrefix("/api/v1/apps/", client.HandleApps()))
	//mux.Handle("/api/v1/statuses/", http.StripPrefix("/api/v1/statuses/", client.HandleStatuses()))
	mux.Handle("/api/v1/timelines/public/", http.StripPrefix("/api/v1/timelines/public/", client.TimelinesPublic()))

	// Web
	mux.Handle("/", web.Handle(env))

	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatalf("Could not start http server: %s", err)
	}
}
