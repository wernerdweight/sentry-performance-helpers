package performance_test

import (
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
	assertion.Nil(performance.CreateSpan("unknown", "unknown"))
}
