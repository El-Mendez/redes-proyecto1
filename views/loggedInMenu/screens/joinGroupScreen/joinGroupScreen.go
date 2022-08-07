package joinGroupScreen

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol"
	utils "github.com/el-mendez/redes-proyecto1/util"
)

type joinGroupScreen struct {
	groupInput    textinput.Model
	aliasInput    textinput.Model
	awaitingGroup bool
}

func (s *joinGroupScreen) Init() tea.Cmd { return nil }

func New() *joinGroupScreen {
	groupInput := textinput.New()
	groupInput.Placeholder = "testing@conference.alumchat.fun"
	groupInput.Prompt = ""

	aliasInput := textinput.New()
	aliasInput.Placeholder = "alias"
	groupInput.Prompt = ""

	return &joinGroupScreen{groupInput, aliasInput, false}
}

func (s *joinGroupScreen) Focus() {
	s.awaitingGroup = true
	s.groupInput.Focus()
}

func (s *joinGroupScreen) Blur() {
	s.groupInput.Reset()
	s.groupInput.Blur()

	s.aliasInput.Reset()
	s.aliasInput.Blur()

	s.awaitingGroup = false
}

func (s *joinGroupScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if s.awaitingGroup {
		// it the user tries to submit the username
		if utils.IsEnter(msg) {
			if jid, ok := protocol.JIDFromString(s.groupInput.Value()); ok && jid.BaseJid() == s.groupInput.Value() {
				s.awaitingGroup = false
				s.groupInput.Blur()
				s.aliasInput.Focus()
			}
			return s, nil
		}

		// handle writing on the group name input
		s.groupInput, cmd = s.groupInput.Update(msg)
		return s, cmd
	}

	if utils.IsEnter(msg) {
		if s.aliasInput.Value() != "" {
			username := s.groupInput.Value()
			content := s.aliasInput.Value()
			s.aliasInput.Blur()
			s.aliasInput.Reset()

			return s, joinGroup(username, content)
		}

		return s, nil
	}

	s.aliasInput, cmd = s.aliasInput.Update(msg)
	return s, cmd
}

func (s *joinGroupScreen) View() string {
	if s.awaitingGroup {
		return fmt.Sprintf("Enter the name of the group you want to sent a message to: %s\n\n\n(press Ctrl+Q to go back)", s.groupInput.View())
	}

	return fmt.Sprintf("Enter the name of the group you want to sent a message to: %s\nWhat do you want to be called by?: %s\n\n(press Ctrl+Q to go back)", s.groupInput.View(), s.aliasInput.View())
}
