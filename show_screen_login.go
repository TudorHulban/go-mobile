package main

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *application) showLoginScreen(window fyne.Window) {
	titleLabel := widget.NewLabel("Account Login")
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	loginButton := widget.NewButton(
		"Sign In",
		func() {
			if usernameEntry.Text == "admin" && passwordEntry.Text == "password" {
				// Save user preference dynamically to device filesystem
				fyne.CurrentApp().
					Preferences().
					SetString("session_user", usernameEntry.Text)

				a.showDashboardScreen(window, usernameEntry.Text)
			} else {
				dialog.ShowError(
					errors.New("Invalid username or password"),
					window,
				)
			}
		},
	)

	loginContent := container.NewVBox(
		titleLabel,
		usernameEntry,
		passwordEntry,
		loginButton,
	)

	window.SetContent(loginContent)
}
