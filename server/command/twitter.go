package command

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-twitter/server/serializers"
	"github.com/mattermost/mattermost-plugin-twitter/server/util"
)

const (
	invalidCommand = "Invalid command parameters. Please use `/twitter help` for more information."

	helpText = "###### Twitter - Slash Command Help\n\n" +
		"* `/twitter connect` - Connect to your twitter account.\n" +
		"* `/twitter disconnect` - Disconnect your twitter account.\n"
)

func GetCommand(iconData string) *model.Command {
	return &model.Command{
		Trigger:              "twitter",
		DisplayName:          "Twitter",
		AutoComplete:         true,
		AutoCompleteDesc:     "Available commands: connect, disconnect, help.",
		AutoCompleteHint:     "[command]",
		AutocompleteData:     getAutoCompleteData(),
		AutocompleteIconData: iconData,
	}
}

func getAutoCompleteData() *model.AutocompleteData {
	twitter := model.NewAutocompleteData("twitter", "[command]", "Available commands: connect, disconnect, help.")

	connect := model.NewAutocompleteData("connect", "", "Connect to your twitter account.")
	twitter.AddCommand(connect)

	disconnect := model.NewAutocompleteData("disconnect", "", "Disconnect your twitter account.")
	twitter.AddCommand(disconnect)

	help := model.NewAutocompleteData("help", "", "Show twitter slash command help")
	twitter.AddCommand(help)

	return twitter
}

var TwitterCommandHandler = Handler{
	handlers: map[string]HandlerFunc{
		"connect":    twitterConnect,
		"disconnect": twitterDisconnect,
		"help":       twitterHelpCommand,
	},
	defaultHandler: func(context *Context, args ...string) (*model.CommandResponse, *model.AppError) {
		return util.SendEphemeralCommandResponse(invalidCommand)
	},
}

func twitterConnect(ctx *Context, args ...string) (*model.CommandResponse, *model.AppError) {
	// If the user is already connected to twitter.
	if twUser, err := ctx.store.GetTwitterUser(ctx.UserId); err == nil && twUser != nil {
		return util.SendEphemeralCommandResponse(fmt.Sprintf("You are already connected as twitter user: %s.\nUse `/twitter disconnect` to disconnect your account.", twUser.Name+" (@"+twUser.Username+")"))
	}

	twitterOAuth1Config := util.GetTwitterOAuth1Config(ctx.api, ctx.manifest)
	token, secret, err := twitterOAuth1Config.RequestToken()
	if err != nil {
		ctx.api.LogError("Failed to connect.", "userID", ctx.UserId, "Error", err.Error())
		return util.SendEphemeralCommandResponse("Failed to connect to twitter. If the problem persists, contact your system administrator.")
	}

	err = ctx.store.StoreOneTimeSecretJSON(ctx.UserId, &serializers.OAuth1aTemporaryCredentials{Token: token, Secret: secret})
	if err != nil {
		ctx.api.LogError("Failed to connect.", "userID", ctx.UserId, "Error", err.Error())
		return util.SendEphemeralCommandResponse("Failed to connect to twitter. If the problem persists, contact your system administrator.")
	}

	authURL, err := twitterOAuth1Config.AuthorizationURL(token)
	if err != nil {
		ctx.api.LogError("Failed to connect.", "userID", ctx.UserId, "Error", err.Error())
		return util.SendEphemeralCommandResponse("Failed to connect to twitter. If the problem persists, contact your system administrator.")
	}

	return util.SendEphemeralCommandResponse(fmt.Sprintf("Click [here](%s) to connect to your twitter account.", authURL))
}

func twitterDisconnect(ctx *Context, args ...string) (*model.CommandResponse, *model.AppError) {
	if err := ctx.store.DeleteTwitterUser(ctx.UserId); err != nil {
		ctx.api.LogError("Failed to disconnect user.", "userID", ctx.UserId, "Error", err.Error())
		return util.SendEphemeralCommandResponse("Failed to disconnect your twitter account. If the problem persists, contact your system administrator.")
	}
	return util.SendEphemeralCommandResponse("Successfully disconnected from your twitter account.")
}

func twitterHelpCommand(_ *Context, _ ...string) (*model.CommandResponse, *model.AppError) {
	return util.SendEphemeralCommandResponse(helpText)
}
