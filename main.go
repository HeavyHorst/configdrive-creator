package main

import (
	"log"
	"net/http"

	"github.com/rs/xaccess"
	"github.com/rs/xhandler"
	"github.com/rs/xlog"
	"github.com/rs/xmux"
)

var (
	mw    xhandler.Chain
	cfg   config
	mkiso *mkisofs
)

func init() {
	cfg.initDefaultConfig()
	mkiso = initMkisofs()

	//plug the xlog handler's input to Go's default logger
	log.SetFlags(0)
	xlogger := xlog.New(cfg.loggerConfig)
	log.SetOutput(xlogger)

	mw.UseC(xlog.NewHandler(cfg.loggerConfig))
	mw.UseC(xaccess.NewHandler())
}

func main() {
	router := xmux.New()

	router.GET("/", xhandler.HandlerFuncC(indexHandler))
	router.POST("/configdrive", xhandler.HandlerFuncC(configdriveHandler))
	log.Fatal(http.ListenAndServe(":3000", mw.Handler(router)))
}
