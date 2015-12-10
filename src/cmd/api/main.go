package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// StringService provides operations on strings.
type StringService interface {
	Uppercase(string) (string, error)
}

type stringService struct{}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func main() {
	ctx := context.Background()
	svc := stringService{}

	uppercaseHandler := httptransport.NewServer(
		ctx,
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	allServicesHandler := httptransport.NewServer(
		ctx,
		allServicesHandler(microServicesService{}),
		passThrough,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/services", allServicesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type microServicesService struct{}

// MicroServicesService dasdas
type MicroServicesService interface {
	AllServices() string
}

func (microServicesService) AllServices() string {
	return "allservices"
}

func allServicesHandler(svc MicroServicesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.AllServices(), nil
	}
}

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func passThrough(r *http.Request) (interface{}, error) {
	return r, nil
}

func decodeUppercaseRequest(r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type allServicesRequest struct{}

type allServicesResponse struct {
	V string `json:"v"`
}

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")
