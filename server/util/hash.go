package util

import (
	"crypto/md5"
	"fmt"
)

// HashKey returns the kvstore kev by appending prefix with the hash of the input key
// From: https://github.com/mattermost/mattermost-plugin-mscalendar/blob/26fe3c5ea965a435e76dfc5b23e7f66fa9e9b592/server/utils/kvstore/hashed_key.go#L47
// TODO: use a more secure hash primitive
func HashKey(prefix, key string) string {
	if key == "" {
		return prefix
	}

	h := md5.New()
	_, _ = h.Write([]byte(key))
	return fmt.Sprintf("%s%x", prefix, h.Sum(nil))
}
