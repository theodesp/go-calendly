# go-calendly - WIP #

<a href="https://godoc.org/github.com/theodesp/go-calendly/calendly">
<img src="https://godoc.org/github.com/theodesp/go-calendly/calendly?status.svg" alt="GoDoc">
</a>

<a href="https://opensource.org/licenses/Apache-2.0">
<img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License"/>
</a>

<a href="https://travis-ci.org/theodesp/go-calendly" rel="nofollow">
<img src="https://travis-ci.org/theodesp/go-calendly.svg?branch=master" />
</a>

<a href="https://ci.appveyor.com/project/theodesp/go-calendly" rel="nofollow">
<img src="https://ci.appveyor.com/api/projects/status/ytwi6bn3ai6tmd7i/branch/master?svg=true" />
</a>

<a href="https://codecov.io/gh/theodesp/go-calendly">
  <img src="https://codecov.io/gh/theodesp/go-calendly/branch/master/graph/badge.svg" />
</a>

<a href="https://goreportcard.com/report/github.com/theodesp/go-calendly">
  <img src="https://goreportcard.com/badge/github.com/theodesp/go-calendly" />
</a>

go-calendly is a Go client library for accessing the [Calendly API v1](https://developer.calendly.com/docs/getting-started).

go-calendly requires Go version 1.8 or greater.


This is ***WIP!***


## Usage ##

```go
import "github.com/theodesp/go-calendly/calendly"
```

Construct a new Calendly client, then use the various services on the client to
access different parts of the Calendly API. For example:

```go

```

Some API methods have optional parameters that can be passed. For example:

```go

```

NOTE: Using the [context](https://godoc.org/context) package, one can easily
pass cancelation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then `context.Background()`
can be used as a starting point.

For more sample code snippets, head over to the
[examples](https://github.com/theodesp/go-calendly/tree/master/examples) directory.

### Authentication ###

The go-calendly library does not directly handle authentication. Instead, when
creating a new client, pass an `http.Client` that can handle authentication for
you. If you have an API access token (for example, a [integrations](https://calendly.com/integrations)), you can use it with the library using:

```go

```


### API docs ###

For complete usage of go-calendly, see the full [package docs](https://godoc.org/github.com/theodesp/go-calendly/calendly])

[Calendly API]: https://developer.calendly.com/docs/getting-started

## Roadmap ##

[Contributing](./CONTRIBUTING)

## Versioning ##

In general, go-calendly follows [semver](https://semver.org/) as closely as we
can for tagging releases of the package. For self-contained libraries, the
application of semantic versioning is relatively straightforward and generally
understood. But because go-calendly is a client library for the Calendly API, which
itself changes behavior, and because we are typically pretty aggressive about
implementing preview features of the Calendly API, we've adopted the following
versioning policy:

* We increment the **major version** with any incompatible change to
	non-preview functionality, including changes to the exported Go API surface
	or behavior of the API.
* We increment the **minor version** with any backwards-compatible changes to
	functionality, as well as any changes to preview functionality in the Calendly
	API.
* We increment the **patch version** with any backwards-compatible bug fixes.

Preview functionality may take the form of entire methods or simply additional
data returned from an otherwise non-preview method. Refer to the Calendly API
documentation for details on preview functionality.

## License ##

This library is distributed under the BSD-style license found in the [LICENSE](./LICENSE)
file.