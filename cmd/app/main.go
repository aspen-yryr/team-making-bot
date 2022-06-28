package main

import (
	"flag"
	"os"
	"team-making-bot/internal/bot"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
)

const defaultEnvFile = "./env/.env"

func main() {
	e := flag.String("env_file", defaultEnvFile, "env variables definition file")
	greet := flag.Bool("greet", false, "If true, bot greet to channel on activated")
	flag.Parse()

	err := godotenv.Load(*e)
	if err != nil {
		glog.Error("Cannot load \"" + *e + "\" file")
		return
	}
	apiKey := os.Getenv("DISCORD_BOT_KEY")
	if apiKey == "" {
		glog.Error("Cannot get bot Key")
		return
	}
	b := bot.New(apiKey, *greet)
	b.Run()
}
