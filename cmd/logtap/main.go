package main

import (
	"log"
	"net/http"
	"os"

	"github.com/lichuan0620/logtap/cmd/logtap/option"
	"github.com/lichuan0620/logtap/pkg/logtap"
	"github.com/lichuan0620/logtap/pkg/logtap/handler"
)

func main() {
	logTap, err := logtap.NewLogTap(option.Spec, option.Name)
	if err != nil {
		os.Exit(1)
	}
	go serveHTTP(logTap)
	if err = logTap.Run(option.StopCh); err != nil {
		panic(err)
	}
}

func serveHTTP(tap logtap.LogTap) {
	if err := http.ListenAndServe(option.WebAddress, handler.NewLogTapHandler(tap)); err != nil {
		log.Fatalln(err.Error())
	}
}
