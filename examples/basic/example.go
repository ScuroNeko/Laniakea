package main

import (
	"log"

	"git.nix13.pw/scuroneko/laniakea"
)

func pong(ctx *laniakea.MsgContext, db *laniakea.NoDB) {
	ctx.Answer(ctx.Msg.Text)
}

func main() {
	bot := laniakea.NewBot[laniakea.NoDB](laniakea.LoadOptsFromEnv())
	defer bot.Close()

	p := laniakea.NewPlugin[laniakea.NoDB]("ping")
	p.NewCommand(pong, "ping")

	bot = bot.ErrorTemplate(
		"Error\n\n%s",
	).AddPlugins(p)

	if err := bot.AutoGenerateCommands(); err != nil {
		log.Println(err)
	}
	bot.Run()
}
