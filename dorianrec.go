package main

import "github.com/hypebeast/go-osc/osc"

func main() {
	for {
		addr := "127.0.0.1:8765"
		server := &osc.Server{Addr: addr}

		server.Handle("/dorian/address", func(msg *osc.Message) {
			osc.PrintMessage(msg)
		})

		server.ListenAndServe()
	}
}


