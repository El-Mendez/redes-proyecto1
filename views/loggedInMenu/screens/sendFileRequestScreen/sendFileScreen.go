package sendFileRequestScreen

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol"
	utils "github.com/el-mendez/redes-proyecto1/util"
)

type sendFriendRequestScreen struct {
	usernameInput textinput.Model
	fileInput     textinput.Model
}

func (s *sendFriendRequestScreen) Init() tea.Cmd { return nil }

func New() *sendFriendRequestScreen {
	usernameInput := textinput.New()
	usernameInput.Placeholder = "testing@alumchat.fun/device"
	usernameInput.Prompt = ""

	fileInput := textinput.New()
	fileInput.Placeholder = "test.txt"
	fileInput.Prompt = ""

	return &sendFriendRequestScreen{usernameInput, fileInput}
}

func (s *sendFriendRequestScreen) Focus() {
	s.usernameInput.Focus()
}

func (s *sendFriendRequestScreen) Blur() {
	s.usernameInput.Reset()
	s.usernameInput.Blur()
	s.fileInput.Reset()
	s.fileInput.Blur()
}

func (s *sendFriendRequestScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// it the user tries to submit the username
	if s.usernameInput.Focused() {
		if utils.IsEnter(msg) {
			username := s.usernameInput.Value()
			if jid, ok := protocol.JIDFromString(username); ok && jid.DeviceName != "" {
				s.usernameInput.Blur()
				cmd = s.fileInput.Focus()
			}
			return s, cmd
		}
		s.usernameInput, cmd = s.usernameInput.Update(msg)
		return s, cmd
	}

	if utils.IsEnter(msg) {
		s.fileInput.Blur()
		return s, sendFileRequest(s.usernameInput.Value(), s.fileInput.Value())
	}
	s.fileInput, cmd = s.fileInput.Update(msg)
	return s, cmd
}

func (s *sendFriendRequestScreen) View() string {
	if s.usernameInput.Focused() {
		return fmt.Sprintf("Who do you want to send the file to? %s \n\n(press Ctrl+Q to go back)", s.usernameInput.View())
	} else {
		return fmt.Sprintf("Which file do you want to send? %s \n\n(press Ctrl+Q to go back)", s.fileInput.View())
	}
}
