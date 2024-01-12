package client

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"yap-pwkeeper/internal/pkg/models"
)

const (
	pageWelcome = "welcome"
	pageLogin   = "cred"
)

// welcomePage creates welcome page and switches to it
func (a *App) welcomePage() {
	p := tview.NewModal().SetBackgroundColor(tcell.ColorBlack)
	p.SetText("Welcome!")
	p.AddButtons([]string{"Register New User", "Login Existing User"})
	p.SetFocus(1)
	p.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonIndex {
		case 0:
			a.registerPage()
		case 1:
			a.loginPage()
		default:
			a.Stop()
		}
		a.pages.RemovePage(pageWelcome)
	})
	a.pages.AddPage(pageWelcome, p, true, true)
	a.pages.SwitchToPage(pageWelcome)
}

// loginPage user login page
func (a *App) loginPage() {
	login := models.UserCredentials{}
	cancelFunc := func() {
		a.welcomePage()
		a.pages.RemovePage(pageLogin)
	}
	form := tview.NewForm().SetCancelFunc(cancelFunc)
	form.AddInputField("login", "", 30, nil, func(text string) {
		login.Login = text
	})
	form.AddPasswordField("password", "", 30, '*', func(text string) {
		login.Password = text
	})
	form.AddButton("login", func() {
		if err := a.store.Login(login.Login, login.Password); err != nil {
			a.modalErr(err.Error())
			return
		}
		a.pages.RemovePage(pageLogin)
		a.pages.SwitchToPage(pageRoot)
		a.ui.SetFocus(a.categories)

	})
	form.AddButton("<<Back", func() {
		cancelFunc()
	})

	form.SetButtonsAlign(tview.AlignCenter)
	form.SetBorder(true)
	form.SetTitle("Enter login credentials")
	//return loginForm
	a.pages.AddPage(pageLogin, center(44, 9, form), true, true)
	a.pages.SwitchToPage(pageLogin)
}

// loginPage user registration page
func (a *App) registerPage() {
	login := models.UserCredentials{}
	rePass := ""
	cancelFunc := func() {
		a.welcomePage()
		a.pages.RemovePage(pageLogin)
	}
	form := tview.NewForm().SetCancelFunc(cancelFunc)
	form.AddInputField("login", "", 30, nil, func(text string) {
		login.Login = text
	})
	form.AddPasswordField("password", "", 30, '*', func(text string) {
		login.Password = text
	})
	form.AddPasswordField("repeat password", "", 30, '*', func(text string) {
		rePass = text
	})
	form.AddButton("Register", func() {
		if rePass != login.Password {
			a.modalErr("Passwords do not match!")
			return
		}
		if len(login.Login) < 3 {
			a.modalErr("login is too short, use at least 3 symbols")
			return
		}
		if err := a.store.Register(login.Login, login.Password); err != nil {
			a.modalErr("Registration failed: " + err.Error())
			return
		}
		if err := a.store.Login(login.Login, login.Password); err != nil {
			a.modalErr("User registered, but: " + err.Error())
			return
		}
		a.pages.RemovePage(pageLogin)
		a.pages.SwitchToPage(pageRoot)
		a.ui.SetFocus(a.categories)

	})
	form.AddButton("<<Back", func() {
		cancelFunc()
	})
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetBorder(true)
	form.SetTitle("New User Registration")
	a.pages.AddPage(pageLogin, center(51, 11, form), true, true)
	a.pages.SwitchToPage(pageLogin)
}
