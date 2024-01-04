package client

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"yap-pwkeeper/internal/pkg/models"
)

func (a *App) welcomePage() {
	welcome := tview.NewModal().SetBackgroundColor(tcell.ColorBlack)
	welcome.SetText("Welcome!")
	welcome.AddButtons([]string{"Register New User", "Login Existing User"})
	welcome.SetFocus(1)
	welcome.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonIndex {
		case 0:
			a.registerPage()
		case 1:
			a.loginPage()
		default:
			a.Stop()
		}
	})
	a.pages.AddPage("welcome", welcome, true, true)
	a.pages.SwitchToPage("welcome")
}

func (a *App) loginPage() {
	login := models.LoginRequest{}
	loginForm := tview.NewForm()
	loginForm.AddInputField("login", "", 30, nil, func(text string) {
		login.Login = text
	})
	loginForm.AddPasswordField("password", "", 30, '*', func(text string) {
		login.Password = text
	})
	loginForm.AddButton("Login", func() {})
	loginForm.AddButton("Back", func() {
		a.pages.SwitchToPage("welcome")
		a.pages.RemovePage("login")
	})
	loginForm.SetButtonsAlign(tview.AlignCenter)
	loginForm.SetBorder(true)
	loginForm.SetTitle("Enter login credentials")
	a.pages.AddPage("login", center(44, 9, loginForm), true, true)
	a.pages.SwitchToPage("login")
}

func (a *App) registerPage() {
	login := models.LoginRequest{}
	rePass := ""
	registerForm := tview.NewForm()
	registerForm.AddInputField("login", "", 30, nil, func(text string) {
		login.Login = text
	})
	registerForm.AddPasswordField("password", "", 30, '*', func(text string) {
		login.Password = text
	})
	registerForm.AddPasswordField("repeat password", "", 30, '*', func(text string) {
		rePass = text
	})
	registerForm.AddButton("Register", func() {
		if rePass != login.Password {
			a.modalErr("Passwords do not match!")
			return
		}
		if len(login.Login) < 3 {
			a.modalErr("Login does not match criteria:")
			return
		}

	})
	registerForm.AddButton("Back", func() {
		a.pages.SwitchToPage("welcome")
		a.pages.RemovePage("register")
	})
	registerForm.SetButtonsAlign(tview.AlignCenter)
	registerForm.SetBorder(true)
	registerForm.SetTitle("Register New User")
	a.pages.AddPage("register", center(51, 11, registerForm), true, true)
	a.pages.SwitchToPage("register")
}

func (a *App) modalErr(text string) {
	m := tview.NewModal().
		SetBackgroundColor(tcell.ColorRed).
		SetText("Error: " + text).
		AddButtons([]string{"Got it"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("error")
		})
	a.pages.AddPage("error", m, true, true)
}

func center(width, height int, p tview.Primitive) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}
