// Package performance provides helpers to get sentry performance monitoring working fast.
package performance

import (
	"context"
	"github.com/getsentry/sentry-go"
)

type transactionContext struct {
	transactions map[string]*sentry.Span
	spans map[string]map[string]*sentry.Span
}

var currentTransactionContext *transactionContext = nil

func (tc *transactionContext) createTransaction(name string, operation string) *sentry.Span {
	var ctx = context.Background()
	var transaction = sentry.StartSpan(ctx, operation, sentry.TransactionName(name))
	tc.transactions[name] = transaction
	tc.spans[name] = make(map[string]*sentry.Span)
	return transaction
}

func (tc *transactionContext) getTransaction(name string) *sentry.Span {
	if transaction, ok := tc.transactions[name]; ok {
		return transaction
	}
	return nil
}

func (tc *transactionContext) createSpan(transactionName string, operation string) *sentry.Span {
	var transaction = tc.getTransaction(transactionName)
	if nil == transaction {
		return nil
	}
	var span = sentry.StartSpan(transaction.Context(), operation)
	tc.spans[transactionName][operation] = span
	return span
}

func (tc *transactionContext) getSpan(transactionName string, operation string) *sentry.Span {
	if transaction, ok := tc.spans[transactionName]; ok {
		if span, ok := transaction[operation]; ok {
			return span
		}
	}
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
		transactions: make(map[string]*sentry.Span),
		spans:        make(map[string]map[string]*sentry.Span),
	}
}

// CreateTransaction creates a transaction that can be attached additional spans.
func CreateTransaction(name string, operation string) *sentry.Span {
	return getCurrentTransactionContext().createTransaction(name, operation)
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
