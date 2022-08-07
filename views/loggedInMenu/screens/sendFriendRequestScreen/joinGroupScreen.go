package sendFriendRequestScreen

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol"
	utils "github.com/el-mendez/redes-proyecto1/util"
)

type sendFriendRequestScreen struct {
	usernameInput textinput.Model
}

func (s *sendFriendRequestScreen) Init() tea.Cmd { return nil }

func New() *sendFriendRequestScreen {
	usernameInput := textinput.New()
	usernameInput.Placeholder = "testing@alumchat.fun"
	usernameInput.Prompt = ""

	return &sendFriendRequestScreen{usernameInput}
}

func (s *sendFriendRequestScreen) Focus() {
	s.usernameInput.Focus()
}

func (s *sendFriendRequestScreen) Blur() {
	s.usernameInput.Reset()
	s.usernameInput.Blur()
}

func (s *sendFriendRequestScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// it the user tries to submit the username
	if utils.IsEnter(msg) {
		username := s.usernameInput.Value()
		if _, ok := protocol.JIDFromString(username); ok {
			s.usernameInput.Blur()
			s.usernameInput.Reset()
			cmd = sendFriendRequest(username)
		}
		return s, cmd
	}

	// handle writing on the group name input
	s.usernameInput, cmd = s.usernameInput.Update(msg)
	return s, cmd
}

func (s *sendFriendRequestScreen) View() string {
	return fmt.Sprintf("Enter the name of your new friend: %s \n\n(press Ctrl+Q to go back)", s.usernameInput.View())
}
