package twitch_prometheus_exporter

import (
	"errors"

	"github.com/nicklaw5/helix"
)

var my_userid = ""

func IsLoggedIn(twitch *helix.Client) bool {
	if config.Login.Enabled {
		resp, err := IsUserAuthorized(twitch)
		if err != nil {
			return false
		}
		return resp
	}
	return config.Login.Enabled
}

//TODO: use this for some user scope request
//nolint: deadcode
func getMyData(twitch *helix.Client) (*helix.User, error) {
	if !IsLoggedIn(twitch) {
		return nil, errors.New("Not logged in.")
	}

	resp, err := twitch.GetUsers(&helix.UsersParams{})
	if err != nil {
		return nil, err
	}

	data := &resp.Data.Users[0]
	my_userid = data.ID
	return data, nil
}

func loadLocalFollows(twitch *helix.Client) ([]helix.User, error) {
	var logger = Log.Child("loadLocalFollow")
	logger.Info("followers=", config.Streamer.Lists, config.API.ClientID)
	users_resp, err := twitch.GetUsers(&helix.UsersParams{
		Logins: config.Streamer.Lists,
	})
	if err != nil {
		return nil, err
	}

	return users_resp.Data.Users, nil
}

func loadRemoteFollowers(twitch *helix.Client) ([]helix.User, error) {
	if !IsLoggedIn(twitch) {
		return []helix.User{}, nil
	}
	resp, err := twitch.GetUsersFollows(&helix.UsersFollowsParams{FromID: my_userid})
	if err != nil {
		return []helix.User{}, err
	}
	follows := resp.Data.Follows
	var follow_ids = []string{}

	for i, val := range follows {
		follow_ids[i] = val.ToID
	}

	users_resp, err := twitch.GetUsers(&helix.UsersParams{
		IDs: follow_ids,
	})
	if err != nil {
		return nil, err
	}

	return users_resp.Data.Users, nil
}

func followsFilter(arr []helix.User) []helix.User {
	var users = []helix.User{}
	for _, val := range arr {
		if !contains(users, val) {
			users = append(users, val)
		}
	}

	return users
}

func contains(s []helix.User, x helix.User) bool {
	for _, v := range s {
		if v.ID == x.ID {
			return true
		}
	}
	return false
}

func LoadFollowers(twitch *helix.Client) ([]helix.User, error) {
	local, err := loadLocalFollows(twitch)
	if err != nil {
		return nil, err
	}

	remote, err := loadRemoteFollowers(twitch)
	if err != nil {
		return nil, err
	}

	var out = []helix.User{}

	out = append(out, local...)
	out = append(out, remote...)

	return followsFilter(out), nil
}

func GetIdsFromHelixUsers(users []helix.User) []string {
	var follow_ids = make([]string, len(users))

	for i, val := range users {
		follow_ids[i] = val.ID
	}

	return follow_ids
}
func GetUserNameFromIds(users []helix.User) []string {
	var follow_ids = make([]string, len(users))

	for i, val := range users {
		follow_ids[i] = val.ID
	}

	return follow_ids
}
