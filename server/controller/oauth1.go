package controller

import (
	"fmt"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-twitter/server/serializers"
	"github.com/mattermost/mattermost-plugin-twitter/server/util"
)

func (c *Controller) twitterLoginCallback(w http.ResponseWriter, r *http.Request) {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(r)
	if err != nil {
		c.api.LogError("twitterLoginCallback: Failed to parse authorisation callback.", "Error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mmUserID := r.Header.Get("Mattermost-User-ID")
	if mmUserID == "" {
		c.api.LogError("twitterLoginCallback: Failed to get mattermost userID.")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	mmUser, appErr := c.api.GetUser(mmUserID)
	if appErr != nil {
		c.api.LogError("twitterLoginCallback: Failed to get mattermost user.", "Error", appErr.Error())
		http.Error(w, appErr.Error(), http.StatusInternalServerError)
		return
	}

	var oauthTmpCredentials serializers.OAuth1aTemporaryCredentials
	if storeErr := c.store.LoadOneTimeSecretJSON(mmUserID, &oauthTmpCredentials); storeErr != nil || len(oauthTmpCredentials.Token) == 0 {
		c.api.LogError(fmt.Sprintf("twitterLoginCallback: Failed to load oauth one-time secret. Error: %v", storeErr))
		http.Error(w, fmt.Sprintf("temporary credentials for %s not found or expired, try to connect again", mmUserID), http.StatusInternalServerError)
		return
	}

	if oauthTmpCredentials.Token != requestToken {
		c.api.LogError("twitterLoginCallback: saved OAuth credentials and request token do not match.")
		http.Error(w, "request token mismatch", http.StatusBadRequest)
		return
	}

	twitterOAuth1Config := util.GetTwitterOAuth1Config(c.api, c.manifest)

	// Twitter ignores the oauth_signature on the access token request. The user
	// to which the request (temporary) token corresponds is already known on the
	// server. The request for a request token earlier was validated signed by
	// the consumer. Consumer applications can avoid keeping request token state
	// between authorization granting and callback handling.
	accessToken, accessSecret, err := twitterOAuth1Config.AccessToken(requestToken, "", verifier)
	if err != nil {
		c.api.LogError("twitterLoginCallback: Failed to get AccessToken from request.", "Error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Twitter client
	client := util.GetTwitterClient(accessToken, accessSecret)
	twUser, resp, err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	if err != nil {
		c.api.LogError("twitterLoginCallback: Failed to verify twitter credentials for connected user.", "Error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	if err := c.store.SaveTwitterUser(mmUserID, &serializers.TwitterUser{
		Name:         twUser.Name,
		Username:     twUser.ScreenName,
		AccessToken:  accessToken,
		AccessSecret: accessSecret,
	}); err != nil {
		c.api.LogError("twitterLoginCallback: Failed to save twitter client to KVStore.", "Error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.renderTemplate(w, "oauth1-complete.html", "text/html", map[string]string{
		"TwitterDisplayName":    twUser.Name + " (@" + twUser.ScreenName + ")",
		"MattermostDisplayName": mmUser.GetDisplayName(model.SHOW_NICKNAME_FULLNAME),
	})
}
