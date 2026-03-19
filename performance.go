// Package performance provides helpers to get sentry performance monitoring working fast.
package performance

import (
	"context"
	"github.com/getsentry/sentry-go"
	"sync"
)

// Span wraps a sentry.Span and cleans up map references on Finish.
type Span struct {
	*sentry.Span
	transactionName string
	operation       string
	isTransaction   bool
	tc              *transactionContext
}

// Finish finishes the underlying sentry span and removes it from the global map
// so the span (and everything it references) can be garbage collected.
func (s *Span) Finish() {
	if s == nil {
		return
	}
	s.Span.Finish()
	if s.isTransaction {
		s.tc.transactionsMutex.Lock()
		if current, ok := s.tc.transactions[s.transactionName]; ok && current == s {
			delete(s.tc.transactions, s.transactionName)
			s.tc.spansMutex.Lock()
			delete(s.tc.spans, s.transactionName)
			s.tc.spansMutex.Unlock()
		}
		s.tc.transactionsMutex.Unlock()
	} else {
		s.tc.spansMutex.Lock()
		if ops, ok := s.tc.spans[s.transactionName]; ok {
			if current, ok := ops[s.operation]; ok && current == s {
				delete(ops, s.operation)
			}
		}
		s.tc.spansMutex.Unlock()
	}
}

type transactionContext struct {
	transactions      map[string]*Span
	spans             map[string]map[string]*Span
	transactionsMutex sync.RWMutex
	spansMutex        sync.RWMutex
}

var currentTransactionContext *transactionContext = nil

func (tc *transactionContext) createTransaction(name string, operation string, ctx context.Context) *Span {
	var sentrySpan = sentry.StartSpan(ctx, operation, sentry.WithTransactionName(name))
	span := &Span{
		Span:            sentrySpan,
		transactionName: name,
		isTransaction:   true,
		tc:              tc,
	}
	tc.transactionsMutex.Lock()
	tc.transactions[name] = span
	tc.transactionsMutex.Unlock()
	tc.spansMutex.Lock()
	tc.spans[name] = make(map[string]*Span)
	tc.spansMutex.Unlock()
	return span
}

func (tc *transactionContext) getTransaction(name string) *Span {
	tc.transactionsMutex.RLock()
	if transaction, ok := tc.transactions[name]; ok {
		tc.transactionsMutex.RUnlock()
		return transaction
	}
	tc.transactionsMutex.RUnlock()
	return nil
}

func (tc *transactionContext) createSpan(transactionName string, operation string) *Span {
	var transaction = tc.getTransaction(transactionName)
	if nil == transaction {
		transaction = tc.createTransaction(transactionName, operation, context.Background())
	}
	var sentrySpan = sentry.StartSpan(transaction.Context(), operation)
	span := &Span{
		Span:            sentrySpan,
		transactionName: transactionName,
		operation:       operation,
		isTransaction:   false,
		tc:              tc,
	}
	tc.spansMutex.Lock()
	if ops, ok := tc.spans[transactionName]; ok {
		ops[operation] = span
	}
	tc.spansMutex.Unlock()
	return span
}

func (tc *transactionContext) getSpan(transactionName string, operation string) *Span {
	tc.spansMutex.RLock()
	if transaction, ok := tc.spans[transactionName]; ok {
		if span, ok := transaction[operation]; ok {
			tc.spansMutex.RUnlock()
			return span
		}
	}
	tc.spansMutex.RUnlock()
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
		transactions:      make(map[string]*Span),
		spans:             make(map[string]map[string]*Span),
		transactionsMutex: sync.RWMutex{},
		spansMutex:        sync.RWMutex{},
	}
}

// CreateTransaction creates a transaction that can be attached additional spans.
func CreateTransaction(name string, operation string) *Span {
	return getCurrentTransactionContext().createTransaction(name, operation, context.Background())
}

// CreateTransactionWithContext creates a transaction that can be attached additional spans and sets an existing context.
func CreateTransactionWithContext(name string, operation string, ctx context.Context) *Span {
	return getCurrentTransactionContext().createTransaction(name, operation, ctx)
}

// GetTransaction returns a transaction by its name (or nil if none exists).
func GetTransaction(name string) *Span {
	return getCurrentTransactionContext().getTransaction(name)
}

// CreateSpan creates a span and attaches it to a transaction specified by its name (returns nil if transaction doesn't exist).
func CreateSpan(transactionName string, operation string) *Span {
	return getCurrentTransactionContext().createSpan(transactionName, operation)
}

// GetSpan returns a transaction by its name (or nil if none exists).
func GetSpan(transactionName string, operation string) *Span {
	return getCurrentTransactionContext().getSpan(transactionName, operation)
}
