# go-nano
golang implementation of nano

# Installation
First you have to get [gorram](https://github.com/natefinch/gorram) as the plugin system requires it.

`go get npf.io/gorram/cli`

Then get discordgo 

`go get github.com/bwmarrin/discordgo`

# Usage

in the config.json file you add the plugins, the discord information, the functions, and the name of the function. You can give the function an alias. You also specify which package it is in (required!).

Create a new folder in your `GOPATH/src` and name it `nano` in the new folder you add another folder called `plugins`, and in there you add the plugins you want. All plugins require a new folder, unless they belong to the same package. Check the `main.go` file for how things work (documentation to be added).
