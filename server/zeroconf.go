package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/grandcat/zeroconf"
)

// TODO: Re-register service
func RegisterService(name string, port int, protocol string, ch chan bool, uuid ...string) {
	if len(uuid) == 0 {
		name = fmt.Sprintf("%s-%s", name, GenerateID())
	}
	if len(uuid) == 1 {
		name = fmt.Sprintf("%s-%s", name, uuid[0])
	}

	meta := []string{
		"txtv=0",
		"lo=1",
		"la=2",
		"id=" + strings.Split(name, "-")[1],
		"fn=" + name,
	}
	server, err := zeroconf.Register(
		name,
		protocol, //"_karaks._tcp",
		"local.",
		port,
		meta,
		nil,
	)

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			<-ch
			fmt.Println(">>> restart zeroconf ...")
			server.Shutdown()
			server, _ = zeroconf.Register(
				name,
				protocol, //"_karaks._tcp",
				"local.",
				port,
				meta,
				nil,
			)
		}
	}()

	defer server.Shutdown()

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
		server.Shutdown()
		log.Println("Exit by user")
		os.Exit(0)
		// Exit by user
	}

	log.Println("Shutting down.")
}

func NewZeroConfChannel() chan bool {
	return make(chan bool)
}
