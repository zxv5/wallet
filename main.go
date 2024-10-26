package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wallet/internal/config"
	"wallet/internal/router"
	"wallet/pkg/config/setting"
	"wallet/pkg/logger"
)

var configFile = flag.String("f", "etc/config.yaml", "config file")

// @title wallet
// @version 1.0.0
// @description wallet
// @host 127.0.0.1:3000
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Please enter your Bearer token
func main() {
	flag.Parse()

	var conf config.Config
	setting.New(*configFile).Setup(&conf)
	logger.Init(nil)

	endPoint := fmt.Sprintf(":%d", conf.Server.Port)
	maxHeaderBytes := 1 << 20

	routersInit := router.Init(&conf)

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: maxHeaderBytes,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	go func() {
		logger.Infof("service start sucessful and listening on port :%s ...", endPoint)
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err.Error())
		}
	}()

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:10000", nil))
	// }()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server Shutdown: %s\n", err.Error())
	}

	logger.Info("Server exiting")
}
