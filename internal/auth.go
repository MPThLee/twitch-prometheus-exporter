package twitch_prometheus_exporter

import (
	"fmt"
	"io/ioutil"

	"github.com/kataras/golog"
	"github.com/nicklaw5/helix"
)

// totally bad design.

var token_file = config.Main.ConfigDir + "/access_token.secret"

func auth_logger(child string) *golog.Logger {
	return Log.Child("auth").Child(child)
}

func IsUserAuthorized(twitch *helix.Client) (bool, error) {
	var logger = auth_logger("isAuthroized")

	token := twitch.GetUserAccessToken()
	if token == "" {
		resp, err := loadUserTokenFromFile(twitch)
		if err != nil {
			return resp, err
		}
	}

	resp, err := validateUserToken(twitch, token)
	if err != nil {
		logger.Error("token is not valid. ", err)
		return resp, err
	}

	return true, nil
}

func validateAndSetUserToken(twitch *helix.Client, token string) (bool, error) {
	isValid, err := validateUserToken(twitch, token)
	if !isValid {
		return false, err
	}
	twitch.SetUserAccessToken(token)
	return setUserTokenToFile(token)
}

func validateUserToken(twitch *helix.Client, token string) (bool, error) {
	var logger = auth_logger("validateUserToken")

	isValid, resp, err := twitch.ValidateToken(token)
	if err != nil {
		logger.Error("failed to validate token", "resp", resp)
		return isValid, err
	}

	logger.Debug("resp", resp)
	return isValid, nil
}

func loadUserTokenFromFile(twitch *helix.Client) (bool, error) {
	var logger = auth_logger("loadTokenFromFile")

	buf, err := ioutil.ReadFile(token_file)
	if err != nil {
		logger.Info("Failed to retrieve token from file.")
	}

	token := string(buf)
	return validateAndSetUserToken(twitch, token)
}

func setUserTokenToFile(token string) (bool, error) {
	var logger = auth_logger("setTokenToFile")
	err := ioutil.WriteFile(token_file, []byte(token), 0644)
	if err != nil {
		logger.Error("failed to save token to file", err)
		return false, err
	}
	logger.Debug("write token file success.")
	return true, nil
}

func RefreshUserToken(twitch *helix.Client) (bool, error) {
	var logger = auth_logger("refreshToken")

	token := twitch.GetUserAccessToken()
	resp, err := validateUserToken(twitch, token)
	if !resp {
		logger.Error("failed to refresh token due to token not validate")
		return resp, err
	}

	response, err := twitch.RefreshUserAccessToken(token)
	if err != nil {
		logger.Error("failed to refresh token due to token not validate",
			err, "resp", response)
		return false, err
	}

	logger.Debug("resp", resp)
	return true, nil
}

func RequestAuthorize(twitch *helix.Client) (bool, error) {
	var logger = auth_logger("requestAuthorize")
	var authcode string

	url := twitch.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "code", // or "token"
		Scopes:       []string{"user:read:email"},
		State:        "some-state",
		ForceVerify:  false,
	})

	logger.Infof("Please visit this url to get authorized code: %s", url)

	fmt.Println("Enter Your Code: ")
	fmt.Scanln(&authcode)

	return requestUserToken(twitch, authcode)
}

func requestUserToken(twitch *helix.Client, code string) (bool, error) {
	var logger = auth_logger("requestUserToken")

	resp, err := twitch.RequestUserAccessToken(code)
	if err != nil {
		logger.Error("failed to request user token", err)
		return false, err
	}

	logger.Debug("resp", resp)

	return validateAndSetUserToken(twitch, resp.Data.AccessToken)
}

func RequestAppToken(twitch *helix.Client) (bool, error) {
	var logger = auth_logger("requestAppToken")

	resp, err := twitch.RequestAppAccessToken([]string{})
	if err != nil {
		logger.Error("failed to request app token", err)
		return false, err
	}

	logger.Debug("resp", resp)

	twitch.SetAppAccessToken(resp.Data.AccessToken)
	return true, nil
}
