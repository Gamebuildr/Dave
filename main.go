package main

import (
	"net/http"
	"os"
	"strconv"

	"fmt"

	"github.com/Gamebuildr/Dave/client"
	"github.com/Gamebuildr/Dave/pkg/config"
	"github.com/robfig/cron"
)

func main() {
	devMode, err := strconv.ParseBool(os.Getenv(config.DevMode))
	if err != nil {
		devMode = false
	}

	c := cron.New()
	daveClient := client.DaveClient{}
	daveClient.DevMode = devMode
	daveClient.Create()

	c.AddFunc("0 * * * * *", func() {
		daveClient.RunClient(client.GogetaSubsystem, os.Getenv(config.GogetaSQSEndpoint), 10)
		daveClient.RunClient(client.MrRobotSubsystem, os.Getenv(config.MrrobotSQSEndpoint), 3)
	})
	c.Start()

	daveClient.Log.Info("Dave client running on port 3001.")
	fmt.Printf("Dave client running on port 3001")
	if err := http.ListenAndServe(":3001", nil); err != nil {
		daveClient.Log.Error(err.Error())
	}
}
