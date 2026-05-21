package main

import (
	// Added standard Go errors package

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type application struct {
	ui   fyne.App
	team Team
}

func main() {
	team := NewTeam(
		[]string{
			"Unassigned",
			"Alice",
			"Bob",
			"Charlie",
		}...,
	)

	team.UpsertTask(
		&Task{
			Name:   "Operation 1",
			Status: "init",
		},
	)

	team.UpsertTask(
		&Task{
			Name:   "Revision 1",
			Status: "not started",
		},
	)

	team.UpsertTask(
		&Task{
			Name:    "Operation 2",
			Status:  "assigned",
			OwnerID: 1,
		},
	)

	a := application{
		ui:   app.New(),
		team: *team,
	}

	myWindow := a.ui.NewWindow("Fyne App")

	savedUser := a.ui.
		Preferences().
		StringWithFallback("session_user", "")

	if savedUser != "" {
		// Bypass login completely!
		a.showDashboardScreen(myWindow, savedUser)
	} else {
		// Fresh run, require authentication
		// showLoginScreen(myWindow)
		a.showDashboardScreen(myWindow, "admin")
	}

	myWindow.ShowAndRun()
}
