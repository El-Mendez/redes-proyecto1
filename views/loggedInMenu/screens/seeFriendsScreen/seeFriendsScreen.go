package seeFriendsScreen

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/views"
	"strings"
)

type seeFriendsScreen struct {
}

func (s *seeFriendsScreen) Init() tea.Cmd { return nil }

func New() *seeFriendsScreen {
	return &seeFriendsScreen{}
}

func (s *seeFriendsScreen) Focus() {}

func (s *seeFriendsScreen) Blur() {}

func (s *seeFriendsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s *seeFriendsScreen) View() string {
	var builder strings.Builder
	views.State.FriendsMutex.Lock()
	defer views.State.FriendsMutex.Unlock()

	for friend, devices := range views.State.Friends {
		builder.WriteString(friend)
		builder.WriteString(" with devices: \n")

		for deviceName, state := range devices {
			builder.WriteString("\n\t")
			builder.WriteString(deviceName)
			builder.WriteString(" status: ")
			builder.WriteString(state.Show)
			builder.WriteString(" ")
			builder.WriteString(fmt.Sprintf("%v", state.Status))
		}
	}

	builder.WriteString("\n\n(Press Ctrl+Q to go back)")
	return builder.String()
}
