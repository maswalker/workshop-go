package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	h "go-seven/pkg/web"

	"github.com/gin-gonic/gin"
)

type HelloReq struct {
	Id int `json:"id" binding:"required"`
}

type HelloRes struct {
	Id int  `json:"id" binding:"required"`
	Ok bool `json:"ok" binding:"required"`
}

func Hello(req *HelloReq) (*HelloRes, error) {
	return &HelloRes{
		Id: req.Id,
		Ok: true,
	}, nil
}

func main() {
	router := gin.Default()

	router.POST("/hello", h.WrapperHandlerJson(Hello))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 10)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
