package twitch_prometheus_exporter

import (
	"github.com/nicklaw5/helix"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	viewersGaugeDesc = prometheus.NewDesc(
		"twitch_viewers_total",
		"Number of Viewers",
		[]string{"username"}, nil,
	)

	startedAtGaugeDesc = prometheus.NewDesc(
		"twitch_started_at",
		"Stream started at",
		[]string{"username"}, nil,
	)

	streamOnlineGaugeDesc = prometheus.NewDesc(
		"twitch_stream_online",
		"Is Stream Online",
		[]string{"username"}, nil,
	)
)

type StreamDataCollector struct {
	client *helix.Client
	users  []helix.User
}

func NewStreamDataCollector(client *helix.Client, users []helix.User) (*StreamDataCollector, error) {
	return &StreamDataCollector{
		client: client,
		users:  users,
	}, nil
}

func (vc StreamDataCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(vc, ch)
}

func (vc StreamDataCollector) Collect(ch chan<- prometheus.Metric) {
	var logger = Log.Child("metrics").Child("app_StreamData")

	ids := GetIdsFromHelixUsers(vc.users)
	out, err := vc.client.GetStreams(&helix.StreamsParams{
		UserIDs: ids,
	})
	if err != nil {
		logger.Error(err)
	}

	streamMap := make(map[string]helix.Stream)
	for _, stream := range out.Data.Streams {
		streamMap[stream.UserID] = stream
	}

	userMap := make(map[string]string)
	for _, user := range vc.users {
		userMap[user.ID] = user.Login
	}

	logger.Debug(out)
	for id, name := range userMap {
		online, viewers, uptime := false, 0, 0
		data, ok := streamMap[id]
		if ok {
			online = true
			viewers = data.ViewerCount
			uptime = int(data.StartedAt.Unix())
		}

		ch <- prometheus.MustNewConstMetric(
			streamOnlineGaugeDesc,
			prometheus.GaugeValue,
			boolToFloat64(online),
			name,
		)

		ch <- prometheus.MustNewConstMetric(
			viewersGaugeDesc,
			prometheus.GaugeValue,
			float64(viewers),
			name,
		)

		ch <- prometheus.MustNewConstMetric(
			startedAtGaugeDesc,
			prometheus.GaugeValue,
			float64(uptime),
			name,
		)
	}
}

// no direct way to convert...
func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
