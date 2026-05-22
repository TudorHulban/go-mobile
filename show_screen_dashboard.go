package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type stateDashboard struct {
	selectorOwner  *widget.Select
	selectorStatus *widget.Select

	btnUpsert *widget.Button

	currentTaskID uint8
	taskSelected  bool
}

func (a *application) checkIfModified(state *stateDashboard) {
	if !state.taskSelected {
		return
	}

	currentTask, errGetTask := a.team.GetTaskBy(state.currentTaskID)
	if errGetTask != nil {
		return
	}

	assignedTeamMember, errGetPerson := a.team.GetTeamMemberBy(currentTask.OwnerID)
	if errGetPerson != nil {
		assignedTeamMember = "" // Handle case where task has no owner yet or owner missing safely

	}

	hasOwnerChanged := state.selectorOwner.Selected != assignedTeamMember
	hasStatusChanged := state.selectorStatus.Selected != currentTask.Status

	if hasOwnerChanged || hasStatusChanged {
		state.btnUpsert.Enable()
	} else {
		state.btnUpsert.Disable()
	}
}

func (a *application) buildDetailPanel(state *stateDashboard) fyne.CanvasObject {
	state.selectorOwner = widget.NewSelect(
		a.team.GetTeamMembers(),
		nil,
	)

	state.selectorStatus = widget.NewSelect(
		[]string{"init", "in work", "closed"},
		nil,
	)

	state.btnUpsert = widget.NewButton(
		"Save Changes",
		nil,
	)

	state.selectorOwner.Hide()
	state.selectorStatus.Hide()
	state.btnUpsert.Disable()
	state.btnUpsert.Hide()

	return container.NewVBox(
		widget.NewLabel("Task Details"),
		container.NewHBox(widget.NewLabel("Owner:"), state.selectorOwner),
		container.NewHBox(widget.NewLabel("Status:"), state.selectorStatus),
		state.btnUpsert,
	)
}

type TaskRow struct {
	fyne.CanvasObject
	name   *widget.Label
	status *canvas.Text
}

func (a *application) showDashboardScreen(window fyne.Window, username string) {
	greetingLabel := widget.NewLabel("Hi, " + username + " 👋")
	greetingLabel.TextStyle = fyne.TextStyle{Bold: true}

	detailTitle := widget.NewLabel("Select a task to view details")

	// Initialize state object ONCE
	state := stateDashboard{}

	onFieldChanged := func(_ string) { // The string parameter is the new value, but we can ignore it since checkIfModified reads from the widgets directly
		a.checkIfModified(&state)
	}

	detailPanel := a.buildDetailPanel(&state)

	// Assign the same callback to both selectors
	state.selectorOwner.OnChanged = onFieldChanged
	state.selectorStatus.OnChanged = onFieldChanged

	// =========================================================
	// Task List
	// =========================================================
	var taskList *widget.List

	taskList = widget.NewList(
		func() int {
			return len(a.team.TeamTasks)
		},

		func() fyne.CanvasObject {
			taskName := widget.NewLabel("Template Task")
			statusBadge := canvas.NewText("STATUS", color.White)
			statusBadge.TextStyle = fyne.TextStyle{Bold: true}

			return container.NewHBox(taskName, statusBadge)
		},

		func(id widget.ListItemID, item fyne.CanvasObject) {
			// Map Fyne's sequence ID to your custom Map uint8 keys
			taskIDs := a.team.GetTasksIDs()
			if int(id) >= len(taskIDs) {
				return
			}
			actualID := taskIDs[id]

			currentTask, errGetTask := a.team.GetTaskBy(actualID)
			if errGetTask != nil {
				return
			}

			box := item.(*fyne.Container)
			labelName := box.Objects[0].(*widget.Label)
			badgeText := box.Objects[1].(*canvas.Text)

			labelName.SetText(currentTask.Name)
			badgeText.Text = "  " + currentTask.Status + "  "
			badgeText.Color = getStatusColor(currentTask.Status)

			badgeText.Refresh()
		},
	)

	// =========================================================
	// Save Button
	// =========================================================
	state.btnUpsert.OnTapped = func() {
		// 1. Check the shared state pointer
		if !state.taskSelected {
			return
		}

		// Guard against double-clicks or mid-process changes
		state.taskSelected = false

		currentTask, errGetTask := a.team.GetTaskBy(state.currentTaskID)
		if errGetTask != nil {
			state.taskSelected = true // Reset guard before leaving
			return
		}

		// Look up the selected string name from the dropdown back into its uint8 ID
		var newOwnerID uint8
		for memberID, name := range a.team.TeamMembers {
			if name == state.selectorOwner.Selected {
				newOwnerID = memberID
				break
			}
		}

		// Update the real data objects
		currentTask.OwnerID = newOwnerID
		currentTask.Status = state.selectorStatus.Selected

		// UI Feedback: immediate disable while processing/unselecting
		state.btnUpsert.Disable()

		// Refresh using Fyne's local row selection index
		taskIDs := a.team.GetTasksIDs()

		for rowIdx, dbID := range taskIDs {
			if dbID == state.currentTaskID {
				taskList.RefreshItem(rowIdx)

				// Note: Unselect will clear the panel, so if your UI hides
				// things on unselect, that logic will execute here.
				taskList.Unselect(rowIdx)

				break
			}
		}

		a.team.UpsertTask(currentTask)
	}

	// =========================================================
	// Row Selection
	// =========================================================
	taskList.OnSelected = func(id widget.ListItemID) {
		taskIDs := a.team.GetTasksIDs()
		if int(id) >= len(taskIDs) {
			return
		}

		// Set this to false FIRST.
		// This acts as a guard. When SetSelected fires OnChanged below,
		// checkIfModified will see this is false and safely return early.
		state.selectorOwner.OnChanged = nil
		state.selectorStatus.OnChanged = nil
		state.taskSelected = false

		state.currentTaskID = taskIDs[id]

		selectedTask, errGetTask := a.team.GetTaskBy(state.currentTaskID)
		if errGetTask != nil {
			return
		}

		detailTitle.SetText("Task - " + selectedTask.Name)

		// Get owner text string from their uint8 ID
		ownerName, _ := a.team.GetTeamMemberBy(selectedTask.OwnerID)

		// Update the selector components programmatically
		state.selectorOwner.SetSelected(ownerName)
		state.selectorStatus.SetSelected(selectedTask.Status)

		// Show the components
		state.selectorOwner.Show()
		state.selectorStatus.Show()
		state.btnUpsert.Show()

		// Force the button to be disabled initially on fresh selection
		state.btnUpsert.Disable()
		state.taskSelected = true

		// Re-bind the change tracking callbacks now that the data is loaded
		state.selectorOwner.OnChanged = onFieldChanged
		state.selectorStatus.OnChanged = onFieldChanged

		// Note: Fyne automatically calls Refresh() inside SetSelected() and Show(),
		// so you don't need to refresh the selectors manually.
		// Just refresh the parent container if it needs a layout recalculation:
		detailPanel.Refresh()
	}

	// =========================================================
	// Layout & Wrappers
	// =========================================================
	splitView := container.NewVSplit(taskList, detailPanel)
	splitView.Offset = 0.55

	logoutBtn := widget.NewButton(
		"Log Out",
		func() {
			fyne.CurrentApp().Preferences().RemoveValue("session_user")
			a.showLoginScreen(window)
		},
	)

	topHeader := container.NewVBox(greetingLabel)
	mainLayout := container.NewBorder(topHeader, logoutBtn, nil, nil, splitView)

	window.SetContent(mainLayout)
}
