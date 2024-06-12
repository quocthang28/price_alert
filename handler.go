package main

import (
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	AddSymbolCmd        = "add"
	RemoveSymbolCmd     = "remove"
	GetCurrentPricesCmd = "current"
	ChangeIntervalCmd   = "interval"
	ListSymbolsCmd      = "list"
	ShowInfoCmd         = "info"
)

func (app *App) registerCommands() {
	_, err := app.DiscordSession.ApplicationCommandBulkOverwrite(app.Config.AppID, app.Config.GuildID, []*discordgo.ApplicationCommand{
		{
			Name:        AddSymbolCmd,
			Description: "[symbol] - Add price alert for this symbol.",
			Options: []*discordgo.ApplicationCommandOption{
				{

					Name:        "token_symbol",
					Description: "The symbol of the token",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        RemoveSymbolCmd,
			Description: "[symbol] - Remove price alert for this symbol.",
			Options: []*discordgo.ApplicationCommandOption{
				{

					Name:        "token_symbol",
					Description: "The symbol of the token",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        GetCurrentPricesCmd,
			Description: "Get current prices",
			Options: []*discordgo.ApplicationCommandOption{
				{

					Name:        "token_symbol",
					Description: "The symbol of the token",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
			},
		},
		{
			Name:        ChangeIntervalCmd,
			Description: "[interval] - Change alert interval. Valid time units: \"m\", \"h\".",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "interval",
					Description: "The time between every price alert",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        ShowInfoCmd,
			Description: "Show current watchlist and alert interval.",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (app *App) addSymbolHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == AddSymbolCmd {
		symbol, err := validateCommandOption(i.ApplicationCommandData().Options)
		if err != nil {
			replyToCommand(s, i, err.Error())
			return
		}

		if slices.Contains(app.AlertConfig.getSymbols(), symbol) {
			replyToCommand(s, i, fmt.Sprintf("Symbol %s is already on watchlist!", symbol))
			return
		}

		err = app.AlertConfig.addSymbol(symbol)
		if err != nil {
			replyToCommand(s, i, err.Error())
		}

		// reschedule with updated symbols
		app.Scheduler.schedulePriceAlertJob(app.alertCryptoPrices, app.AlertConfig.getSymbols(), app.AlertConfig.getInterval(), true)

		replyToCommand(s, i, fmt.Sprintf("Added %s to watchlist.", symbol))
	}
}

func (app *App) removeSymbolHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == RemoveSymbolCmd {
		symbol, err := validateCommandOption(i.ApplicationCommandData().Options)
		if err != nil {
			replyToCommand(s, i, err.Error())
			return
		}

		if !slices.Contains(app.AlertConfig.getSymbols(), symbol) {
			replyToCommand(s, i, fmt.Sprintf("Symbol %s is not on watchlist!", symbol))
			return
		}

		err = app.AlertConfig.removeSymbol(symbol)
		if err != nil {
			replyToCommand(s, i, err.Error())
		}

		// reschedule with updated symbols
		app.Scheduler.schedulePriceAlertJob(app.alertCryptoPrices, app.AlertConfig.getSymbols(), app.AlertConfig.getInterval(), true)

		replyToCommand(s, i, fmt.Sprintf("Remove %s from watchlist.", symbol))
	}
}

func (app *App) getCurrentPricesHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == GetCurrentPricesCmd {
		symbol, err := validateCommandOption(i.ApplicationCommandData().Options)
		if err != nil {
			symsbols := app.AlertConfig.getSymbols()
			prices := getCryptoPrices(symsbols)
			replyToCommand(s, i, prices)
		} else {
			price := getCryptoPrices([]string{symbol})
			replyToCommand(s, i, price)
		}
	}
}

func (app *App) changeIntervalHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == ChangeIntervalCmd {
		interval, err := validateCommandOption(i.ApplicationCommandData().Options)
		if err != nil {
			replyToCommand(s, i, err.Error())
			return
		}

		re := regexp.MustCompile(`\b\d+[mh]\b`)
		if !re.MatchString(interval) {
			replyToCommand(s, i, "Interval is invalid!")
			return
		}

		err = app.AlertConfig.changeInterval(interval)
		if err != nil {
			replyToCommand(s, i, err.Error())
			return
		}

		// reschedule with updated interval
		app.Scheduler.schedulePriceAlertJob(app.alertCryptoPrices, app.AlertConfig.getSymbols(), interval, true)

		unit := "hour(s)"
		if strings.Contains(interval, "m") {
			unit = "minute(s)"
		}

		replyToCommand(s, i, fmt.Sprintf("Change applied, prices will be alerted every %s %s.", string(interval[:len(interval)-1]), unit))
	}
}

func (app *App) showInfoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == ShowInfoCmd {
		symbols := app.AlertConfig.getSymbols()
		interval := app.AlertConfig.getInterval()

		unit := "hour(s)"
		if strings.Contains(interval, "m") {
			unit = "minute(s)"
		}

		replyToCommand(s, i, fmt.Sprintf("Current watchlist: %s\nAlert interval: %s %s.", strings.Join(symbols, ", "), string(interval[:len(interval)-1]), unit))
	}
}

// func (app App) handleRequestLogFile() {

// }
