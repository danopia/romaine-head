package head

import (
	"log"
	"os"
	"os/signal"
)

func WaitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan

		ShutdownLeaves()
		os.Exit(0)
	}()
}

func ShutdownLeaves() {
	log.Println("Shutting down leaves")

	for _, leaf := range leaves {
		if leaf != nil && (leaf.State == "launching" || leaf.State == "running") {
			leaf.Anchor.Process.Kill()
		}
	}
}
