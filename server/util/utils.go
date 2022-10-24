package util

import (
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	twAuth "github.com/dghubble/oauth1/twitter"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-twitter/server/config"
	"github.com/mattermost/mattermost-plugin-twitter/server/constant"
)

func GetSiteURL(api plugin.API) string {
	return *api.GetConfig().ServiceSettings.SiteURL
}

func GetPluginURLPath(manifest *model.Manifest) string {
	return "/plugins/" + manifest.Id
}

func GetPluginURL(api plugin.API, manifest *model.Manifest) string {
	return strings.TrimRight(GetSiteURL(api), "/") + GetPluginURLPath(manifest)
}

func GetPluginAPIURL(api plugin.API, manifest *model.Manifest) string {
	return GetPluginURL(api, manifest) + "/api/v1"
}

func GetTwitterOAuth1Config(api plugin.API, manifest *model.Manifest) oauth1.Config {
	conf := config.GetConfig()

	return oauth1.Config{
		ConsumerKey:    conf.OAuthClientID,
		ConsumerSecret: conf.OAuthClientSecret,
		CallbackURL:    GetPluginAPIURL(api, manifest) + constant.PathTwitterOAuth1Callback,
		Endpoint:       twAuth.AuthorizeEndpoint,
	}
}

func GetTwitterClient(accessToken, accessSecret string) *twitter.Client {
	conf := config.GetConfig()
	oauth1Config := oauth1.NewConfig(conf.OAuthClientID, conf.OAuthClientSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := oauth1Config.Client(oauth1.NoContext, token)

	// Twitter client
	return twitter.NewClient(httpClient)
}
