package main

import (
	"context"
	"errors"
	"sync"
)

type CloserFunc func(ctx context.Context) error

type Graceful struct {
	mu *sync.Mutex
	closers []CloserFunc
}

func NewGraceful() Graceful {
	return Graceful{
		mu: &sync.Mutex{},
		closers: make([]CloserFunc, 0, 1),
	}
}

func (g *Graceful) Add(closerFunc CloserFunc) {
	g.mu.Lock()
	g.closers = append(g.closers, closerFunc)
	g.mu.Unlock()
}

func (g *Graceful) startGraceful() {
	g.mu.Lock()
}

func (g *Graceful) StartGraceful(shutdownCtx context.Context) ([]error) {	
	errs := make([]error, 0, 1)

	select {
	case <- shutdownCtx.Done():
		err := errors.New("shutdown ctx canceled")
		errs = append(errs, err)
		return errs
	default:
	}

	for _, close := range g.closers {
		select {
		case <- shutdownCtx.Done():
			err := errors.New("shutdown ctx canceled")
			errs = append(errs, err)
			return errs
		default:
		}

		err := close(shutdownCtx)
		if err != nil {
			errs  = append(errs, err)
		}
	} 

	return errs
}