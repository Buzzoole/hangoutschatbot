# Hangouts Chat Bot

An Hangouts Chat bot, written in Go 🤖

## What it is

_hangoutschatbot_ is a framework useful to provide developers with a boilerplate application they can rely on and create their own chat bots based on [Google's Hangouts Chat](https://chat.google.com/), a service available in the G Suite.

## Usage

### Requirements

Environment:

- [go environment 1.10](https://golang.org/doc/install)
- [dep](https://github.com/golang/dep)

In order to make the chat bot work on your G Suite, you'll need:

- An active [Google project](https://console.cloud.google.com/projectcreate)
- [Hangouts Chat API](https://console.cloud.google.com/apis/api/chat.googleapis.com) enabled on that project
- A [service account](https://console.cloud.google.com/iam-admin/serviceaccounts), enabled on those APIs (you'll need the JSON credentials file to run the chat bot)
- To configure the [Hangouts Chat API](https://console.cloud.google.com/apis/api/chat.googleapis.com/hangouts-chat) to use the HTTP address where the HTTP bot server will be accessible, listening for new messages (an example is provided [here](./docs/bot_configuration.png))

### Building & Execution

To simplify the process, a convenient Makefile is provided.

#### Getting dependencies

Run:

```bash
make get-deps
```

#### Building the bot

Run:

```bash
make build
```

## Concepts

### Application Structure

At its very core, _hangoutschatbot_ consists of:

- An HTTP server, used as webhook from Hangouts Chat to forward all the messages coming from chat rooms
- A SpacesMessagesService, used to publish new messages or replies to existing messages to Hangouts Chat
- A list of registered plugins, which are in charge of reacting to chat messages

### Plugins

Plugins do actually implement the behaviour of every chat bot.

Each plugin is responsible of deciding whether to reply (or not) to a message, according to its own business logic.

### Creating a Plugin

New Plugins can be added by implementing the following interface:

```go
func (p *SomePlugin) Handle(channel *chat.SpacesMessagesService, message *chat.Message) error {}
```

Where `chat` is the `google.golang.org/api/chat/v1` package.

Do not forget to _register_ the new plugin in the `main.go` file, e.g.

```go
cs.RegisterPlugin(new(HelloWorldPlugin))
```

## Reference

- [Hangouts Chat API main documentation page](https://developers.google.com/hangouts/chat/)
- [Creating new bots](https://developers.google.com/hangouts/chat/how-tos/bots-develop)