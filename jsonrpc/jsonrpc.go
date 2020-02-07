package jsonrpc

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type NotificationsProvider interface {
	Notifications() chan *Notification
}

type Headers map[string]string

type Request struct {
	Headers Headers         `json:"-"`
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type Response struct {
	Result interface{}    `json:"result"`
	Error  *ResponseError `json:"error"`
	ID     int            `json:"id"`
}

type Notification struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type Resolution struct {
	ID    int
	Reply bool
	Res   interface{}
	Err   *ResponseError
}

type Handler func(json.RawMessage) (interface{}, *ResponseError)
type NotificationHandler func(json.RawMessage)

type Server struct {
	handlers              map[string]Handler
	notificationHandlers  map[string]NotificationHandler
	requests              chan *Request
	out                   chan *Resolution
	Logger                *log.Logger
	Reader                io.Reader
	Writer                io.Writer
	notificationsProvider NotificationsProvider
}

func NewServer(in io.Reader, out io.Writer) *Server {
	return &Server{
		handlers:              make(map[string]Handler),
		notificationHandlers:  make(map[string]NotificationHandler),
		requests:              make(chan *Request, 8),
		out:                   make(chan *Resolution, 8),
		Logger:                log.New(ioutil.Discard, "[jsonrpc] ", log.Ldate|log.Ltime),
		Reader:                in,
		Writer:                out,
		notificationsProvider: nil,
	}
}

func (s *Server) SetNotificationsProvider(provider NotificationsProvider) {
	s.notificationsProvider = provider
}

// AddHandler registers a handler for a given method
func (s *Server) AddHandler(method string, handler Handler) {
	_, exists := s.handlers[method]
	if exists {
		s.HandleError(fmt.Errorf("Handler redefinition: %s", method))
		return
	}
	s.handlers[method] = handler
}

// AddHandler registers a notification handler for a given method
func (s *Server) AddNotificationHandler(method string, handler NotificationHandler) {
	_, exists := s.notificationHandlers[method]
	if exists {
		s.HandleError(fmt.Errorf("Notification handler redefinition: %s", method))
		return
	}
	s.notificationHandlers[method] = handler
}

// GetResponse will match a handler to the request/notification,
// resolve the response, then send a response if there are any
func (s *Server) GetResponse(r *Request) (bool, interface{}, *ResponseError) {
	// Try to match a request/response handler
	handler := s.handlers[r.Method]
	s.Logger.Printf("Found Handler for method %s\n", r.Method)
	if handler == nil {
		// None found, try to match a notification handler
		handler := s.notificationHandlers[r.Method]
		if handler == nil {
			// Handler not found at all
			s.Logger.Printf("Method not found %s\n", r.Method)
			return true, nil, NewError(MethodNotFound, "", nil)
		}
		// Call the notification handler
		// Use a goroutine to resolve this request as fast as possible
		go handler(r.Params)
		return false, nil, nil
	}
	// Call the request handler
	res, err := handler(r.Params)
	if err != nil {
		return true, nil, err
	}
	return true, res, nil
}

// HandleRequest resolves the response and send it down the output
func (s *Server) HandleRequest(req *Request) {
	reply, res, err := s.GetResponse(req)
	resolution := &Resolution{Reply: reply, Res: res, Err: err, ID: req.ID}
	s.HandleResponse(resolution)
}

// HandleResponse will send a response or error down the output
// based on the properties of a given resolution
func (s *Server) HandleResponse(resolution *Resolution) {
	if resolution.Err != nil {
		err := s.PrintError(resolution.ID, resolution.Err)
		if err != nil {
			s.HandleError(err)
		}
	} else if resolution.Reply {
		// Notifications, will return false for reply
		err := s.PrintResponse(resolution.ID, resolution.Res, nil)
		if err != nil {
			s.HandleError(err)
		}
	}
}

// HandleError simply prints the error and exits the process
func (s *Server) HandleError(err error) {
	s.Logger.Panic(err)
}

// Listen continuously reads the input for requests
// It uses goroutines to handle requests as they come
func (s *Server) Listen() {
	go s.SendNotifications()
	for {
		req, readErr := ReadRequest(s.Reader)
		if readErr != nil {
			go s.HandleResponse(&Resolution{Err: readErr, Reply: true, Res: nil})
			return
		}
		go s.HandleRequest(req)
	}
}

func (s *Server) SendNotifications() {
	if s.notificationsProvider == nil {
		return
	}
	for {
		notifsChan := s.notificationsProvider.Notifications()
		n := <-notifsChan
		err := s.PrintNotification(n)
		if err != nil {
			s.HandleError(err)
		}
	}
}

// PrintError send an error back to the client
func (s *Server) PrintError(id int, err *ResponseError) error {
	return s.PrintResponse(id, nil, err)
}

// PrintResponse sends a response back to the client
func (s *Server) PrintResponse(id int, contents interface{}, resErr *ResponseError) error {
	// Build the response object
	res := &Response{ID: id, Error: resErr, Result: contents}
	jsonString, err := json.Marshal(res)
	if err != nil {
		return err
	}
	s.Logger.Printf("Sending response: %s\n", jsonString)
	_, err = fmt.Fprintf(s.Writer, "Content-Length: %d\r\n\r\n%s", len(jsonString), jsonString)
	if err != nil {
		return err
	}
	return nil
}

// PrintNotification sends a notification back to the client
func (s *Server) PrintNotification(notification *Notification) error {
	jsonString, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	s.Logger.Printf("Sending notification: %s\n", jsonString)
	_, err = fmt.Fprintf(s.Writer, "Content-Length: %d\r\n\r\n%s", len(jsonString), jsonString)
	if err != nil {
		return err
	}
	return nil
}
