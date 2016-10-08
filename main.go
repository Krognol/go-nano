package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"npf.io/gorram/cli"

	"github.com/bwmarrin/discordgo"
)

type Info struct {
	// Optional: good to have to make sure the sender of a message isn't the bot itself
	ClientID string `json:"clientid"`

	// Required if running a bot account for authorization
	Token string `json:"token"`

	// Optional: if you want to know who owns the bot
	OwnerID string `json:"owner"`

	// Required: so no annoyances happen when listening
	CommandPrefix string `json:"prefix"`

	Plugins []Plugin `json:"plugins"`
}

type Plugin struct {
	Functions []Func `json:"functions"`
}

type Func struct {
	// Required: the package name which holds the function
	Package string `json:"package"`

	// Required: the function name, has to start with a capital letter
	Name string `json:"name"`

	// Optional: an alias for the function
	Alias string `json:"alias"`

	// Optional: good to have if you want to send info about the function
	Info string `json:"info"`
}

// Gets assigned to at init
var this *Info

func init() {

	// Open the config.json file and unmarshal it to 'this', named for simplicity
	f, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(f, &this)

}

// Store the plugin functions in a map for easier use
var funcs = make(map[string]Func)

func main() {
	// Some info
	fmt.Printf("ID: %s\nTok: %s\nOwn: %s\nPre: %s\n\n======PLUGINS======\n\n", this.ClientID, this.Token, this.OwnerID, this.CommandPrefix)

	// Iterate through all the plugins and functions
	for _, p := range this.Plugins {
		for _, f := range p.Functions {

			fmt.Printf("Package: %s\n  Func: %s\n", f.Package, f.Name)
			// If the alias is not an empty string, add it to the funcs map with the alias as a key
			if f.Alias != "" {
				fmt.Printf("    Alias: %q\n", f.Alias)
				funcs[f.Alias] = f
			} else {
				funcs[f.Name] = f
			}
		}
		fmt.Printf("\n")
	}

	// Initialize a new discord client
	// the "Bot " part is required if you're running a bot account
	client, _ := discordgo.New("Bot " + this.Token)

	err := client.Open()
	if err != nil {
		panic(err)
	}
	// Add the message listener as an event handler
	client.AddHandler(Listen)
	println("connected")

	// Easy was to have the program run until ctrl+c is pressed
	<-make(chan struct{})
}

func Listen(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Check so the author of the message is not itself, and the message has the designated prefix
	if m.Author.ID != this.ClientID && strings.HasPrefix(m.Message.Content, this.CommandPrefix) {
		// Trim the prefix and leave everything else, split it, and put it into a string array
		str := strings.TrimPrefix(m.Message.Content, this.CommandPrefix)
		new_str := strings.Split(str, " ")

		// Check if the first argument matches anything
		if new_str[0] == "help" {
			s.ChannelMessageSend(m.ChannelID, NanoHelp(new_str))
		}

		if f, ok := funcs[new_str[0]]; ok {
			s.ChannelMessageSend(m.ChannelID, CallPlugin(f, new_str[1:]))
		}
	}
}

//I thank nate finch for making gorram so this could be possible.
func CallPlugin(fun Func, args []string) string {
	out := CopyFromStdout(fun, args)
	return out
}

func CopyFromStdout(fun Func, args []string) string {
	// Create a backup for Stdout
	old := os.Stdout

	// Initialize a new Pipe
	r, w, _ := os.Pipe()
	os.Stdout = w

	// We run gorram, which will find the plugin in the specified directory
	// and then run the function called with any arguments we supplied to it
	cli.Run(args, "gorram", "nano/plugins/"+fun.Package, fun.Name)
	outC := make(chan string)

	// Copy the buffer to 'buf' and send the buffer string to outC
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// Close the Pipe
	w.Close()

	// Assign the old Stdout to Stdout
	os.Stdout = old

	// Declare out and return it
	out := <-outC

	return out
}

// Sends the 'info' property of a plugin function
func NanoHelp(args []string) string {
	if len(args) > 2 {
		if f, ok := funcs[args[1]]; ok {
			return f.Info
		}
		return "function not found"
	}
	return "add a function name to the end to get more info about it"
}
