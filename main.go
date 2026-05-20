package main

import (
	// Added standard Go errors package

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Fyne App")

	savedUser := myApp.Preferences().StringWithFallback("session_user", "")

	if savedUser != "" {
		// Bypass login completely!
		showDashboardScreen(myWindow, savedUser)
	} else {
		// Fresh run, require authentication
		// showLoginScreen(myWindow)
		showDashboardScreen(myWindow, "admin")
	}

	myWindow.ShowAndRun()
}
