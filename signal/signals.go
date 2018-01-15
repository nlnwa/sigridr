package signal

import (
	"os"
	"os/signal"
)

func Receive(sigs ...os.Signal) <-chan os.Signal {
	sigc := make(chan os.Signal, 1)

	signal.Notify(sigc, sigs...)

	return sigc
}
