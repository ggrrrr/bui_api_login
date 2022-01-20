package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ggrrrr/bui_api_login/cli"
	"github.com/ggrrrr/bui_api_login/cli/users"
	"github.com/ggrrrr/bui_api_login/controlers/auth"
	res "github.com/ggrrrr/bui_api_login/resources/auth"

	"github.com/ggrrrr/bui_lib/api"
	db "github.com/ggrrrr/bui_lib/db/cassandra"
	"github.com/ggrrrr/bui_lib/token"
	"github.com/ggrrrr/bui_lib/token/sign"
)

var (
	root = context.Background()
	err  error

	commands = map[string]cli.CliCommand{}
	command  string
)

func init() {
	flag.StringVar(&command, "cli", "", "CLI commands users")
	commands["users"] = users.New()
}

func main() {
	flag.Parse()

	auth.Configure()

	err = db.Configure()
	if err != nil {
		log.Fatalf(err.Error())
	}
	session, err := db.Connect()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer session.Close()
	db.CreateSchema("passwd")
	defer db.Session.Close()

	if command != "" {
		v, ok := commands[command]
		if !ok {
			fmt.Println("command not found.")
		}
		v.Exec()
		fmt.Println("done.")
		return
	}

	server()
	// log.Printf("end.")
}

func server() {

	if token.Configure() != nil {
		log.Fatalln("token config")
	}

	if sign.Configure() != nil {
		log.Fatalf(err.Error())
	}

	err = api.Configure()
	if err != nil {
		log.Fatalf(err.Error())
	}

	go func() {
		time.Sleep(5 * time.Second)
		api.Ready()
	}()

	err = api.Create(root, false)
	if err != nil {
		log.Fatalf(err.Error())
	}

	api.HandleFunc("/auth/login/user", res.LoginUserRequest)
	api.HandleFunc("/auth/login/oauth2", res.LoginOauth2Request)
	api.HandleFunc("/auth/token", res.TokenVerifyRequest)
	api.HandleFunc("/auth/me/password", res.ChangePasswordRequest)

	api.HandleFunc("/auth/request/list", res.ListRequest)
	api.HandleFunc("/auth/request/activate", res.ActivateRequest)

	osSignals := make(chan os.Signal, 1)
	go func() {
		err := api.Start()
		defer api.Shutdown()
		if err != nil {
			log.Printf("http error: %+v", err)
			osSignals <- os.Kill
		}
	}()

	signal.Notify(osSignals, os.Interrupt)
	log.Printf("os.signal: %v", <-osSignals)
	api.Shutdown()
	db.Session.Close()
	log.Printf("end.")

}
