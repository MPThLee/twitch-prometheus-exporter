# Twitch Prometheus Exporter
For who want to track your favorate streamer.

# How to run
WIP \
Build first. Configure the config. run.

Note: Login feature is WIP. may not working collectly.

# Docker
WIP

Note: should run with `-it` first to get user token.

# TODO
* Refactor bad design
* Fix user access system
* Grafana dashboard template
* Collect IRC data for certain emotes or characters during stream.
* Collect user-based channel stats (like channel points)
  - This would be hard as they doesn't expose this API to helix. should use GraphQL API which doesn't documented.
  - [go-twitch](https://github.com/Adeithe/go-twitch) has (deprecated) GQL library.
