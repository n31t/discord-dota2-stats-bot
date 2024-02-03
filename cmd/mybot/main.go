package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/n31t/godiscord/pkg/models"

	"github.com/bwmarrin/discordgo"
)

const prefix string = "!neit"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	sess, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		args := strings.Split(m.Content, " ")
		if args[0] != prefix {
			return
		}

		if args[1] == "ping" {
			s.ChannelMessageSend(m.ChannelID, "pong")
		}
		if args[1] == "embed" {
			embed := discordgo.MessageEmbed{
				Title:       "Title",
				Description: "Description",
			}
			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		}
		if args[1] == "stratz" {
			if len(args) < 3 {
				s.ChannelMessageSend(m.ChannelID, "Please provide an ID.")
				return
			}

			id := args[2]
			url := fmt.Sprintf("https://api.stratz.com/api/v1/Player/%s", id)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatal(err)
			}

			req.Header.Set("Authorization", "Bearer "+os.Getenv("STRATZ_TOKEN"))
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error making the request.")
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error reading the response.")
				return
			}

			var player models.Player
			err = json.Unmarshal(body, &player)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error parsing the response.")
				return
			}

			rankName, ok := models.SeasonalRankMap[player.SteamAccount.SeasonRank]
			if !ok {
				rankName = "Unknown"
			}

			embed := &discordgo.MessageEmbed{
				Title: "Player Data",
				Fields: []*discordgo.MessageEmbedField{
					{Name: "ID", Value: strconv.Itoa(player.SteamAccount.ID), Inline: true},
					{Name: "Last Active Time", Value: player.SteamAccount.LastActiveTime, Inline: true},
					{Name: "Dota Plus Subscriber", Value: strconv.FormatBool(player.SteamAccount.IsDotaPlusSubscriber), Inline: true},
					{Name: "Smurf Flag", Value: strconv.Itoa(player.SteamAccount.SmurfFlag), Inline: true},
					{Name: "Match Count", Value: strconv.Itoa(player.MatchCount), Inline: true},
					{Name: "Win Count", Value: strconv.Itoa(player.WinCount), Inline: true},
					{Name: "Behavior Score", Value: strconv.Itoa(player.BehaviorScore), Inline: true},
					{Name: "Season Rank", Value: rankName, Inline: true},
				},
			}

			embed2 := &discordgo.MessageEmbed{
				Title:  "Rank history",
				Fields: []*discordgo.MessageEmbedField{},
			}
			for i, rank := range player.Ranks {
				if i >= 25 {
					break
				}

				rankName, ok := models.SeasonalRankMap[rank.Rank]
				if !ok {
					rankName = "Unknown"
				}

				embed2.Fields = append(embed2.Fields, &discordgo.MessageEmbedField{
					Name:   fmt.Sprintf("Rank History %d", i+1),
					Value:  fmt.Sprintf("Date Time: %s\nRank: %s", rank.AsOfDateTime, rankName),
					Inline: true,
				})
			}

			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			s.ChannelMessageSendEmbed(m.ChannelID, embed2)
		}
	})

	sess.Identify.Intents = discordgo.IntentsGuildMessages

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	fmt.Println("Bot is running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, os.Kill)
	<-sc
}
