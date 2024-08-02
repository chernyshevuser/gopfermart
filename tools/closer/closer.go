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

func GracefulShutdown(closers ...Closer) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-s
	closeAll(closers...)
}

func closeAll(closers ...Closer) {
	for i := range closers {
		if err := closers[i].Close(); err != nil {
			fmt.Printf("Error closing resource: %v\n", err)
		}
	}
	fmt.Println("Graceful shutdown!")
}
