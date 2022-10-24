# Plugin Starter Template [![CircleCI branch](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-twitter/master.svg)](https://circleci.com/gh/mattermost/mattermost-plugin-starter-template)

A Mattermost plugin to connect to twitter.

## Getting Started

To learn more about plugins, see [our plugin documentation](https://developers.mattermost.com/extend/plugins/).

Build your plugin:
```
make
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/com.mattermost.twitter.tar.gz
```

## Configuration

Getting the Twitter Consumer Key (API Key) and Consumer Secret key is very simple, just follow the below 4 steps and you are ready to go.

- Go to https://dev.twitter.com/apps/new and log in, if necessary
- Supply the necessary required fields, accept the Terms Of Service, and solve the CAPTCHA.
- Submit the form
- Go to the API Keys tab, there you will find your Consumer key and Consumer secret keys.
- Copy the consumer key (API key) and consumer secret from the screen into our application.


Enable the 3-legged OAuth.
- In your app settings page of the app you just created, select `Enable 3-legged OAuth`.
https://developer.twitter.com/en/portal/projects/<project-id>/apps/<app-id>/auth-settings

- Set the callbackURL to `<your-mattermost-url>/plugins/com.mattermost.twitter/twitter/callback`.
- Set the Website URL to `your-mattermost-url`.

## Development

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options.

### Deploying with Local Mode

If your Mattermost server is running locally, you can enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode) to streamline deploying your plugin. Edit your server configuration as follows:

```json
{
    "ServiceSettings": {
        ...
        "EnableLocalMode": true,
        "LocalModeSocketLocation": "/var/tmp/mattermost_local.socket"
    }
}
```

and then deploy your plugin:
```
make deploy
```

You may also customize the Unix socket path:
```
export MM_LOCALSOCKETPATH=/var/tmp/alternate_local.socket
make deploy
```

If developing a plugin with a webapp, watch for changes and deploy those automatically:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make watch
```

### Deploying with credentials

Alternatively, you can authenticate with the server's API with credentials:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```

or with a [personal access token](https://docs.mattermost.com/developer/personal-access-tokens.html):
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
```
