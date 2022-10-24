// Package store ...
// Loosely adapted from: https://github.com/mattermost/mattermost-plugin-mscalendar/blob/7d80765da31ac3483354b99197f099d805c3b4a9/server/utils/kvstore/ots.go
package store

import (
	"encoding/json"

	"github.com/mattermost/mattermost-plugin-twitter/server/util"
)

const (
	prefixOneTimeSecret = "ots_" // + unique key that will be deleted after the first verification

	// Expire in 15 minutes
	otsExpiration = 15 * 60
)

func (s *Store) StoreOneTimeSecret(token, secret string) error {
	return s.StoreTTL(util.HashKey(prefixOneTimeSecret, token), []byte(secret), otsExpiration)
}

func (s *Store) LoadOneTimeSecret(key string) (data []byte, returnErr error) {
	data, err := s.Load(util.HashKey(prefixOneTimeSecret, key))
	if len(data) != 0 {
		_ = s.Delete(util.HashKey(prefixOneTimeSecret, key))
	}
	return data, err
}

func (s *Store) StoreOneTimeSecretJSON(token string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.StoreTTL(util.HashKey(prefixOneTimeSecret, token), data, otsExpiration)
}

func (s *Store) LoadOneTimeSecretJSON(key string, v interface{}) (returnErr error) {
	data, err := s.Load(util.HashKey(prefixOneTimeSecret, key))
	if err != nil {
		return err
	}

	// If the key expired, appErr is nil, but the data is also nil
	if len(data) == 0 {
		return ErrNotFound
	}

	_ = s.Delete(util.HashKey(prefixOneTimeSecret, key))
	return json.Unmarshal(data, v)
}
