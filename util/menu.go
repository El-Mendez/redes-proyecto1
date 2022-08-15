package utils

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// ViewMenu is a reusable menu view using the tea library. Does not handle logic, it only displays it.
func ViewMenu(title string, counter int, options *[]string, selectedStyle *lipgloss.Style, footer *string) string {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "%s \n\n", title)

	for i, option := range *options {
		if i == counter {
			fmt.Fprintf(&builder, selectedStyle.Render("[✓] %s"), option)
			fmt.Fprintf(&builder, "\n")
		} else {
			fmt.Fprintf(&builder, "[ ] %s \n", option)
		}
	}
	if footer != nil {
		fmt.Fprintf(&builder, "\n\n%s \n\n", *footer)
	} else {
		fmt.Fprintf(&builder, "\n\n")
	}
	return builder.String()
}
