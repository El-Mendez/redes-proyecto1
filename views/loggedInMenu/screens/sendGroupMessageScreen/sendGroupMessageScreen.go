package sendGroupMessageScreen

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol"
	utils "github.com/el-mendez/redes-proyecto1/util"
)

type sendMessageScreen struct {
	usernameInput    textinput.Model
	contentInput     textarea.Model
	awaitingUsername bool
}

func (s *sendMessageScreen) Init() tea.Cmd { return nil }

func New() *sendMessageScreen {
	usernameInput := textinput.New()
	usernameInput.Placeholder = "testing@conference.alumchat.fun"

	contentInput := textarea.New()
	contentInput.Placeholder = "Enter your message here..."
	contentInput.ShowLineNumbers = false

	return &sendMessageScreen{usernameInput, contentInput, false}
}

func (s *sendMessageScreen) Focus() {
	s.awaitingUsername = true
	s.usernameInput.Focus()
}

func (s *sendMessageScreen) Blur() {
	s.usernameInput.Reset()
	s.usernameInput.Blur()

	s.contentInput.Reset()
	s.contentInput.Blur()

	s.awaitingUsername = false
}

func (s *sendMessageScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if s.awaitingUsername {
		// it the user tries to submit the username
		if utils.IsEnter(msg) {
			if _, ok := protocol.JIDFromString(s.usernameInput.Value()); ok {
				s.awaitingUsername = false
				s.usernameInput.Blur()
				s.contentInput.Focus()
			}
			return s, nil
		}

		// handle writing on the username field
		s.usernameInput, cmd = s.usernameInput.Update(msg)
		return s, cmd
	}

	if utils.IsEnter(msg) {
		if s.contentInput.Value() != "" {
			username := s.usernameInput.Value()
			content := s.contentInput.Value()

			s.contentInput.Reset()
			return s, sendGroupMessage(username, content)
		}

		return s, nil
	}

	s.contentInput, cmd = s.contentInput.Update(msg)
	return s, cmd
}

func (s *sendMessageScreen) View() string {
	if s.awaitingUsername {
		return fmt.Sprintf("Enter the name of the group you want to sent a message to:\n%s\n\n(press Ctrl+Q to go back)", s.usernameInput.View())
	}

	return s.contentInput.View() + "\n\n(press Ctrl+Q to go back)"
}
