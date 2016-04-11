## go-vcr

[![Build Status](https://travis-ci.org/dnaeon/go-vcr.svg)](https://travis-ci.org/dnaeon/go-vcr)

`go-vcr` simplifies testing by recording your HTTP interactions and
replaying them in future runs in order to provide fast, deterministic
and accurate testing of your code.

`go-vcr` was inspired by the [VCR library for Ruby](https://github.com/vcr/vcr).

## Installation

Install `go-vcr` by executing the command below:

```bash
$ go get github.com/dnaeon/go-vcr/recorder
```

## Usage

Here is a simple example of recording and replaying
[etcd](https://github.com/coreos/etcd) HTTP interactions.

You can find other examples in the `example` directory of this
repository as well.

```go
package main

import (
	"log"
	"time"

	"github.com/dnaeon/go-vcr/recorder"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func main() {
	// Start our recorder
	r, err := recorder.New("fixtures/etcd")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	// Create an etcd configuration using our transport
	cfg := client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		HeaderTimeoutPerRequest: time.Second,
		Transport:               r.Transport, // Inject our transport!
	}

	// Create an etcd client using the above configuration
	c, err := client.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create etcd client: %s", err)
	}

	// Get an example key from etcd
	etcdKey := "/foo"
	kapi := client.NewKeysAPI(c)
	resp, err := kapi.Get(context.Background(), etcdKey, nil)

	if err != nil {
		log.Fatalf("Failed to get etcd key %s: %s", etcdKey, err)
	}

	log.Printf("Successfully retrieved etcd key %s: %s", etcdKey, resp.Node.Value)
}
```

## License

`go-vcr` is Open Source and licensed under the
[BSD License](http://opensource.org/licenses/BSD-2-Clause)
