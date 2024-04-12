package main

import (
	"bufio"
	"changeme/tools/SearchRegistry"
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows/registry"
	"os"
	"strings"
	"sync"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx

	go func() {
		for {
			select {
			case percentage := <-SearchRegistry.Percentage:
				runtime.EventsEmit(ctx, "percentage", percentage)
			case _searcher := <-SearchRegistry.SearchChan:
				runtime.EventsEmit(ctx, "SearchRegistry", _searcher)
				//case _error := <-SearchRegistry.Error:
				//	runtime.EventsEmit(ctx, "DeleteError", _error)
			}
		}
	}()

}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) SearchRegistry(input string) {
	SearchRegistry.SearchRegistry(input)
}

func (a *App) DeleteRegistry(input string) string {
	result := SearchRegistry.DeleteRegistry(input)
	return result
}



