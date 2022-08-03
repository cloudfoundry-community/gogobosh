# Go Go BOSH - BOSH client API for golang applications

This project is a golang library for applications wanting to talk to a BOSH/MicroBOSH or bosh-lite.

[![GoDoc](https://godoc.org/github.com/cloudfoundry-community/gogobosh?status.png)](https://godoc.org/github.com/cloudfoundry-community/gogobosh)

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

The documentation is published to [https://godoc.org/github.com/cloudfoundry-community/gogobosh](https://godoc.org/github.com/cloudfoundry-community/gogobosh). Also, [view the documentation locally](http://localhost:6060/pkg/github.com/cloudfoundry-community/gogobosh/) with:

```shell
$ godoc
```

## Usage

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

##Development

Tests are currently all local and do not run against bosh or bosh-lite. All the tests can be run using:
```shell
$ make test
```

Before submitting a PR make sure all the tests pass, the code is properly formatted and linted:
```shell
$ make
```

## Contributing

Contributions from the community are welcomed. This is a rough outline of what a contributor's workflow looks like:

- Create a topic branch from where you want to base your work
- Make commits of logical units
- Make sure your commit messages are in the proper format (see below)
- Push your changes to a topic branch in your fork of the repository
- Submit a pull request