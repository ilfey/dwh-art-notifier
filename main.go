package main

import (
	"main/nekos"
	"main/webhook"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

func getEmbeds(category string, errChan chan error) []webhook.Embed {
	// Get art
	res, err := nekos.GetArt(category)
	if err != nil {
		errChan <- err
		return nil
	}

	embeds := make([]webhook.Embed, len(res.Results))

	// Fill embeds
	for i, result := range res.Results {
		embeds[i] = webhook.Embed{
			Title: &category,
			Url:   result.SourceURL,
			Image: &webhook.Image{
				Url: result.URL,
			},
			Author: &webhook.Author{
				Name: result.ArtistName,
				Url:  result.ArtistHref,
			},
		}
	}

	return embeds
}

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	viper.SetDefault("DELAY", "1h")
	viper.SetDefault("CATEGORY", "waifu")

	delay := viper.GetString("DELAY")
	category := viper.GetString("CATEGORY")
	wh_url := viper.GetString("WEBHOOK_URL")
	if wh_url == "" {
		panic("WEBHOOK_URL is not set")
	}

	errChan := make(chan error)
	runner := cron.New()

	// Add cron job
	_, err := runner.AddFunc("@every "+delay, func() {
		embeds := getEmbeds(category, errChan)

		msg := &webhook.Message{
			Embeds: &embeds,
		}

		// Send message
		err := webhook.SendMessage(wh_url, msg)
		if err != nil {
			errChan <- err
			return
		}
	})
	if err != nil {
		panic(err)
	}

	// Start first run
	embeds := getEmbeds(category, errChan)

	msg := &webhook.Message{
		Embeds: &embeds,
	}

	// Send message
	err = webhook.SendMessage(wh_url, msg)
	if err != nil {
		errChan <- err
		return
	}

	// Start runner
	runner.Start()

	for {
		err := <-errChan
		if err != nil {
			panic(err)
		}
	}
}
