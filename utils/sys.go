package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func MakeSysChan() chan os.Signal {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	return sigCh
}
