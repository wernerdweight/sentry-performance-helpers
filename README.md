Performance Helpers for Sentry Performance Monitoring
====================================

Helpers to get sentry performance monitoring working easily.

[![Build Status](https://travis-ci.org/wernerdweight/sentry-performance-helpers.svg?branch=master)](https://travis-ci.org/wernerdweight/sentry-performance-helpers)
[![Latest Stable Version](https://poser.pugx.org/wernerdweight/sentry-performance-helpers/v/stable)](https://packagist.org/packages/wernerdweight/sentry-performance-helpers)
[![Total Downloads](https://poser.pugx.org/wernerdweight/sentry-performance-helpers/downloads)](https://packagist.org/packages/wernerdweight/sentry-performance-helpers)
[![License](https://poser.pugx.org/wernerdweight/sentry-performance-helpers/license)](https://packagist.org/packages/wernerdweight/sentry-performance-helpers)


Installation
------------

### 1. Download using composer

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
