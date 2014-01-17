# Go Go BOSH - BOSH client API for golang applications

This project is a golang library for applications wanting to talk to a BOSH/MicroBOSH or bosh-lite.

## API compatibility

The BOSH Core team nor its Product Managers have never claimed that a BOSH director has a public API; and they want to make changes to the API in the future. I have no idea how to support a future, un-documented, un-versioned API. It'll be tricky.

The best way to describe the API support in this library is to document what version of bosh-lite is being tested against, the date that it was published. Hopefully bosh-lite is always approximately parallel (via rebasing) in its API with the main BOSH project; and the same timestamps can map to the continuously delivered releases of BOSH & its RubyGems.

Trying to write a client library for an API without any versioning strategy could get messy for client applications. Please write your own integration tests that work against running BOSHes that you'll use in production.

If you are using this library, or the Ruby library within the `bosh_cli` rubygem, or talking directly with the BOSH director API - please announce yourself on the bosh-users google group and/or to the PM of BOSH. This way they can be aware of who many be affected by API changes.

## Tests

The integration tests assume that bosh-lite is running locally.