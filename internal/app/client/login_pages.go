package client

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"yap-pwkeeper/internal/pkg/models"
)

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
		a.pages.RemovePage("welcome")
	})
	a.pages.AddPage("welcome", p, true, true)
	a.pages.SwitchToPage("welcome")
}

func (a *App) loginPage() {
	login := models.UserCredentials{}
	cFunc := func() {
		a.welcomePage()
		a.pages.RemovePage("login")
	}
	form := tview.NewForm().SetCancelFunc(cFunc)
	form.AddInputField("login", "", 30, nil, func(text string) {
		login.Login = text
	})
	form.AddPasswordField("password", "", 30, '*', func(text string) {
		login.Password = text
	})
	form.AddButton("login", func() {
		if err := a.authServer.Login(login.Login, login.Password); err != nil {
			a.modalErr("login failed: " + err.Error())
			return
		}
		a.browser()
	})
	form.AddButton("<<Back", func() {
		cFunc()
	})
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetBorder(true)
	form.SetTitle("Enter login credentials")
	//return loginForm
	a.pages.AddPage("login", center(44, 9, form), true, true)
	a.pages.SwitchToPage("login")
}

func (a *App) registerPage() {
	login := models.UserCredentials{}
	rePass := ""
	cFunc := func() {
		a.welcomePage()
		a.pages.RemovePage("register")
	}
	form := tview.NewForm().SetCancelFunc(cFunc)
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
		if err := a.authServer.Register(login.Login, login.Password); err != nil {
			a.modalErr("Registration failed: " + err.Error())
			return
		}
		a.browser()
	})
	form.AddButton("<<Back", func() {
		cFunc()
	})
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetBorder(true)
	form.SetTitle("New User Registration")
	a.pages.AddPage("register", center(51, 11, form), true, true)
	a.pages.SwitchToPage("register")
}
