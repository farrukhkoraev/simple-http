package request

import (
	"fmt"
	"strings"
)

var StatusMap = map[int]string{
	200: "Success",
	201: "Created",
	404: "Not Found",
}

type Request struct {
	Method   string
	Path     string
	Protocol string
	Headers  map[string]string
	Body     string
}

type Response struct {
	Protocol   string
	StatusCode int
	Headers    map[string]string
	Body       string
}
type ParseError struct {
	Msg string
}

func (err *ParseError) Error() string {
	return err.Msg
}

func ParseRequest(message string) (*Request, error) {
	req := &Request{Headers: make(map[string]string)}

	// parse start line (method, path, version)
	requestLine, rest, found := strings.Cut(message, "\n")
	if !found {
		return nil, &ParseError{"Error parsing start-line"}
	}

	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, &ParseError{"Error parsing start-line"}
	}
	req.Method = parts[0]
	req.Path = parts[1]
	req.Protocol = parts[2]

	// parse headers
	headers, body, found := strings.Cut(rest, "\n\n")
	if !found {
		return nil, &ParseError{"Error parsing headers"}
	}
	for _, header := range strings.Split(headers, "\n") {
		key, value, _ := strings.Cut(header, ":")
		req.Headers[key] = strings.TrimPrefix(value, " ")
	}

	req.Body = body

	return req, nil
}

func SerializeResponse(response *Response) (string, error) {
	builder := strings.Builder{}
	_, err := builder.WriteString(
		fmt.Sprintf("%s %d %s\n",
			response.Protocol,
			response.StatusCode,
			StatusMap[response.StatusCode]),
	)
	if err != nil {
		return "", err
	}

	for key, value := range response.Headers {
		_, err = builder.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		if err != nil {
			return "", err
		}
	}

	builder.WriteString("\n")
	_, err = builder.WriteString(response.Body)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}
