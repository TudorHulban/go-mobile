package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func showDashboardScreen(window fyne.Window, username string) {
	// Header Area
	greetingLabel := widget.NewLabel("Hello, " + username + " 👋")
	greetingLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Define Detail Widgets (Using unified regular labels for font consistency)
	detailTitle := widget.NewLabel("Select a task to view details")

	ownerLabel := widget.NewLabel("Owner: ")
	ownerSelector := widget.NewSelect(
		[]string{
			"Alice", "Bob", "Charlie", "Unassigned", // Add your typical process owners here
		},
		nil,
	)
	ownerSelector.Hide()

	statusLabel := widget.NewLabel("Status:")
	statusSelector := widget.NewSelect(
		[]string{
			"init", "not started", "assigned", "in work", "work done", "to bill", "invoiced", "closed",
		},
		nil,
	)
	statusSelector.Hide()

	// Create the Update Button (Disabled by default)
	updateBtn := widget.NewButton("Save Changes", nil)
	updateBtn.Disable()
	updateBtn.Hide()

	// Track transactional row mutations within this block scope
	var currentTaskID widget.ListItemID
	var taskSelected = false

	// Helper function to check if anything actually changed
	checkIfModified := func() {
		if !taskSelected {
			return
		}
		hasOwnerChanged := ownerSelector.Selected != _Tasks[currentTaskID].Owner
		hasStatusChanged := statusSelector.Selected != _Tasks[currentTaskID].Status

		if hasOwnerChanged || hasStatusChanged {
			updateBtn.Enable()
		} else {
			updateBtn.Disable()
		}
	}

	// Attach modification checks to BOTH dropdown changes
	ownerSelector.OnChanged = func(newOwner string) {
		checkIfModified()
	}
	statusSelector.OnChanged = func(newStatus string) {
		checkIfModified()
	}

	// Group widgets cleanly into a structured card layout
	// Combines the text label and combo selector
	// onto a single horizontal line
	ownerInlineRow := container.NewHBox(ownerLabel, ownerSelector)
	statusInlineRow := container.NewHBox(statusLabel, statusSelector)

	detailPanel := container.NewVBox(
		detailTitle,
		ownerInlineRow,
		statusInlineRow,
		updateBtn,
	)

	// Build Top List Widget
	var taskList *widget.List

	taskList = widget.NewList(
		func() int {
			return len(_Tasks)
		},

		func() fyne.CanvasObject {
			taskName := widget.NewLabel("Template Task Name")
			statusBadge := canvas.NewText("STATUS", color.White)
			statusBadge.TextStyle = fyne.TextStyle{Bold: true}

			return container.NewBorder(nil, nil, nil, statusBadge, taskName)
		},

		func(id widget.ListItemID, item fyne.CanvasObject) {
			task := _Tasks[id]

			box := item.(*fyne.Container)
			nameLabel := box.Objects[0].(*widget.Label)
			badgeText := box.Objects[1].(*canvas.Text)

			nameLabel.SetText(task.Name)

			badgeText.Text = "  " + task.Status + "  "
			badgeText.Color = getStatusColor(task.Status)

			box.Refresh()
		},
	)

	// Define the save operation routine
	updateBtn.OnTapped = func() {
		if !taskSelected {
			return
		}

		// 1. Drop the guard flag BEFORE mutating data or refreshing components
		taskSelected = false

		// commit chosen state
		_Tasks[currentTaskID].Owner = ownerSelector.Selected
		_Tasks[currentTaskID].Status = statusSelector.Selected

		// Gray out the button immediately after successful update
		updateBtn.Disable()

		// Refresh entire window layout to repaint the list row colors instantly
		// window.Canvas().Refresh(window.Content())
		taskList.RefreshItem(currentTaskID)

		// 2. Clear any active row selection highlight so it doesn't leave ghost states floating
		taskList.Unselect(currentTaskID)

		// 3. Re-engage your modification listeners safely now that updates are done
		taskSelected = true
	}

	// Intercept user selection to populate information fields down below
	taskList.OnSelected = func(id widget.ListItemID) {
		currentTaskID = id

		// Set guard flag to false briefly
		// to prevent SetSelected from falsely triggering OnChanged
		taskSelected = false
		task := _Tasks[id]

		// Format output context string explicitly
		detailTitle.SetText("Task - " + task.Name)

		// Synchronize dropdown state with currently stored asset values
		ownerSelector.SetSelected(task.Owner)
		statusSelector.SetSelected(task.Status)

		// Reveal bottom controls on first selection
		ownerSelector.Show()
		statusSelector.Show()
		updateBtn.Show()

		// Lock button back down until a new option is explicitly chosen
		updateBtn.Disable()

		// targeted refresh
		ownerSelector.Refresh()
		statusSelector.Refresh()
		detailPanel.Refresh()

		taskSelected = true
	}

	// Assemble Layout Architecture
	splitView := container.NewVSplit(taskList, detailPanel)
	splitView.Offset = 0.55 // Split layout cleanly near the screen center

	logoutBtn := widget.NewButton(
		"Log Out",
		func() {
			fyne.CurrentApp().
				Preferences().
				RemoveValue("session_user")

			showLoginScreen(window)
		},
	)

	topHeader := container.NewVBox(greetingLabel)
	mainLayout := container.NewBorder(topHeader, logoutBtn, nil, nil, splitView)

	window.SetContent(mainLayout)
}
