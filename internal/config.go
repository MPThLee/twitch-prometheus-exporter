package twitch_prometheus_exporter

import (
	"flag"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

var config RootConfig

type RootConfig struct {
	API      ApiConfig      `koanf:"api"`
	Login    LoginConfig    `koanf:"login"`
	Main     MainConfig     `koanf:"main"`
	Streamer StreamerConfig `koanf:"streamer"`
	Scrape   ScrapeConfig   `koanf:"scrape"`
}

type ApiConfig struct {
	ClientID     string `koanf:"client_id"`
	ClientSecret string `koanf:"client_secret"`
	RedirectUrl  string `koanf:"redirect_url"`
}

type LoginConfig struct {
	Enabled bool `koanf:"enabled"`
}

type MainConfig struct {
	LogLevel   string `koanf:"log_level"`
	ListenPort int    `koanf:"listen_port"`
}

type StreamerConfig struct {
	LoadFollowers bool     `koanf:"load_followers"`
	Lists         []string `koanf:"list"`
}

type ScrapeConfig struct {
	Viewers bool `koanf:"viewers"`
}

func LoadConfig(path string) RootConfig {
	var logger = Log.Child("config").Child("loadConfig")
	var k = koanf.New(".")

	k.Load(structs.Provider(RootConfig{}, "koanf"), nil)

	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		logger.Fatalf("error loading config: ", err)
		panic(err)
	}

	var out RootConfig
	k.Unmarshal("", &out)
	return out
}

// Bad design.
func InitConfig() RootConfig {
	var configPath = flag.String("config", "./config.yaml", "Config path")
	flag.Parse()
	config = LoadConfig(*configPath)
	return config
}
