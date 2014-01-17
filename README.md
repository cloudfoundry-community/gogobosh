# Go Go BOSH - BOSH client API for golang applications

This project is a golang library for applications wanting to talk to a BOSH/MicroBOSH or bosh-lite.

## API compatibility

The BOSH Core team nor its Product Managers have never claimed that a BOSH director has a public API; and they want to make changes to the API in the future. I have no idea how to support a future, un-documented, un-versioned API. It'll be tricky.

The best way to describe the API support in this library is to document what version of bosh-lite is being tested against, the date that it was published. Hopefully bosh-lite is always approximately parallel (via rebasing) in its API with the main BOSH project; and the same timestamps can map to the continuously delivered releases of BOSH & its RubyGems.

Trying to write a client library for an API without any versioning strategy could get messy for client applications. Please write your own integration tests that work against running BOSHes that you'll use in production.

If you are using this library, or the Ruby library within the `bosh_cli` rubygem, or talking directly with the BOSH director API - please announce yourself on the bosh-users google group and/or to the PM of BOSH. This way they can be aware of who many be affected by API changes.

## Install

```
go get github.com/cloudfoundry-community/gogobosh
````

### Use

``` golang
package main

import (
    bosh "github.com/cloudfoundry-community/gogobosh"
)

func main() {
    director := bosh.New("https://192.168.50.4:25555", "admin", "admin")
    fmt.Println("Director")
    fmt.Printf("  Name       %s", director.Name)
    fmt.Printf("  URL        %s", director.URL)
    fmt.Printf("  Version    %s", director.Version)
    fmt.Printf("  User       %s", director.User)
    fmt.Printf("  UUID       %s", director.UUID)
    fmt.Printf("  CPI        %s", director.CPI)
    fmt.Printf("  dns        %s", director.DNSEnabled)
    fmt.Printf("  compiled_package_cache %s (provider: %s)", director.CompiledPackageCacheEnabled, director.CompiledPackageCacheProvider)
    fmt.Printf("  snapshots  %s", director.SnapshotsEnabled)
}
```

## Tests

The integration tests assume that bosh-lite is running locally.