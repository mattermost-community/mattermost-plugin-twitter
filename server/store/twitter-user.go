package store

import (
	"github.com/mattermost/mattermost-plugin-twitter/server/serializers"
	"github.com/mattermost/mattermost-plugin-twitter/server/util"
)

const (
	twitterUserPrefix = "twitter-user-"
)

func (s Store) SaveTwitterUser(mmUserID string, user *serializers.TwitterUser) error {
	return s.StoreJSON(util.HashKey(twitterUserPrefix, mmUserID), user)
}

func (s Store) GetTwitterUser(mmUserID string) (*serializers.TwitterUser, error) {
	var user serializers.TwitterUser
	err := s.LoadJSON(util.HashKey(twitterUserPrefix, mmUserID), &user)
	if err != nil {
		if err != ErrNotFound {
			s.api.LogError("Failed to get connected twitter user.", "userID", mmUserID, "error", err.Error())
		}
		return nil, err
	}
	return &user, nil
}

func (s Store) DeleteTwitterUser(mmUserID string) error {
	return s.Delete(util.HashKey(twitterUserPrefix, mmUserID))
}
