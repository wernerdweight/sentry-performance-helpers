// Package performance provides helpers to get sentry performance monitoring working fast.
package performance

import (
	"context"
	"github.com/getsentry/sentry-go"
	"sync"
)

type transactionContext struct {
	transactions      map[string]*sentry.Span
	spans             map[string]map[string]*sentry.Span
	transactionsMutex sync.RWMutex
	spansMutex        sync.RWMutex
}

var currentTransactionContext *transactionContext = nil

func (tc *transactionContext) createTransaction(name string, operation string, ctx context.Context) *sentry.Span {
	var transaction = sentry.StartSpan(ctx, operation, sentry.WithTransactionName(name))
	tc.transactionsMutex.Lock()
	tc.transactions[name] = transaction
	tc.transactionsMutex.Unlock()
	tc.spansMutex.Lock()
	tc.spans[name] = make(map[string]*sentry.Span)
	tc.spansMutex.Unlock()
	return transaction
}

func (tc *transactionContext) getTransaction(name string) *sentry.Span {
	tc.transactionsMutex.Lock()
	if transaction, ok := tc.transactions[name]; ok {
		tc.transactionsMutex.Unlock()
		return transaction
	}
	tc.transactionsMutex.Unlock()
	return nil
}

func (tc *transactionContext) createSpan(transactionName string, operation string) *sentry.Span {
	var transaction = tc.getTransaction(transactionName)
	if nil == transaction {
		transaction = tc.createTransaction(transactionName, operation, context.Background())
	}
	var span = sentry.StartSpan(transaction.Context(), operation)
	tc.spansMutex.Lock()
	tc.spans[transactionName][operation] = span
	tc.spansMutex.Unlock()
	return span
}

func (tc *transactionContext) getSpan(transactionName string, operation string) *sentry.Span {
	tc.spansMutex.Lock()
	if transaction, ok := tc.spans[transactionName]; ok {
		if span, ok := transaction[operation]; ok {
			tc.spansMutex.Unlock()
			return span
		}
	}
	tc.spansMutex.Unlock()
	return nil
}

func getCurrentTransactionContext() *transactionContext {
	if nil == currentTransactionContext {
		Refresh()
	}
	return currentTransactionContext
}

// Refresh clears current transaction context.
func Refresh() {
	currentTransactionContext = &transactionContext{
		transactions:      make(map[string]*sentry.Span),
		spans:             make(map[string]map[string]*sentry.Span),
		transactionsMutex: sync.RWMutex{},
		spansMutex:        sync.RWMutex{},
	}
}

// CreateTransaction creates a transaction that can be attached additional spans.
func CreateTransaction(name string, operation string) *sentry.Span {
	return getCurrentTransactionContext().createTransaction(name, operation, context.Background())
}

// CreateTransactionWithContext creates a transaction that can be attached additional spans and sets an existing context.
func CreateTransactionWithContext(name string, operation string, ctx context.Context) *sentry.Span {
	return getCurrentTransactionContext().createTransaction(name, operation, ctx)
}

// GetTransaction returns a transaction by its name (or nil if none exists).
func GetTransaction(name string) *sentry.Span {
	return getCurrentTransactionContext().getTransaction(name)
}

// CreateSpan creates a span and attaches it to a transaction specified by its name (returns nil if transaction doesn't exist).
func CreateSpan(transactionName string, operation string) *sentry.Span {
	return getCurrentTransactionContext().createSpan(transactionName, operation)
}

// GetSpan returns a transaction by its name (or nil if none exists).
func GetSpan(transactionName string, operation string) *sentry.Span {
	return getCurrentTransactionContext().getSpan(transactionName, operation)
}
