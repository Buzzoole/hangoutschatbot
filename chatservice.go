package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	chat "google.golang.org/api/chat/v1"
)

// ChatService is responsible of receiving and sending Hangouts Chat messages,
// handling them via one or more registered plugins.
type ChatService struct {
	broker            chan *chat.Message
	registeredPlugins []Plugin
	quitError         chan error
	sms               *chat.SpacesMessagesService
	httpPort          int
}

// Plugin is the interface which has to be respected to create a new plugin.
type Plugin interface {
	Handle(channel *chat.SpacesMessagesService, message *chat.Message) error
}

// NewChatService returns a new ChatService object relying on Hangouts Chat.
func NewChatService(httpPort int) (*ChatService, error) {
	httpClient, err := google.DefaultClient(oauth2.NoContext, "https://www.googleapis.com/auth/chat.bot")
	if err != nil {
		return nil, fmt.Errorf("error creating httpClient: %v", err)
	}

	chatService, err := chat.New(httpClient)
	if err != nil {
		return nil, fmt.Errorf("error creating chatService: %v", err)
	}

	sms := chat.NewSpacesMessagesService(chatService)

	return &ChatService{
		broker:   make(chan *chat.Message, 0),
		sms:      sms,
		httpPort: httpPort,
	}, nil
}

// RegisterPlugin should be used to add a new Plugin to the ChatService object.
func (cs *ChatService) RegisterPlugin(p Plugin) {
	cs.registeredPlugins = append(cs.registeredPlugins, p)
}

// Serve does actually start the ChatService.
// Each retrieved message is passed to all the registered plugins, which
// are responsible of managing it.
func (cs *ChatService) Serve() error {
	go cs.startHTTPServer()

	for {
		message := <-cs.broker

		for _, plugin := range cs.registeredPlugins {
			err := plugin.Handle(cs.sms, message)
			if err != nil {
				return <-cs.quitError
			}
		}
	}
}

// startHTTPServer is a private method used from the ChatService to run the
// internal HTTP server, which is responsible of getting messages from the
// Hangouts Chat engine.
func (cs *ChatService) startHTTPServer() {
	r := mux.NewRouter()
	hangoutsHandler := &hangoutsHTTPHandler{
		broker: cs.broker,
	}

	r.Handle("/", hangoutsHandler)

	http.Handle("/", r)

	cs.quitError <- http.ListenAndServe(fmt.Sprintf(":%d", cs.httpPort), nil)
}

// hangoutsHTTPHandler is responsible of serving HTTP request coming from
// Hangouts Chat, in order to retrieve new Chat events.
type hangoutsHTTPHandler struct {
	broker chan *chat.Message
}

// ServeHTTP implements the HTTP Handler
func (hhh *hangoutsHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	receivedEvent := &chat.DeprecatedEvent{}
	err = json.Unmarshal(body, receivedEvent)
	if err != nil {
		log.Println(err)
	}

	if receivedEvent.Message != nil {
		hhh.broker <- receivedEvent.Message
	}

	emptyResponse, _ := new(chat.Message).MarshalJSON()
	fmt.Fprintf(w, string(emptyResponse))
}
