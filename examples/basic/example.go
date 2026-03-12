package main

import (
	"log"

	"git.nix13.pw/scuroneko/laniakea"
)

func echo(ctx *laniakea.MsgContext, db *laniakea.NoDB) {
	ctx.Answer(ctx.Text) // User input WITHOUT command
}

func main() {
	opts := &laniakea.BotOpts{Token: "TOKEN"}
	bot := laniakea.NewBot[laniakea.NoDB](opts)
	defer bot.Close()

	p := laniakea.NewPlugin[laniakea.NoDB]("ping")
	p.AddCommand(p.NewCommand(echo, "echo"))
	p.AddCommand(p.NewCommand(func(ctx *laniakea.MsgContext, db *laniakea.NoDB) {
		ctx.Answer("Pong")
	}, "ping"))

	bot = bot.ErrorTemplate("Error\n\n%s").AddPlugins(p)

	if err := bot.AutoGenerateCommands(); err != nil {
		log.Println(err)
	}
	bot.Run()
}
