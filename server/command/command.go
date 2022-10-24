package command

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-twitter/server/store"
)

// Context includes the context in which the slash command is executed and allows access to
// plugin API, helpers and services
type Context struct {
	*model.CommandArgs
	context  *plugin.Context
	api      plugin.API
	helpers  plugin.Helpers
	manifest *model.Manifest
	store    *store.Store
}

func NewContext(args *model.CommandArgs, context *plugin.Context, api plugin.API, helpers plugin.Helpers, manifest *model.Manifest, store *store.Store) *Context {
	return &Context{
		args,
		context,
		api,
		helpers,
		manifest,
		store,
	}
}

type HandlerFunc func(context *Context, args ...string) (*model.CommandResponse, *model.AppError)

type Handler struct {
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

func (ch Handler) Handle(context *Context, args ...string) (*model.CommandResponse, *model.AppError) {
	for n := len(args); n > 0; n-- {
		h := ch.handlers[strings.Join(args[:n], "/")]
		if h != nil {
			return h(context, args[n:]...)
		}
	}
	return ch.defaultHandler(context, args...)
}
