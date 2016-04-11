// Copyright (c) 2015 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer
//    in this position and unchanged.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package recorder

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	"github.com/dnaeon/go-vcr/cassette"
)

// Recorder states
const (
	ModeRecording = iota
	ModeReplaying
)

// Recorder is capable of recording an interaction and replaying a previously
// recorded one.
type Recorder struct {
	// Operating mode of the recorder
	mode int

	// HTTP server used to mock requests
	server *httptest.Server

	// Cassette used by the recorder
	cassette *cassette.Cassette

	// Transport that can be used by clients to inject
	Transport *http.Transport
}

// Proxies client requests to their original destination
func requestHandler(r *http.Request, c *cassette.Cassette, mode int) (*cassette.Interaction, error) {
	if r.Method == "CONNECT" {
		r.URL.Scheme = "https"
	}
	// Return interaction from cassette if in replay mode
	if mode == ModeReplaying {
		return c.GetInteraction(r)
	}

	// Else, perform client request to their original
	// destination and record interactions
	body := ioutil.NopCloser(r.Body)
	req, err := http.NewRequest(r.Method, r.URL.String(), body)
	fmt.Printf("request made using the proxy %v", req)
	if err != nil {
		return nil, err
	}

	req.Header = r.Header
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Record the interaction and add it to the cassette
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Add interaction to cassette
	interaction := &cassette.Interaction{
		Request: cassette.Request{
			Body:    string(reqBody),
			Headers: req.Header,
			URL:     req.URL.String(),
			Method:  req.Method,
		},
		Response: cassette.Response{
			Body:    string(respBody),
			Headers: resp.Header,
			Status:  resp.Status,
			Code:    resp.StatusCode,
		},
	}
	c.AddInteraction(interaction)

	return interaction, nil
}

// New Creates a new recorder
func New(cassetteName string) (*Recorder, error) {
	var mode int
	var c *cassette.Cassette
	cassetteFile := fmt.Sprintf("%s.yaml", cassetteName)

	// Depending on whether the cassette file exists or not we
	// either create a new empty cassette or load from file
	if _, err := os.Stat(cassetteFile); os.IsNotExist(err) {
		// Create new cassette and enter in recording mode
		c = cassette.New(cassetteName)
		mode = ModeRecording
	} else {
		// Load cassette from file and enter replay mode
		c, err = cassette.Load(cassetteName)
		if err != nil {
			return nil, err
		}
		mode = ModeReplaying
	}

	// Handler for client requests
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pass cassette and mode to handler, so that interactions can be
		// retrieved or recorded depending on the current recorder mode
		interaction, err := requestHandler(r, c, mode)

		if err != nil {
			log.Fatalf("Failed to process request for URL %s: %s", r.URL, err)
		}

		w.WriteHeader(interaction.Response.Code)
		fmt.Fprintln(w, interaction.Response.Body)
	})

	// HTTP server used to mock requests
	server := httptest.NewServer(handler)

	// A proxy function which routes all requests through our HTTP server
	// Can be used by clients to inject into their own transports
	proxyURL, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}

	// A transport which can be used by clients to inject
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	r := &Recorder{
		mode:      mode,
		server:    server,
		cassette:  c,
		Transport: transport,
	}

	return r, nil
}

// Stop Stops the recorder gracefully, closing its dependencies as well.
func (r *Recorder) Stop() error {
	r.server.Close()

	if r.mode == ModeRecording {
		if err := r.cassette.Save(); err != nil {
			return err
		}
	}

	return nil
}
