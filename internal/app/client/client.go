package client

import (
	"github.com/rivo/tview"

	"yap-pwkeeper/internal/pkg/models"
)

const (
	logsPage = "logsPage"
	mainPage = "finder"
)

type AuthServer interface {
	Register(login, password string) error
	Login(login, password string) error
}

type DataStore interface {
	Clear()
	Update() error
	GetCardsList() []*models.Card
	GetCredentialsList() []*models.Credential
	GetNotesList() []*models.Note
	GetNote(id string) *models.Note
	AddNote(note models.Note) error
	UpdateNote(note models.Note) error
	DeleteNote(note models.Note) error
	GetCard(id string) *models.Card
	GetCredential(id string) *models.Credential
}

type App struct {
	ui         *tview.Application
	pages      *tview.Pages
	authServer AuthServer
	debug      bool
	store      DataStore
}

func New(options ...func(a *App)) *App {
	app := &App{
		ui:    tview.NewApplication(),
		pages: tview.NewPages(),
	}
	for _, opt := range options {
		opt(app)
	}
	app.ui.SetRoot(app.pages, true).EnableMouse(false)
	return app
}

func WithAuthServer(aaa AuthServer) func(a *App) {
	return func(a *App) {
		a.authServer = aaa
	}
}
func WithDataStore(ds DataStore) func(a *App) {
	return func(a *App) {
		a.store = ds
	}
}
func WithDebug(level int) func(a *App) {
	return func(a *App) {
		if level > 0 {
			a.debug = true
		}
	}
}

func (a *App) Run() error {
	//log.SetOutput(io.Discard)
	//a.welcomePage()
	a.browser()
	return a.ui.Run()
}

func (a *App) Stop() {
	a.ui.Stop()
}
