package main

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func showLoginScreen(window fyne.Window) {
	titleLabel := widget.NewLabel("Account Login")
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	loginButton := widget.NewButton(
		"Sign In",
		func() {
			if usernameEntry.Text == "admin" && passwordEntry.Text == "password" {
				showDashboardScreen(window)
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
