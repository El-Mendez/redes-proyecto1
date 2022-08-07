package gotFriendRequest

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	utils "github.com/el-mendez/redes-proyecto1/util"
)

var options = []string{"Accept", "Reject"}

type FriendRequestScreen struct {
	accepted      int
	username      string
	selectedStyle lipgloss.Style
}

func (s *FriendRequestScreen) Init() tea.Cmd {
	return nil
}

func New(username string) *FriendRequestScreen {
	return &FriendRequestScreen{
		username: username,
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true),
	}
}

func (s *FriendRequestScreen) Focus() {
	s.accepted = 0
}

func (s *FriendRequestScreen) Blur() {
	s.accepted = 0
}

func (s *FriendRequestScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyEnter:
			if s.accepted == 0 {
				return s, acceptFriendRequest(s.username)
			} else {
				return s, rejectFriendshipRequest(s.username)
			}
		case tea.KeyUp, tea.KeyDown:
			s.accepted = utils.EuclideanModule(s.accepted+1, 2)
		}
	}
	return s, nil
}

func (s *FriendRequestScreen) View() string {
	return utils.ViewMenu("You got a new friend request from "+s.username, s.accepted, &options, &s.selectedStyle, nil)
}
