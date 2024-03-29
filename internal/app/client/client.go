package client

import (
	"github.com/rivo/tview"

	"yap-pwkeeper/internal/pkg/models"
)

const (
	pageRoot = "root"
)

// DataStore defines all store methods
type DataStore interface {
	Register(login, password string) error
	Login(login, password string) error

	Update() error

	GetCardsList() []*models.Card
	GetCard(id string) *models.Card
	AddCard(note models.Card) error
	UpdateCard(note models.Card) error
	DeleteCard(note models.Card) error

	GetCredentialsList() []*models.Credential
	GetCredential(id string) *models.Credential
	AddCredential(note models.Credential) error
	UpdateCredential(note models.Credential) error
	DeleteCredential(note models.Credential) error

	GetNotesList() []*models.Note
	GetNote(id string) *models.Note
	AddNote(note models.Note) error
	UpdateNote(note models.Note) error
	DeleteNote(note models.Note) error

	GetFilesList() []*models.File
	GetFileInfo(id string) *models.File
	GetFile(documentId string, path string) error
	AddFile(d models.File, filename string) error
	UpdateFileInfo(d models.File) error
	UpdateFile(d models.File, filename string) error
	DeleteFile(note models.File) error
}

type App struct {
	ui         *tview.Application
	pages      *tview.Pages
	store      DataStore
	categories *tview.List
	itemsList  *tview.List
	form       *tview.Form
	statusBar  *tview.Form
	useMouse   bool
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

// WithDataStore attaches storage to app instance
func WithDataStore(ds DataStore) func(a *App) {
	return func(a *App) {
		a.store = ds
	}
}

// WithMouse enables mouse in tui (may be not stable)
func WithMouse(m bool) func(a *App) {
	return func(a *App) {
		a.useMouse = m
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
	a.ui.SetRoot(a.pages, true).EnableMouse(a.useMouse)
	a.mainPage()
	a.welcomePage()
}

// Run starts application
func (a *App) Run() error {
	return a.ui.Run()
}

// Stop terminates application
func (a *App) Stop() {
	a.ui.Stop()
}
