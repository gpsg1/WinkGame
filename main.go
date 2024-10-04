package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// TEST_SERVER : 1271088825560727592
// AZ Guild : 948807733199642645

var (
	BotToken       = flag.String("token", "", "")
	GuildID        = flag.String("guild", "1271088825560727592", "")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var session *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	session, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

// Commands
var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "wink",
			Description: "Start wink game",
		},
	}
)

func init() {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			switch i.ApplicationCommandData().Name {
			case "wink":
				selectUserHandler(s, i)
			}
		case discordgo.InteractionMessageComponent:
			switch i.MessageComponentData().CustomID {
			case "user_select_menu":
				handleSelectMenu(s, i)
			case "start_button":
				handleStartButton(s, i)
			case "check", "cancel":
				followUpHandler(s, i)
			}
		}
	})
}

func main() {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		for _, v := range registeredCommands {
			err := session.ApplicationCommandDelete(session.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Shutting down.")
}
