package extension

import (
	"io/ioutil"
	"log"
	"net/url"
	"runtime"
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
	// Remove the leding slash if on windows
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
	taskfile.PreloadWithBytes(path, []byte(text))
	return nil
}

type TaskfileExtension struct {
	Logger *log.Logger
}

func New() *TaskfileExtension {
	return &TaskfileExtension{Logger: log.New(ioutil.Discard, "[taskfile]", 0)}
}
