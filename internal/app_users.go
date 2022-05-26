package twitch_prometheus_exporter

import (
	"github.com/nicklaw5/helix"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	viewsCountDesc = prometheus.NewDesc(
		"twitch_views_total",
		"Number of Views",
		[]string{"username"}, nil,
	)
)

type UsersCollector struct {
	client *helix.Client
	users  []helix.User
}

func NewUsersCollector(client *helix.Client, users []helix.User) (*UsersCollector, error) {
	return &UsersCollector{
		client: client,
		users:  users,
	}, nil
}

func (vc UsersCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(vc, ch)
}

func (vc UsersCollector) Collect(ch chan<- prometheus.Metric) {
	var logger = Log.Child("metrics").Child("app_users")

	ids := GetIdsFromHelixUsers(vc.users)
	out, err := vc.client.GetUsers(&helix.UsersParams{
		IDs: ids,
	})
	if err != nil {
		logger.Error(err)
	}

	logger.Debug(out)
	for _, user := range out.Data.Users {
		ch <- prometheus.MustNewConstMetric(
			viewsCountDesc,
			prometheus.CounterValue,
			float64(user.ViewCount),
			user.Login,
		)
	}
}
