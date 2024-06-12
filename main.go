package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := newApp()

	log.Println("Creating websocket connection to Discord...")
	err := app.DiscordSession.Open()
	if err != nil {
		log.Fatal("Error opening connection, ", err)
	}

	log.Println("Bot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Shut down signal received, stopping...")
	app.DiscordSession.Close()
	log.Println("Bot is now offline.")
}
