package twitch_prometheus_exporter

import (
	"sync"

	"github.com/nicklaw5/helix"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	followerGaugeDesc = prometheus.NewDesc(
		"twitch_followers_total",
		"Number of Followers",
		[]string{"username"}, nil,
	)
)

type FollowerCollector struct {
	client    *helix.Client
	streamers map[string]helix.User // key is the id, value is the username
}

func NewFollowersCollector(client *helix.Client, users []helix.User) (*FollowerCollector, error) {
	ids := GetIdsFromHelixUsers(users)
	resp, err := client.GetUsers(&helix.UsersParams{
		IDs: ids,
	})
	if err != nil {
		return nil, err
	}

	streamers := make(map[string]helix.User, len(ids))
	for _, streamer := range resp.Data.Users {
		streamers[streamer.ID] = streamer
	}

	return &FollowerCollector{
		client:    client,
		streamers: streamers,
	}, nil
}

func (fc FollowerCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(fc, ch)
}

func (fc FollowerCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	wg.Add(len(fc.streamers))
	for _, streamer := range fc.streamers {
		go fc.fetchFollowerCount(streamer.ID, ch, &wg)
	}

	wg.Wait()
}

func (fc FollowerCollector) fetchFollowerCount(id string, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	var logger = Log.Child("metrics").Child("app_followers")
	defer wg.Done()

	out, err := fc.client.GetUsersFollows(&helix.UsersFollowsParams{
		ToID:  id,
		First: 1,
	})
	if err != nil {
		logger.Error(err)
	}

	logger.Debug(out)

	ch <- prometheus.MustNewConstMetric(
		followerGaugeDesc,
		prometheus.GaugeValue,
		float64(out.Data.Total),
		fc.streamers[id].Login,
	)
}
