package jsonrpc

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Header struct {
	Name  string
	Value string
}

// Get the Key and Value from a header string
// Tested against a strings.split alternative, performs much better
func GetKV(h string) (*Header, error) {
	m := regexp.MustCompile(`(.+?): (.+)\r\n`)
	f := m.FindStringSubmatch(h)
	if len(f) != 3 {
		return nil, fmt.Errorf("Could not parse header %s", h)
	}
	return &Header{Name: f[1], Value: f[2]}, nil
}

// ReadHeader consumes bytes until it matches a Header line
func ReadHeader(r io.Reader) (*Header, error) {
	// One byte at a time
	nextChar := make([]byte, 1)
	buffer := ""
	for {
		_, err := r.Read(nextChar)
		if err != nil {
			return nil, err
		}
		// Build up buffer
		buffer += string(nextChar)
		// Found the character indicating end of header
		if strings.HasSuffix(buffer, "\n") {
			if len(buffer) == 2 {
				// Buffer is exactly \r\n, this is the end of the header block
				return nil, nil
			}
			return GetKV(buffer)
		}
	}
}

// ReadHeaders Reads all the headers from a io.Reader
func ReadHeaders(r io.Reader) (Headers, error) {
	headers := make(Headers)
	for {
		// Get the next header
		h, err := ReadHeader(r)
		if err != nil {
			return nil, err
		}
		// No header and no error means the end of the header section, exit
		if h == nil {
			return headers, nil
		}
		// Normalise the keys in the map
		headers[strings.ToLower(h.Name)] = h.Value
	}
}

// GetContentLength returns the value of Content-Length in the provided headers as int
func GetContentLength(headers Headers) (int, error) {
	contentLengthString, ok := headers["content-length"]
	if !ok {
		return -1, fmt.Errorf("Missing Content-Length in request")
	}

	size, err := strconv.Atoi(contentLengthString)
	if err != nil {
		return -1, err
	}
	return size, nil
}

// ReadRequest consumes the input until it parsed all parts of a request
func ReadRequest(r io.Reader) (*Request, *ResponseError) {
	headers, err := ReadHeaders(r)
	if err != nil {
		return nil, NewError(ParseError, err.Error(), nil)
	}

	size, err := GetContentLength(headers)
	if err != nil {
		return nil, NewError(ParseError, err.Error(), nil)
	}
	request, err := ReadBody(r, size)
	if err != nil {
		return nil, NewError(ParseError, err.Error(), nil)
	}
	request.Headers = headers
	return request, nil
}

// ReadBody consumes the input for a given size and parses the extracted JSON
func ReadBody(r io.Reader, size int) (*Request, error) {
	packet := make([]byte, size)
	_, err := r.Read(packet)
	if err != nil {
		return nil, err
	}
	var request *Request
	err = json.Unmarshal(packet, &request)
	if err != nil {
		return nil, err
	}
	return request, nil
}
