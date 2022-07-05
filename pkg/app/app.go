package app

import (
	"context"

	"go.uber.org/fx"
)

type App struct {
	options []fx.Option
}

func (a *App) Provides(constructors ...interface{}) *App {
	a.options = append(a.options, fx.Provide(constructors...))
	return a
}

func (a *App) Populates(targets ...interface{}) *App {
	a.options = append(a.options, fx.Populate(targets...))
	return a
}

func (a *App) Supply(values ...interface{}) *App {
	a.options = append(a.options, fx.Supply(values...))
	return a
}

func (a *App) Invoke(invokes ...interface{}) *App {
	a.options = append(a.options, fx.Invoke(invokes...))
	return a
}

func (a *App) Options(options ...fx.Option) *App {
	a.options = append(a.options, options...)
	return a
}

func (a *App) Logger(constructor interface{}) *App {
	a.options = append(a.options, fx.WithLogger(constructor))
	return a
}

func (a App) Start() error {
	return fx.New(
		a.options...,
	).Start(context.Background())
}
