package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token") // Take in a flag and store the value at token reference
	flag.Parse()
}

var token string
var buffer = make([][]byte, 0)
var commandFound bool = false

func main() {

	if token == "" {
		fmt.Println("No token provided. Please run soundboard -t <bot token>")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register  ready as a callback for the ready events
	dg.AddHandler(ready)

	// Register messageCreate as a callback for the messageCreate events
	dg.AddHandler(messageCreate)

	// Register guildCreate as a callback for the guildCreate events
	dg.AddHandler(guildCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received
	fmt.Println("Soundboard is now running. Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)                                             // create a channel that can only take in os.Signal types with a length of 1
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill) // relays package signals to SC channel
	<-sc                                                                      //Blocks code execution until a signal is received

	// Cleanly close down disco session
	dg.Close()

}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the game played status message
	s.UpdateStatus(0, "!soundboard")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This prevent infinite looping
	if m.Author.ID == s.State.User.ID {
		return // Do nothing because its the bot message
	}

	// check if the message starts with "!"
	if strings.HasPrefix(m.Content, "!") {
		commandFound = true

		// Find channel message derived from
		c, err := s.State.Channel(m.ChannelID) // Returns a pointer to a state object and an error
		if err != nil {
			// can't find channel
			return
		}

		// Find the guild for that channel
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			// cannot find guild
			return
		}

		// TODO: Add load sound here
		switch strings.ToLower(m.Content) {
		case "!airhorn":
			loadSound("airhorn")
		case "!truckstop":
			loadSound("truckstop")
		case "!birthdaysad":
			loadSound("birthday_sadhorn")
		case "!airhorn4":
			loadSound("airhorn_fourtap")
		case "!airhorntruck":
			loadSound("airhorn_truck")
		case "!airhornclown":
			loadSound("airhorn_clownfull")
		case "!airhornfart":
			loadSound("airhorn_highfartlong")
		case "!birthday":
			loadSound("birthday_horn3")
		case "!moo":
			loadSound("cow_moo")
		case "!cows":
			loadSound("cow_herd")
		case "!crazy":
			loadSound("ethan_areyou_classic")
		case "!soda":
			loadSound("ethan_sodiepop")
		case "!ethanclassic":
			loadSound("ethan_classic")
		case "!ethancut":
			loadSound("ethan_cuts")
		case "!jc":
			loadSound("jc_full")
		case "!anotherone":
			loadSound("another_one_classic")
		case "!realfast":
			loadSound("realfast")
		case "!spaghetti":
			loadSound("spaghetti")
		case "!ohmygoodness":
			loadSound("ohmygoodness")
		case "!tellusthetruth":
			loadSound("tellusthetruth")
		case "!american":
			loadSound("american")
		case "!assuming":
			loadSound("assuming")
		case "!bummer":
			loadSound("bummer")
		case "!crack":
			loadSound("crack")
		case "!ding":
			loadSound("ding")
		case "!erik":
			loadSound("erik")
		case "!fantasy":
			loadSound("fantasy")
		case "!jiras":
			loadSound("jiras")
		case "!lossofwords":
			loadSound("lossofwords")
		case "!niceguy":
			loadSound("niceguy")
		case "!nonono":
			loadSound("nonono")
		case "!ooo":
			loadSound("ooo")
		case "!retrograde":
			loadSound("retrograde")
		case "!sweetjesus":
			loadSound("sweetjesus")
		case "!talktome":
			loadSound("talktome")
		case "!whileitlasted":
			loadSound("whileitlasted")
		case "!butthole":
			loadSound("butthole")
		case "!help":
			commandFound = false
			s.ChannelMessageSend(m.ChannelID, "I found the following sounds in the database:\n\n!truckstop, !airhorn, !birthdaysad\n!airhorn4, !airhornclown, !airhornfart\n!birthday, !moo, !cows\n!crazy, !soda, !ethanclassic\n !ethancut, !jc, !anotherone")
			s.ChannelMessageSend(m.ChannelID, "\n!spaghetti, !ohmygoodness, !tellusthetruth, !american\n!assuming, !bummer, !crack, !ding\n!erik, !fantasy, !jiras, !lossofwords\n!niceguy, !nonono, !ooo, !retrograde\n!sweetjesus, !talktome, !whileitlasted\n!butthole")
		default:
			commandFound = false
			s.ChannelMessageSend(m.ChannelID, "Sound not found. Please type in !help for a list of available sounds")

		}

		// Look for message sender in that guilds current voice states
		// TODO: replace test with _
		if commandFound {
			for _, vs := range g.VoiceStates {
				if vs.UserID == m.Author.ID {
					err = playSound(s, g.ID, vs.ChannelID)
					if err != nil {
						fmt.Println("Error playing sound: ", err)
					}
					return
				}
			}
		}
	}
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Soundboard is now ready! Type in !help to see a list of available sounds")
			return
		}
	}
}

func loadSound(sound string) error {
	file, err := os.Open(sound + ".dca")
	if err != nil {
		fmt.Println("Error opening dca file: ", err)
	}

	var opuslen int16

	for {
		//read opus frame length from DCA file
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// if this is the end of the file, just return
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from DCA file: ", err)
			return err
		}

		// read encoded PCM from DCA file
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Shouldn't be any EOF errors
		if err != nil {
			fmt.Println("Error reading from DCA file: ", err)
			return err
		}

		// append encoded pcm data to the buffer
		buffer = append(buffer, InBuf)
	}
}

// playSound plays the current buffer to the provided channel
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {
	// join the provided voice channel user who intiated the command is in
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep a short period of time before playing sound
	time.Sleep(250 * time.Millisecond)

	// start speaking
	vc.Speaking(true)

	// send buffer data
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// stop speaking
	vc.Speaking(false)

	// sleep again
	time.Sleep(250 * time.Millisecond)

	// disconnect
	vc.Disconnect()
	commandFound = false
	buffer = make([][]byte, 0)

	return nil
}
