# Go Go BOSH - BOSH client API for golang applications

This project is a golang library for applications wanting to talk to a BOSH/MicroBOSH or bosh-lite.

* [![GoDoc](https://godoc.org/github.com/cloudfoundry-community/gogobosh?status.png)](https://godoc.org/github.com/cloudfoundry-community/gogobosh)
* Test status [![Build Status](https://travis-ci.org/cloudfoundry-community/gogobosh.svg)](https://travis-ci.org/cloudfoundry-community/gogobosh)


## API

The following client functions are available, as a subset of the full BOSH Director API.

* client.GetInfo()
* client.GetStemcells()
* client.GetReleases()
* client.GetDeployments()
* client.GetDeployment("cf-warden")
* client.GetDeploymentVMs("cf-warden")
* client.GetTasks()
* client.GetTask(123)
* client.GetTaskResult(123)

## Install

```
go get github.com/cloudfoundry-community/gogobosh
````

## Documentation

The documentation is published to [https://godoc.org/github.com/cloudfoundry-community/gogobosh](https://godoc.org/github.com/cloudfoundry-community/gogobosh).

Also, view the documentation locally with:

```
godoc -goroot=$GOPATH github.com/cloudfoundry-community/gogobosh
```

### Use

As a short getting started guide:

``` golang
package main

import (
  "github.com/cloudfoundry-community/gogobosh"
  "fmt"
)

func main() {
  c, _ := gogobosh.NewClient(gogobosh.DefaultConfig())
  info, _ := c.GetInfo()

  fmt.Println("Director")
  fmt.Printf("  Name       %s\n", info.Name)
  fmt.Printf("  Version    %s\n", info.Version)
  fmt.Printf("  User       %s\n", info.User)
  fmt.Printf("  UUID       %s\n", info.UUID)
  fmt.Printf("  CPI        %s\n", info.CPI)
}
```

## Tests

Tests are all local currently; and do not test against a running bosh or bosh-lite. I'd like to at least do integration tests against a bosh-lite in future.
