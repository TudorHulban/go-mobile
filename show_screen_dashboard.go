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

func (a *application) saveTaskChanges(state *stateDashboard) error {
	return a.team.UpsertTask(
		&Task{
			Status:  state.selectorStatus.Selected,
			OwnerID: uint8(state.selectorOwner.SelectedIndex()),
		},
	)
}

func (a *application) buildDetailPanel(state *stateDashboard) fyne.CanvasObject {
	state.selectorOwner = widget.NewSelect(
		a.team.GetTeamMembers(),
		func(s string) {
			a.checkIfModified(state)
		},
	)

	state.selectorStatus = widget.NewSelect(
		[]string{"init", "in work", "closed"},
		func(s string) {
			a.checkIfModified(state)
		},
	)

	state.btnUpsert = widget.NewButton(
		"Save Changes",
		func() {
			a.saveTaskChanges(state)
		},
	)

	state.selectorOwner.Hide()
	state.selectorStatus.Hide()
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
	ownerLabel := widget.NewLabel("Owner:")

	// Use team members for the dropdown selector
	selectorOwner := widget.NewSelect(a.team.GetTeamMembers(), nil)
	selectorOwner.Hide()

	statusLabel := widget.NewLabel("Status:")
	selectorStatus := widget.NewSelect(
		[]string{"init", "not started", "assigned", "in work", "work done", "to bill", "invoiced", "closed"},
		nil,
	)
	selectorStatus.Hide()

	btnUpsert := widget.NewButton("Save Changes", nil)
	btnUpsert.Disable()
	btnUpsert.Hide()

	// 1. Initialize state object ONCE
	state := &stateDashboard{
		selectorOwner:  selectorOwner,
		selectorStatus: selectorStatus,
		btnUpsert:      btnUpsert,
		taskSelected:   false,
	}

	onFieldChanged := func(_ string) { // The string parameter is the new value, but we can ignore it since checkIfModified reads from the widgets directly
		a.checkIfModified(state)
	}

	// 3. Assign the same callback to both selectors
	selectorOwner.OnChanged = onFieldChanged
	selectorStatus.OnChanged = onFieldChanged

	ownerInlineRow := container.NewHBox(ownerLabel, selectorOwner)
	statusInlineRow := container.NewHBox(statusLabel, selectorStatus)

	detailPanel := container.NewVBox(
		detailTitle,
		ownerInlineRow,
		statusInlineRow,
		btnUpsert,
	)

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
	btnUpsert.OnTapped = func() {
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
			if name == selectorOwner.Selected {
				newOwnerID = memberID
				break
			}
		}

		// Update the real data objects
		currentTask.OwnerID = newOwnerID
		currentTask.Status = selectorStatus.Selected

		// UI Feedback: immediate disable while processing/unselecting
		btnUpsert.Disable()

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

		// 2. We don't turn state.taskSelected back to true here because
		// unselecting the list item means no task is active anymore!
		// Instead, we leave it false and let OnSelected flip it to true when they click a new one.
	}

	// =========================================================
	// Row Selection
	// =========================================================
	taskList.OnSelected = func(id widget.ListItemID) {
		taskIDs := a.team.GetTasksIDs()
		if int(id) >= len(taskIDs) {
			return
		}

		// 1. Set this to false FIRST.
		// This acts as a guard. When SetSelected fires OnChanged below,
		// checkIfModified will see this is false and safely return early.
		state.taskSelected = false

		// 2. Assign the ID to your shared state struct
		state.currentTaskID = taskIDs[id]

		// Pull directly from your structured map framework
		task, errGetTask := a.team.GetTaskBy(state.currentTaskID)
		if errGetTask != nil {
			return
		}

		detailTitle.SetText("Task - " + task.Name)

		// Get owner text string from their uint8 ID
		ownerName, _ := a.team.GetTeamMemberBy(task.OwnerID)

		// 3. Update the selector components programmatically
		selectorOwner.SetSelected(ownerName)
		selectorStatus.SetSelected(task.Status)

		// 4. Show the components
		selectorOwner.Show()
		selectorStatus.Show()
		btnUpsert.Show()

		// 5. Now that the UI fully matches the database, turn the guard off
		state.taskSelected = true

		// 6. Force an explicit check to ensure the button is correctly disabled
		a.checkIfModified(state)

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
