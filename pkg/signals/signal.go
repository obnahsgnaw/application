package signals

import (
	"os"
	"os/signal"
	"syscall"
)

var done = make(chan bool, 1)

// Listen Listen os signal, default:syscall.SIGINT  syscall.SIGTERM
func Listen(cb func(), sig ...os.Signal) {
	if len(sig) == 0 {
		sig = append(sig, syscall.SIGINT)
		sig = append(sig, syscall.SIGTERM)
	}
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, sig...)
	go func(signalCh chan os.Signal, done chan bool) {
		<-signalCh
		done <- true
		if cb != nil {
			cb()
		}
	}(signalCh, done)
}

func Wait() {
	<-done
}
