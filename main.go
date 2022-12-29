package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go-seven/pkg/config"
	"go-seven/pkg/log"
	"go-seven/pkg/models"
	"go-seven/pkg/web"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func loadConfig() *models.Config {
	var file string
	flag.StringVar(&file, "c", "", "choose config file.")
	flag.Parse()

	var path = make([]string, 0)
	if file != "" {
		path = append(path, file)
	}
	err := config.New(path...)
	if err != nil {
		panic(err.Error())
	}

	return config.GetConfig()
}

func loadRouters(router *gin.Engine) {
	// TODO
}

func startWeb(cfg *models.Config) {
	mode := gin.ReleaseMode
	if cfg.Service.Mode == "dev" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	router := gin.New()

	logger := log.Default()
	router.Use(
		cors.Default(),
		web.Logger(logger),
		web.Recovery(logger),
	)

	loadRouters(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Web.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen ", log.String("err", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", log.String("err", err.Error()))
	}
}

func main() {
	cfg := loadConfig()

	log.ResetDefault(log.New(&cfg.Log))

	// 3. TODO: connect db

	startWeb(cfg)
}
