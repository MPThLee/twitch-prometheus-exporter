package twitch_prometheus_exporter

import (
	"github.com/kataras/golog"
	"github.com/nicklaw5/helix"
	"github.com/prometheus/client_golang/prometheus"
)

func metrics_logger(child string) *golog.Logger {
	return Log.Child("metrics").Child(child)
}

func GetAppMetricsCollectors(twitch *helix.Client) []prometheus.Collector {
	var logger = metrics_logger("GetAppMetricsCollectors")
	errHandler := func(err error) {
		if err != nil {
			logger.Error(err)
			panic(err)
		}
	}

	followers, err := LoadFollowers(twitch)
	errHandler(err)

	followersCollector, err := NewFollowersCollector(twitch, followers)
	errHandler(err)
	streamDataCollector, err := NewStreamDataCollector(twitch, followers)
	errHandler(err)

	return []prometheus.Collector{
		followersCollector,
		streamDataCollector,
	}
}
