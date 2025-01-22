package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func parseAmount(s string) (int, int, error) {
	sp := strings.Split(s, "@")
	if len(sp) == 2 {
		i, err := strconv.Atoi(sp[0])
		if err != nil {
			return 64, 0, fmt.Errorf("strconv error: sp0")
		}

		i2, err := strconv.Atoi(sp[1])
		if err != nil {
			return 64, 0, fmt.Errorf("strconv error: sp1")
		}

		return i2, i, nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 64, 0, fmt.Errorf("strconv error")
	}
	return 64, i, nil
}

func mod(x int, y int) int {
	return int(math.Mod(float64(x), float64(y)))
}

func calc(amount int, size int) string {
	lc := amount / (54 * size)
	amount = mod(amount, 54*size)
	sb := amount / (27 * size)
	amount = mod(amount, 27*size)
	st := amount / size
	amount = mod(amount, size)

	var res []string
	if lc > 0 {
		res = append(res, fmt.Sprintf("%dLC", lc))
	}
	if sb > 0 {
		res = append(res, fmt.Sprintf("%dc", sb))
	}
	if st > 0 {
		res = append(res, fmt.Sprintf("%dst", st))
	}
	if amount > 0 {
		res = append(res, fmt.Sprintf("%d", amount))
	}

	return strings.Join(res, "+")
}

func onMessage(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.Bot {
		return
	}

	if strings.HasSuffix(msg.Content, "?=") {
		a := strings.TrimSuffix(msg.Content, "?=")
		size, amount, err := parseAmount(a)
		if err == nil {
			s.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
				AllowedMentions: &discordgo.MessageAllowedMentions{Parse: []discordgo.AllowedMentionType{}},
				Reference:       msg.Reference(),
				Content:         calc(amount, size),
			})
		} else {
			s.MessageReactionAdd(msg.ChannelID, msg.ID, "❌")
		}
	}
}

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	discord.Identify.Intents = discordgo.IntentMessageContent | discordgo.IntentGuildMessages
	discord.AddHandler(onMessage)

	err = discord.Open()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Bot is now running.")
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	<-sigch

	err = discord.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}
