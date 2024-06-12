package main

import (
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type App struct {
	Config         Config
	AlertConfig    *AlertConfig
	DiscordSession *discordgo.Session
	Scheduler      *AppScheduler
	//Firestore      *firestore.Client
}

func newApp() *App {
	app := App{
		Config:      loadAppConfig(),
		AlertConfig: loadAlertConfig(),
	}

	ds, err := discordgo.New("Bot " + string(decrypt(os.Getenv("BOT_TOKEN_K"), os.Getenv("BOT_TOKEN"))))
	if err != nil {
		log.Fatal("Error creating Discord session, ", err)
	}

	app.DiscordSession = ds
	app.DiscordSession.Identify.Intents = discordgo.IntentsGuildMessages
	app.DiscordSession.AddHandler(app.addSymbolHandler)
	app.DiscordSession.AddHandler(app.removeSymbolHandler)
	app.DiscordSession.AddHandler(app.getCurrentPricesHandler)
	app.DiscordSession.AddHandler(app.changeIntervalHandler)
	app.DiscordSession.AddHandler(app.showInfoHandler)
	app.registerCommands()

	// get price-alert channel id
	channels, err := app.DiscordSession.GuildChannels(app.Config.GuildID)
	if err != nil {
		log.Fatal("Error getting channels, ", err)
	}

	for _, channel := range channels {
		if strings.Contains(channel.Name, "price-alert") {
			app.Config.AlertChannelId = channel.ID
		}
	}

	// init scheduler
	app.Scheduler = newAppScheduler()
	app.Scheduler.schedulePriceAlertJob(app.alertCryptoPrices, app.AlertConfig.getSymbols(), app.AlertConfig.getInterval(), false)

	// app.Firestore, err = NewFirebaseClient()
	// if err != nil {
	// 	log.Fatal("Error creating Firestore client, ", err)
	// }

	return &app
}

func (app App) alertCryptoPrices(symbols []string) {
	_, err := app.DiscordSession.ChannelMessageSend(app.Config.AlertChannelId, getCryptoPrices(symbols))
	if err != nil {
		log.Println(err)
	}
}
