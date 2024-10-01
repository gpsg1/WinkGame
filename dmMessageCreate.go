package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func sendPrivateMessage(s *discordgo.Session, userID string, message string) {
	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		fmt.Println("error creating channel:", err)
		return
	}

	_, err = s.ChannelMessageSend(channel.ID, message)
	if err != nil {
		fmt.Println("error sending DM message:", err)
	}
}
