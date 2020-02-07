package extension

import (
	"io/ioutil"
	"log"
	"net/url"
	"runtime"
	"taskfile-language-server/jsonrpc"
	"taskfile-language-server/taskfile"

	"github.com/sourcegraph/go-lsp"
)

func GetPath(uri lsp.DocumentURI) (string, error) {
	path := ""
	url, err := url.Parse(string(uri))
	if err != nil {
		return path, err
	}
	path = url.Path
	// Remove the laeding slash if on windows
	if runtime.GOOS == "windows" {
		path = url.Path[1:]
	}
	return path, nil
}

func reloadTaskfile(docUri lsp.DocumentURI, text string) error {
	path, err := GetPath(docUri)
	if err != nil {
		return err
	}
	taskfile.Invalidate(path, text)
	return nil
}

type TaskfileExtension struct {
	Logger        *log.Logger
	notifications chan *jsonrpc.Notification
}

func New() *TaskfileExtension {
	return &TaskfileExtension{
		Logger:        log.New(ioutil.Discard, "[taskfile]", log.Ldate|log.Ltime),
		notifications: make(chan *jsonrpc.Notification),
	}
}

func (t *TaskfileExtension) RegisterHandlers(s *jsonrpc.Server) {
	s.AddHandler("extension/getTasks", t.GetTasks)
}

func (t *TaskfileExtension) SendNotification(method string, contents interface{}) {
	t.notifications <- &jsonrpc.Notification{
		Method: method,
		Params: contents,
	}
}

func (t *TaskfileExtension) Notifications() chan *jsonrpc.Notification {
	return t.notifications
}
