package main

import chat "google.golang.org/api/chat/v1"

// HelloWorldPlugin replies to every message with a 'Hello world!' message.
// Pretty dumby... but working :D
type HelloWorldPlugin struct{}

// Handle method implements the Plugin interface.
func (p *HelloWorldPlugin) Handle(channel *chat.SpacesMessagesService, message *chat.Message) error {
	responseMessage := &chat.Message{}
	if message.Thread != nil {
		responseMessage.Thread = message.Thread
	}
	responseMessage.Text = "Hello world!"

	channel.Create(message.Space.Name, responseMessage).Do()

	return nil
}
