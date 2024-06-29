package closer

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Closer interface {
	Close() error
}

func GracefulShutdown(c chan struct{}, closers ...Closer) {
	s := make(chan os.Signal, 1)

	signal.Notify(s, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-s

	closeAll(closers...)
	close(c)
}

func closeAll(closers ...Closer) {
	for i := range closers {
		closers[i].Close()
	}
	fmt.Println("Graceful shutdown!")
}
