package client

import (
	"github.com/rivo/tview"

	"yap-pwkeeper/internal/pkg/models"
)

type Auther interface {
	Register(request models.LoginRequest) (models.AuthToken, error)
	Login(request models.LoginRequest) (models.AuthToken, error)
	Refresh(token models.AuthToken) (models.AuthToken, error)
}

type App struct {
	ui      *tview.Application
	pages   *tview.Pages
	token   models.AuthToken
	authAPI Auther
	data    models.Wallet
}

func New(options ...func(*App)) *App {
	app := &App{
		ui:    tview.NewApplication(),
		pages: tview.NewPages(),
		data:  models.Store,
	}
	for _, opt := range options {
		opt(app)
	}
	app.ui.SetRoot(app.pages, true).EnableMouse(true)
	return app
}

func (a *App) Run() error {
	a.welcomePage()
	//a.pages.SwitchToPage("welcome")
	return a.ui.Run()
}

func (a *App) Stop() {
	a.ui.Stop()
}
