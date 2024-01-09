package client

import (
	"github.com/rivo/tview"

	"yap-pwkeeper/internal/pkg/models"
)

const (
	pageRoot = "root"
)

type DataStore interface {
	Register(login, password string) error
	Login(login, password string) error
	Update() error
	GetCardsList() []*models.Card
	GetCredentialsList() []*models.Credential
	GetNotesList() []*models.Note
	GetNote(id string) *models.Note
	AddNote(note models.Note) error
	UpdateNote(note models.Note) error
	DeleteNote(note models.Note) error
	GetCard(id string) *models.Card
	AddCard(note models.Card) error
	UpdateCard(note models.Card) error
	DeleteCard(note models.Card) error
	GetCredential(id string) *models.Credential
	AddCredential(note models.Credential) error
	UpdateCredential(note models.Credential) error
	DeleteCredential(note models.Credential) error
}

type App struct {
	ui         *tview.Application
	pages      *tview.Pages
	store      DataStore
	categories *tview.List
	itemsList  *tview.List
	form       *tview.Form
	statusBar  *tview.Form
}

// New is UI app constructor
func New(options ...func(a *App)) *App {
	app := &App{
		ui: tview.NewApplication(),
	}
	for _, opt := range options {
		opt(app)
	}
	app.bootstrap()
	return app
}

func WithDataStore(ds DataStore) func(a *App) {
	return func(a *App) {
		a.store = ds
	}
}

// bootstrap creates all ui pages and should be run only once
func (a *App) bootstrap() {
	a.pages = tview.NewPages()
	a.statusBar = tview.NewForm().SetHorizontal(true)
	a.categories = tview.NewList()
	a.itemsList = tview.NewList()
	a.form = tview.NewForm()
	a.statusBar.AddTextView("Quit: `Esc`, Sync `S`", "", 1, 1, true, false)
	a.ui.SetRoot(a.pages, true).EnableMouse(false)
	a.mainPage()
	a.welcomePage()
}

// Run starts application
func (a *App) Run() error {
	//a.pages.SwitchToPage(pageRoot)
	//a.ui.SetFocus(a.categories)
	return a.ui.Run()
}

// Stop terminates application
func (a *App) Stop() {
	a.ui.Stop()
}
