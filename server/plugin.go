package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gorilla/mux"
	cmd2 "github.com/mattermost/mattermost-plugin-api/experimental/command"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-twitter/server/api"
	"github.com/mattermost/mattermost-plugin-twitter/server/command"
	"github.com/mattermost/mattermost-plugin-twitter/server/config"
	"github.com/mattermost/mattermost-plugin-twitter/server/store"
	"github.com/mattermost/mattermost-plugin-twitter/server/util"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
// See https://developers.mattermost.com/extend/plugins/server/reference/
type Plugin struct {
	plugin.MattermostPlugin

	router *mux.Router
	store  *store.Store
}

func (p *Plugin) OnActivate() error {
	if err := p.registerCommand(); err != nil {
		p.API.LogError(err.Error())
		return err
	}

	p.store = store.NewStore(p.API, p.Helpers)
	p.router = api.NewController(p.API, p.Helpers, manifest, p.store).InitAPI()
	return nil
}

// OnConfigurationChange is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	var configuration config.Configuration

	// Load the public configuration fields from the Mattermost server configuration.
	if err := p.API.LoadPluginConfiguration(&configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}

	if err := configuration.IsValid(); err != nil {
		return errors.Wrap(err, "failed to validate plugin configuration")
	}

	config.SetConfig(&configuration)
	return nil
}

func (p *Plugin) registerCommand() error {
	iconData, err := cmd2.GetIconData(p.API, "assets/logo.svg")
	if err != nil {
		return errors.Wrap(err, "failed to get icon data")
	}

	cmd := command.GetCommand(iconData)
	if err := p.API.RegisterCommand(cmd); err != nil {
		return errors.Wrap(err, "failed to register slash command: "+cmd.Trigger)
	}

	return nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	defer func() {
		if x := recover(); x != nil {
			p.API.LogError("Recovered from a panic while executing slash command.",
				"commandArgs", fmt.Sprintf("%v", args),
				"error", x,
				"stack", string(debug.Stack()))
		}
	}()

	split, argErr := util.SplitArgs(args.Command)
	if argErr != nil {
		return util.SendEphemeralCommandResponse(argErr.Error())
	}

	cmdName := split[0][1:]
	var params []string

	if len(split) > 1 {
		params = split[1:]
	}

	cmd := command.GetCommand("")
	if cmd.Trigger != cmdName {
		return util.SendEphemeralCommandResponse("Unknown command: [" + cmdName + "] encountered")
	}

	p.API.LogDebug("Executing command: " + cmdName + " with params: [" + strings.Join(params, ", ") + "]")
	cmdContext := command.NewContext(args, c, p.API, p.Helpers, manifest, p.store)
	return command.TwitterCommandHandler.Handle(cmdContext, params...)
}

// ServeHTTP handles HTTP requests for the plugin.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.API.LogDebug("New request:", "Host", r.Host, "RequestURI", r.RequestURI, "Method", r.Method)

	if err := config.GetConfig().IsValid(); err != nil {
		p.API.LogError("This plugin is not configured.", "Error", err.Error())
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	p.router.ServeHTTP(w, r)
}
