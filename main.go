package main

import (
	"net/http"
	"os"

	"fmt"

	"github.com/Gamebuildr/Dave/client"
	"github.com/Gamebuildr/Dave/pkg/config"
	"github.com/Gamebuildr/Dave/pkg/scaler"
	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	daveClient := client.DaveClient{}
	daveClient.Create()

	gogetaScaler := createGogetaScaler()
	mrrobotScaler := createMrRobotScaler()

	c.AddFunc("0 * * * * *", func() {
		daveClient.RunClient(gogetaScaler, os.Getenv(config.GogetaSQSEndpoint))
		daveClient.RunClient(mrrobotScaler, os.Getenv(config.MrrobotSQSEndpoint))
	})
	c.Start()

	daveClient.Log.Info("Dave client running on port 3001.")
	fmt.Printf("Dave client running on port 3001")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		daveClient.Log.Error(err.Error())
	}
}

func createGogetaScaler() *scaler.ScalableSystem {
	loadAPI := os.Getenv(config.HalGogetaAPI) + "api/container/count"
	addLoadAPI := os.Getenv(config.HalGogetaAPI) + "api/container/run"

	gogetaScaler := scaler.HTTPScaler{
		LoadAPIUrl:    loadAPI,
		AddLoadAPIUrl: addLoadAPI,
		Client:        &http.Client{},
	}
	system := scaler.ScalableSystem{
		System:  gogetaScaler,
		MaxLoad: 10,
	}
	return &system
}

func createMrRobotScaler() *scaler.ScalableSystem {
	loadAPI := os.Getenv(config.HalMrRobotAPI) + "api/container/count"
	addLoadAPI := os.Getenv(config.HalMrRobotAPI) + "api/container/run"

	mrrobotScaler := scaler.HTTPScaler{
		LoadAPIUrl:    loadAPI,
		AddLoadAPIUrl: addLoadAPI,
		Client:        &http.Client{},
	}
	system := scaler.ScalableSystem{
		System:  mrrobotScaler,
		MaxLoad: 3,
	}
	return &system
}
