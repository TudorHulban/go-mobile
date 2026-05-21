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

	updateBtn := widget.NewButton("Save Changes", nil)
	updateBtn.Disable()
	updateBtn.Hide()

	var currentTaskID uint8
	var taskSelected bool

	checkIfModified := func() {
		if !taskSelected {
			return
		}

		currentTask, errGetTask := a.team.GetTaskBy(currentTaskID)
		if errGetTask != nil {
			return
		}

		assignedTeamMember, errGetPerson := a.team.GetTeamMemberBy(currentTask.OwnerID)
		if errGetPerson != nil {
			assignedTeamMember = "" // Handle case where task has no owner yet or owner missing safely

		}

		hasOwnerChanged := selectorOwner.Selected != assignedTeamMember
		hasStatusChanged := selectorStatus.Selected != currentTask.Status

		if hasOwnerChanged || hasStatusChanged {
			updateBtn.Enable()
		} else {
			updateBtn.Disable()
		}
	}

	selectorOwner.OnChanged = func(newOwner string) {
		checkIfModified()
	}

	selectorStatus.OnChanged = func(newStatus string) {
		checkIfModified()
	}

	ownerInlineRow := container.NewHBox(ownerLabel, selectorOwner)
	statusInlineRow := container.NewHBox(statusLabel, selectorStatus)

	detailPanel := container.NewVBox(
		detailTitle,
		ownerInlineRow,
		statusInlineRow,
		updateBtn,
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
	updateBtn.OnTapped = func() {
		if !taskSelected {
			return
		}

		taskSelected = false

		currentTask, errGetTask := a.team.GetTaskBy(currentTaskID)
		if errGetTask != nil {
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

		updateBtn.Disable()

		// Refresh using Fyne's local row selection index
		taskIDs := a.team.GetTasksIDs()

		for rowIdx, dbID := range taskIDs {
			if dbID == currentTaskID {
				taskList.RefreshItem(rowIdx)
				taskList.Unselect(rowIdx)
				break
			}
		}

		taskSelected = true
	}

	// =========================================================
	// Row Selection
	// =========================================================
	taskList.OnSelected = func(id widget.ListItemID) {
		taskIDs := a.team.GetTasksIDs()
		if int(id) >= len(taskIDs) {
			return
		}

		// Map the UI index to your actual Database ID
		currentTaskID = taskIDs[id]
		taskSelected = false

		// Pull directly from your structured map framework
		task, errGetTask := a.team.GetTaskBy(currentTaskID)
		if errGetTask != nil {
			return
		}

		// fmt.Printf("owner raw ID = %d\n", task.OwnerID)
		// fmt.Printf("selected status=[%s]\n", task.Status)

		detailTitle.SetText("Task - " + task.Name)

		// Get owner text string from their uint8 ID to update selector component safely
		ownerName, _ := a.team.GetTeamMemberBy(task.OwnerID)

		selectorOwner.SetSelected(ownerName)
		selectorStatus.SetSelected(task.Status)

		selectorOwner.Show()
		selectorStatus.Show()
		updateBtn.Show()
		updateBtn.Disable()

		selectorOwner.Refresh()
		selectorStatus.Refresh()
		detailPanel.Refresh()

		taskSelected = true
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
