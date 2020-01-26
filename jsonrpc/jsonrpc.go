package jsonrpc

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

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

type Handler func(json.RawMessage) (interface{}, *ResponseError)
type NotificationHandler func(json.RawMessage)

type Server struct {
	handlers             map[string]Handler
	notificationHandlers map[string]NotificationHandler
	requests             chan *Request
	out                  chan *Resolution
	Logger               *log.Logger
}

func NewServer() *Server {
	return &Server{
		handlers:             make(map[string]Handler),
		notificationHandlers: make(map[string]NotificationHandler),
		requests:             make(chan *Request, 8),
		out:                  make(chan *Resolution, 8),
		Logger:               log.New(ioutil.Discard, "[jsonrpc] ", 0),
	}
}

func (s *Server) AddHandler(method string, handler Handler) {
	s.handlers[method] = handler
}

func (s *Server) AddNotificationHandler(method string, handler NotificationHandler) {
	s.notificationHandlers[method] = handler
}

func (s *Server) GetResponse(r *Request) (bool, interface{}, *ResponseError) {
	handler := s.handlers[r.Method]
	s.Logger.Printf("Found Handler for method %s\n", r.Method)
	if handler == nil {
		handler := s.notificationHandlers[r.Method]
		if handler == nil {
			s.Logger.Printf("Method not found %s\n", r.Method)
			return true, nil, NewError(MethodNotFound, "", nil)
		}
		handler(r.Params)
		return false, nil, nil
	}
	res, err := handler(r.Params)
	if err != nil {
		return true, nil, err
	}
	return true, res, nil
}

type Resolution struct {
	ID    int
	Reply bool
	Res   interface{}
	Err   *ResponseError
}

func (s *Server) HandleRequest(req *Request, out io.Writer) {
	reply, res, err := s.GetResponse(req)
	resolution := &Resolution{Reply: reply, Res: res, Err: err, ID: req.ID}
	s.HandleResponse(resolution, out)
}

func (s *Server) HandleResponse(resolution *Resolution, out io.Writer) {
	if resolution.Err != nil {
		err := s.PrintError(out, resolution.ID, resolution.Err)
		if err != nil {
			s.HandleError(err)
		}
	} else if resolution.Reply {
		// Notifications, will return false for reply
		err := s.PrintResponse(out, resolution.ID, resolution.Res, nil)
		if err != nil {
			s.HandleError(err)
		}
	}
}

func (s *Server) HandleError(err error) {
	s.Logger.Fatal(err)
}

func (s *Server) Listen(in io.Reader, out io.Writer) {
	for {
		req, readErr := ReadRequest(in)
		if readErr != nil {
			go s.HandleResponse(&Resolution{Err: readErr, Reply: true, Res: nil}, out)
			return
		}
		go s.HandleRequest(req, out)
	}
}

func (s *Server) PrintError(w io.Writer, id int, err *ResponseError) error {
	return s.PrintResponse(w, id, nil, err)
}

func (s *Server) PrintResponse(w io.Writer, id int, contents interface{}, resErr *ResponseError) error {
	res := &Response{ID: id, Error: resErr, Result: contents}
	jsonString, err := json.Marshal(res)
	if err != nil {
		return err
	}
	s.Logger.Printf("Sending response: %s\n", jsonString)
	_, err = fmt.Fprintf(w, "Content-Length: %d\r\n\r\n%s", len(jsonString), jsonString)
	if err != nil {
		return err
	}
	return nil
}
