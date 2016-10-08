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
	ClientID      string   `json:"clientid"`
	Token         string   `json:"token"`
	OwnerID       string   `json:"owner"`
	CommandPrefix string   `json:"prefix"`
	Plugins       []Plugin `json:"plugins"`
}

type Plugin struct {
	Functions []Func `json:"functions"`
}

type Func struct {
	Package string `json:"package"`
	Name    string `json:"name"`
	Alias   string `json:"alias"`
	Info    string `json:"info"`
}

var this *Info

func init() {

	f, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(f, &this)

}

var funcs = make(map[string]Func)

func main() {

	fmt.Printf("ID: %s\nTok: %s\nOwn: %s\nPre: %s\n\n======PLUGINS======\n\n", this.ClientID, this.Token, this.OwnerID, this.CommandPrefix)

	for _, p := range this.Plugins {
		for _, f := range p.Functions {

			fmt.Printf("Package: %s\n  Func: %s\n", f.Package, f.Name)
			if f.Alias != "" {
				fmt.Printf("    Alias: %q\n", f.Alias)
				funcs[f.Alias] = f
			} else {
				funcs[f.Name] = f
			}
		}
		fmt.Printf("\n")
	}

	client, _ := discordgo.New("Bot " + this.Token)

	err := client.Open()
	if err != nil {
		panic(err)
	}
	client.AddHandler(Listen)
	println("connected")

	<-make(chan struct{})
}

func Listen(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID != this.ClientID && strings.HasPrefix(m.Message.Content, this.CommandPrefix) {
		str := strings.TrimPrefix(m.Message.Content, this.CommandPrefix)
		new_str := strings.Split(str, " ")

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
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cli.Run(args, "gorram", "nano/plugins/"+fun.Package, fun.Name)
	outC := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()

	os.Stdout = old

	out := <-outC

	return out
}

func NanoHelp(args []string) string {
	if len(args) > 2 {
		if f, ok := funcs[args[1]]; ok {
			return f.Info
		}
		return "function not found"
	}
	return "add a function name to the end to get more info about it"
}
