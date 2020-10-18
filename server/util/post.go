package util

import (
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-twitter/server/constant"
)

func EphemeralPost(channelID, message string) *model.Post {
	post := &model.Post{
		ChannelId: channelID,
		Message:   message,
	}
	post.SetProps(model.StringInterface{
		"from_webhook":      "true",
		"override_username": constant.BotUsername,
		"override_icon_url": constant.BotIconURL,
	})
	return post
}
