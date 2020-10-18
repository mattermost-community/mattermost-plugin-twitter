package constant

const (
	// TODO: use manifest.id instead
	PluginName = "com.mattermost.twitter"

	URLPluginBase = "/plugins/" + PluginName
	URLStaticBase = URLPluginBase + "/static"

	HeaderMattermostUserID = "Mattermost-User-Id"

	BotUsername    = "twitter"
	BotDisplayName = "Twitter"
	BotIconURL     = URLStaticBase + "/twitter.png"

	PathTwitterOAuth1Callback = "/twitter/callback"
)
