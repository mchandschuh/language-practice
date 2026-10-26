package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/mchandschuh/language-practice/go/temporal/hello"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal()
	}

	defer c.Close()

	w := worker.New(c, hello.TaskQueue, worker.Options{})
	w.RegisterWorkflow(hello.HelloWorkflow)
	w.RegisterActivity(hello.HelloActivity)

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Unable to start hello/worker", err)
	}
}
