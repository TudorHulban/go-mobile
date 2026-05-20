package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func showDashboardScreen(window fyne.Window) {
	welcomeLabel := widget.NewLabel("Welcome to your clean Dashboard! 🚀")

	logoutBtn := widget.NewButton(
		"Log Out",
		func() {
			showLoginScreen(window)
		},
	)

	dashboardContent := container.NewVBox(
		welcomeLabel,
		logoutBtn,
	)

	window.SetContent(dashboardContent)
}
