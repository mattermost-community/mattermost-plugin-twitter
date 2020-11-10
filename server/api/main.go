package api

import (
	"net/http"
	"path/filepath"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-twitter/server/constant"
	"github.com/mattermost/mattermost-plugin-twitter/server/store"
)

const (
	HeaderMattermostUserID = "Mattermost-User-Id"
)

type Controller struct {
	api      plugin.API
	helpers  plugin.Helpers
	manifest *model.Manifest
	store    *store.Store
}

func NewController(api plugin.API, helpers plugin.Helpers, manifest *model.Manifest, store *store.Store) *Controller {
	return &Controller{
		api,
		helpers,
		manifest,
		store,
	}
}

// InitAPI initializes the REST API
func (c *Controller) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.Use(c.withRecovery)

	c.handleStaticFiles(r)
	s := r.PathPrefix("/api/v1").Subrouter()

	// Add the custom plugin routes here
	s.HandleFunc(constant.PathTwitterOAuth1Callback, handleAuthRequired(c.twitterLoginCallback)).Methods(http.MethodGet)

	// 404 handler
	r.Handle("{anything:.*}", http.NotFoundHandler())
	return r
}

// From: https://github.com/mattermost/mattermost-plugin-github/blob/42185ff874963bed1efd8bc84c81462184d7cca8/server/plugin/api.go#L135
func (c *Controller) withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				c.api.LogError("Recovered from a panic",
					"url", r.URL.String(),
					"error", x,
					"stack", string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// handleStaticFiles handles the static files under the assets directory.
func (c *Controller) handleStaticFiles(r *mux.Router) {
	bundlePath, err := c.api.GetBundlePath()
	if err != nil {
		c.api.LogWarn("Failed to get bundle path.", "Error", err.Error())
		return
	}

	// This will serve static files from the 'assets' directory under '/static/<filename>'
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(bundlePath, "assets")))))
}

// handleAuthRequired verifies if provided request is performed by a logged-in Mattermost user.
func handleAuthRequired(handleFunc func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(HeaderMattermostUserID)
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		handleFunc(w, r)
	}
}
