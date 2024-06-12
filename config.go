package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"slices"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	BotPrefix      string `json:"BotPrefix"`
	AppID          string `json:"AppID"`
	GuildID        string `json:"GuildID"`
	BotID          string `json:"BotID"`
	AlertChannelId string
}

func loadAppConfig() Config {
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := Config{}

	err = json.Unmarshal(decrypt(os.Getenv("APP_CONFIG_K"), os.Getenv("APP_CONFIG")), &config)
	if err != nil {
		log.Fatal("Error parsing config file:", err.Error())
	}

	return config
}

type AlertConfig struct {
	Symbols  []string `json:"symbols"`
	Interval string   `json:"interval"`

	sync.RWMutex
}

func (ac *AlertConfig) getSymbols() []string {
	ac.RLock()
	defer ac.RUnlock()

	return ac.Symbols
}

func (ac *AlertConfig) getInterval() string {
	ac.RLock()
	defer ac.RUnlock()

	return ac.Interval
}

func loadAlertConfig() *AlertConfig {
	jsonFile, err := os.Open("/app/config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	b, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var alertConfig AlertConfig
	err = json.Unmarshal(b, &alertConfig)
	if err != nil {
		log.Fatal(err)
	}

	return &alertConfig
}

func (ac *AlertConfig) addSymbol(symbol string) error {
	//modify in-memory symbols
	ac.Lock()
	ac.Symbols = append(ac.Symbols, symbol)
	ac.Unlock()

	// modify config json file
	err := ac.rewriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func (ac *AlertConfig) removeSymbol(symbol string) error {
	//modify in-memory symbols
	ac.Lock()
	idx := slices.Index(ac.Symbols, symbol)
	ac.Symbols = slices.Delete(ac.Symbols, idx, idx+1)
	ac.Unlock()

	// modify config json file
	err := ac.rewriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func (ac *AlertConfig) changeInterval(interval string) error {
	//modify in-memory interval
	ac.Lock()
	ac.Interval = interval
	ac.Unlock()

	// modify config json file
	err := ac.rewriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func (ac *AlertConfig) rewriteConfig() error {
	data, err := json.Marshal(ac)
	if err != nil {
		return err
	}

	err = os.WriteFile("/app/config/config.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}
