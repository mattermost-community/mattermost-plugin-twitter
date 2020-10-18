package util

import (
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-twitter/server/constant"
)

// Min - since math.Min is for floats and casting to and from floats is dangerous.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// SplitArgs is used to split a string to an array of arguments with separators: "(quotes) and spaces
// We cant use strings.split as it includes empty string if deliminator is the last character in input string
func SplitArgs(s string) ([]string, error) {
	indexes := regexp.MustCompile("\"").FindAllStringIndex(s, -1)
	if len(indexes)%2 != 0 {
		return []string{}, errors.New("quotes not closed")
	}

	indexes = append([][]int{{0, 0}}, indexes...)

	if indexes[len(indexes)-1][1] < len(s) {
		indexes = append(indexes, [][]int{{len(s), 0}}...)
	}

	var args []string
	for i := 0; i < len(indexes)-1; i++ {
		start := indexes[i][1]
		end := Min(len(s), indexes[i+1][0])

		if i%2 == 0 {
			args = append(args, strings.Split(strings.Trim(s[start:end], " "), " ")...)
		} else {
			args = append(args, s[start:end])
		}
	}

	cleanedArgs := make([]string, len(args))
	count := 0

	for _, arg := range args {
		if arg != "" {
			cleanedArgs[count] = strings.TrimSpace(arg)
			count++
		}
	}

	return cleanedArgs[0:count], nil
}

// SendEphemeralCommandResponse can be used to return an ephemeral message as the response for a slash command
func SendEphemeralCommandResponse(message string) (*model.CommandResponse, *model.AppError) {
	return &model.CommandResponse{
		Type:     model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:     message,
		Username: constant.BotUsername,
		IconURL:  constant.BotIconURL,
	}, nil
}

// BaseCommandResponse returns the base slash command response to post in channel
func BaseCommandResponse() *model.CommandResponse {
	return &model.CommandResponse{
		Type:     model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
		Username: constant.BotUsername,
		IconURL:  constant.BotIconURL,
	}
}
