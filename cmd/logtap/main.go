package main

import (
	"github.com/lichuan0620/logtap/cmd/logtap/option"
	"github.com/lichuan0620/logtap/pkg/logtap"
)

func main() {
	logTap, err := logtap.NewLogTap(option.Spec, option.Name)
	if err != nil {
		panic(err)
	}
	if err = logTap.Run(option.StopCh); err != nil {
		panic(err)
	}
}
