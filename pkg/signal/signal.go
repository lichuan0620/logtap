package signal

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
)

// SetupStopSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned which is closed on one
// of these signals. If a second signal is caught, the program is terminated with exit code 1.
func SetupStopSignalHandler() chan struct{} {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
