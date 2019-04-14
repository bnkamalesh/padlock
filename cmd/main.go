package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bnkamalesh/padlock/api"
	"github.com/bnkamalesh/padlock/api/http"
	"github.com/bnkamalesh/padlock/pkg/appcontext"
	"github.com/bnkamalesh/padlock/pkg/apps"
	"github.com/bnkamalesh/padlock/pkg/platform/cache"
	"github.com/bnkamalesh/padlock/pkg/platform/logger"
	"github.com/bnkamalesh/padlock/pkg/totp"
	"github.com/bnkamalesh/padlock/pkg/users"
)

func main() {
	l := logger.New("*")

	cacheHandler, err := cache.New(cache.Config{
		Hosts:        []string{""},
		DialTimeout:  time.Second * 3,
		ReadTimeout:  time.Millisecond * 100,
		WriteTimeout: time.Second * 1,
	})

	if err != nil {
		l.Fatal(err)
		return
	}

	pgdb, err := postgresHandler(postgreConfig{
		Host:     "127.0.0.1",
		Port:     "5432",
		Username: "padlockpostgres",
		Password: "z1A3hell0WorLd553",
		DBName:   "padlock",
	})
	if err != nil {
		l.Fatal(err)
		return
	}

	appCtx := appcontext.New(l)
	appCtx.Logging = true
	if appCtx.Logging {
		if !appCtx.Debug {
			l = logger.New("info", "warn", "error", "fatal")
		} else {
			l = logger.New("*")
		}
		appCtx.Logger = l
	}

	appsHandler := apps.New(appCtx, pgdb)
	usersHandler := users.New(appCtx, pgdb, cacheHandler)

	test(l, appsHandler, usersHandler)

	api := api.New(appCtx, appsHandler, usersHandler)

	httpServer, err := http.NewServer(
		"",
		"8080",
		"localhost",
		api,
		appCtx,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = httpServer.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func test(l logger.Logger, aH *apps.Apps, uH *users.Users) {
	ctx := context.Background()
	u, err := uH.Create(
		ctx,
		users.User{
			Email: "bnkamalesh@gmail.com",
			Phone: "+91-8792485305",
			Name:  "KBN",
		},
		"123456",
	)
	if err != nil {
		l.Error(err)
		return
	}

	a, err := aH.CreateAndSetOwner(
		ctx,
		apps.App{
			Name:        "KBN-App",
			Description: "",
			TOTP:        totp.New("KBN-App", 6, 30, totp.AlgoSHA1),
		},
		*u,
	)

	if err != nil {
		l.Error(err)
		return
	}
	fmt.Println("Created app!", a.Name)

	u, _, err = uH.Login(ctx, "localhost", "bnkamalesh@gmail.com", "123456")
	if err != nil {
		l.Error(err)
		return
	}
	fmt.Println("Logged in!")
}
