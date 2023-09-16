package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func main() {
	serverUrl := flag.String("url", "", "the url of the TCP server")
	lPort := flag.Int("port", 9090, "Listen port")
	debug := flag.Bool("debug", false, "Show debug logs")

	flag.Parse()

	log.SetLevel(log.InfoLevel)
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	tttApp := GameApp{}
	tttApp.Initialize()

	if *serverUrl == "" {
		log.WithField("func", "main").Infoln("Starting TCP server...")

		tttApp.SetSign(XSign)   // set the current player's sign to X
		tttApp.SetOpSign(OSign) // set opponent's sign to O

		if *lPort != 0 {
			go tttApp.StartServer(":" + strconv.Itoa(*lPort))
		} else {
			go tttApp.StartServer(":9090")
		}
	} else {
		log.WithField("func", "main").Infoln("Connecting to ", *serverUrl)
		go tttApp.ConnectToServer(*serverUrl)
	}

	tttApp.Show()
}
