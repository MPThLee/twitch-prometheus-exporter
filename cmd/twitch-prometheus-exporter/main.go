package main

import (
	"fmt"
	"net/http"
	"time"

	internal "github.com/mpthlee/twitch-prometheus-exporter/internal"

	"github.com/go-co-op/gocron"
	"github.com/nicklaw5/helix"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := internal.Log.Child("main")
	config := internal.InitConfig()

	internal.Log.SetLevel(config.Main.LogLevel)
	c := gocron.NewScheduler(time.UTC)

	client, err := helix.NewClient(&helix.Options{
		ClientID:     config.API.ClientID,
		ClientSecret: config.API.ClientSecret,
	})
	if err != nil {
		panic(err)
	}

	internal.RequestAppToken(client)

	if config.Login.Enabled {
		res, err := internal.RequestAuthorize(client)
		if err != nil || res == false {
			panic(err)
		}

		internal.RefreshUserToken(client)
		c.Every(1).Days().Do(func() {
			internal.RefreshUserToken(client)
		})
	}

	appCollectors := internal.GetAppMetricsCollectors(client)
	for _, val := range appCollectors {
		prometheus.MustRegister(val)
	}

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{},
	))
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Main.ListenPort), nil))
}
