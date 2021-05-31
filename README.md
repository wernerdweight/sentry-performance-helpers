Performance Helpers for Sentry Performance Monitoring
====================================

Helpers to get sentry performance monitoring working easily.

[![Build Status](https://www.travis-ci.com/wernerdweight/sentry-performance-helpers.svg?branch=master)](https://www.travis-ci.com/wernerdweight/sentry-performance-helpers)
[![Go Report Card](https://goreportcard.com/badge/github.com/wernerdweight/sentry-performance-helpers)](https://goreportcard.com/report/github.com/wernerdweight/sentry-performance-helpers)
[![GoDoc](https://godoc.org/github.com/wernerdweight/sentry-performance-helpers?status.svg)](https://godoc.org/github.com/wernerdweight/sentry-performance-helpers)
[![go.dev](https://img.shields.io/badge/go.dev-pkg-007d9c.svg?style=flat)](https://pkg.go.dev/github.com/wernerdweight/sentry-performance-helpers)


Installation
------------

### 1. Installation

```bash
go get github.com/wernerdweight/sentry-performance-helpers
```

Configuration
------------

The package itself needs no configuration. Check Sentry [documentation](https://docs.sentry.io/platforms/go/performance/) for the setup related to Sentry.

Usage
------------

**Basic usage**

```go
// create a transaction
transaction := performance.CreateTransaction("transaction-name", "operation")
doSomething()
// create a span attached to the transaction
span := performance.CreateSpan("transaction-name", "operation")
doSomethingElse()
// finish span (otherwise, it will not get to Sentry)
span.Finish()
doSomethingCompletelyDifferent()
// finish transaction (otherwise, it will not get to Sentry)
transaction.Finish()
```

**Accessing transaction/span properties**

```go
// get transaction by its name
transaction := performance.GetTransaction("transaction-name")
// get span by operation
span := performance.GetSpan("transaction-name", "operation")
```

**Usage with defer**

```go
package main

import "github.com/wernerdweight/sentry-performance-helpers"

func doSomething() {
    defer performance.CreateSpan("my-transaction", "doSomething").Finish()
    // put your code here, the transaction/span is created now
    // due to how defer works, the 'Finish' method will be called at the end
}

func main() {
    defer performance.CreateTransaction("my-transaction", "main").Finish()
    doSomething()
}
```

License
-------
This package is under the MIT license. See the complete license in the root directiory of the bundle.
