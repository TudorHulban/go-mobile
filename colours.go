package main

import (
	"image/color"

	"fyne.io/fyne/v2/theme"
)

var statusColors map[string]color.Color

func initStatusColors() {
	statusColors = map[string]color.Color{
		"init":        theme.Color(theme.ColorNamePlaceHolder),
		"not started": theme.Color(theme.ColorNamePlaceHolder),
		"assigned":    theme.Color(theme.ColorNamePrimary),
		"in work":     color.RGBA{R: 255, G: 165, B: 0, A: 255}, // Orange
		"work done":   theme.Color(theme.ColorNameSuccess),
		"to bill":     theme.Color(theme.ColorNameSuccess),
		"invoiced":    color.RGBA{R: 138, G: 43, B: 226, A: 255},  // Purple
		"closed":      color.RGBA{R: 128, G: 128, B: 128, A: 255}, // Gray
	}
}

func getStatusColor(status string) color.Color {
	if statusColors == nil {
		initStatusColors()
	}

	if c, exists := statusColors[status]; exists {
		return c
	}

	return theme.Color(theme.ColorNameForeground)
}
