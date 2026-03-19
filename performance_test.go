package performance_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	performance "github.com/wernerdweight/sentry-performance-helpers"
	"testing"
)

func TestRefresh(t *testing.T) {
	assertion := assert.New(t)
	performance.CreateTransaction("test", "test")
	performance.CreateSpan("test", "test")
	assertion.NotNil(performance.GetTransaction("test"))
	assertion.NotNil(performance.GetSpan("test", "test"))
	performance.Refresh()
	assertion.Nil(performance.GetTransaction("test"))
	assertion.Nil(performance.GetSpan("test", "test"))
}

func TestTransactions(t *testing.T) {
	assertion := assert.New(t)
	performance.CreateTransaction("test", "test")
	assertion.NotNil(performance.GetTransaction("test"))
	assertion.Nil(performance.GetTransaction("unknown"))
}

func TestSpans(t *testing.T) {
	assertion := assert.New(t)
	performance.CreateTransaction("test", "test")
	assertion.NotNil(performance.CreateSpan("test", "test"))
	assertion.NotNil(performance.GetSpan("test", "test"))
	assertion.Nil(performance.GetSpan("test", "unknown"))
	assertion.Nil(performance.GetSpan("unknown", "unknown"))
	assertion.NotNil(performance.CreateSpan("unknown", "unknown"))
}

func TestFinishCleansUpTransaction(t *testing.T) {
	assertion := assert.New(t)
	span := performance.CreateTransaction("cleanup-test", "test.op")
	assertion.NotNil(performance.GetTransaction("cleanup-test"))

	span.Finish()

	assertion.Nil(performance.GetTransaction("cleanup-test"))
}

func TestFinishCleansUpTransactionWithContext(t *testing.T) {
	assertion := assert.New(t)
	span := performance.CreateTransactionWithContext("ctx-test", "test.op", context.Background())
	assertion.NotNil(performance.GetTransaction("ctx-test"))

	span.Finish()

	assertion.Nil(performance.GetTransaction("ctx-test"))
	assertion.Nil(performance.GetSpan("ctx-test", "any"))
}

func TestFinishCleansUpSpan(t *testing.T) {
	assertion := assert.New(t)
	performance.CreateTransaction("span-test", "test.op")
	span := performance.CreateSpan("span-test", "child.op")
	assertion.NotNil(performance.GetSpan("span-test", "child.op"))

	span.Finish()

	assertion.Nil(performance.GetSpan("span-test", "child.op"))
	assertion.NotNil(performance.GetTransaction("span-test"))
}

func TestFinishTransactionCleansUpAllSpans(t *testing.T) {
	assertion := assert.New(t)
	tx := performance.CreateTransaction("full-cleanup", "test.op")
	performance.CreateSpan("full-cleanup", "child1")
	performance.CreateSpan("full-cleanup", "child2")
	assertion.NotNil(performance.GetSpan("full-cleanup", "child1"))
	assertion.NotNil(performance.GetSpan("full-cleanup", "child2"))

	tx.Finish()

	assertion.Nil(performance.GetTransaction("full-cleanup"))
	assertion.Nil(performance.GetSpan("full-cleanup", "child1"))
	assertion.Nil(performance.GetSpan("full-cleanup", "child2"))
}

func TestSpanEmbedsSentrySpan(t *testing.T) {
	assertion := assert.New(t)
	span := performance.CreateTransaction("embed-test", "test.op")
	assertion.NotNil(span.Context())
	span.Finish()
}
