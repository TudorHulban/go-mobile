package main

import (
	// Added standard Go errors package

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Fyne App")

	// Start by showing the login screen
	showLoginScreen(myWindow)

	myWindow.ShowAndRun()
}
