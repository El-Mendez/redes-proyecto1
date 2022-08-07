package setStatusScreen

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	utils "github.com/el-mendez/redes-proyecto1/util"
)

var status = []string{"None", "Available", "Away", "Extended Away", "Do not Disturb"}
var show = [5]string{"", "chat", "away", "xa", "dnd"}

type SetStatusScreen struct {
	s               int
	awaitingMessage bool
	messageInput    textinput.Model
	selectedStyle   lipgloss.Style
}

func (s *SetStatusScreen) Init() tea.Cmd {
	return nil
}

func New() *SetStatusScreen {
	messageInput := textinput.New()
	messageInput.Placeholder = "Listening to music..."
	messageInput.Prompt = ""

	return &SetStatusScreen{
		messageInput: messageInput,
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true),
	}
}

func (s *SetStatusScreen) Focus() {
	s.s = 0
	s.messageInput.Focus()
}

func (s *SetStatusScreen) Blur() {
	s.s = 0
	s.awaitingMessage = false
	s.messageInput.Blur()
	s.messageInput.Reset()
}

func (s *SetStatusScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if s.awaitingMessage {
		if utils.IsEnter(msg) {
			message := s.messageInput.Value()
			s.messageInput.Blur()
			s.messageInput.Reset()
			s.awaitingMessage = false
			return s, setStatus(show[s.s], message)
		}
		s.messageInput, cmd = s.messageInput.Update(msg)
	} else {
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.Type {
			case tea.KeyEnter:
				s.awaitingMessage = true
				s.messageInput.Focus()
			case tea.KeyUp:
				s.s = utils.EuclideanModule(s.s-1, len(status))
			case tea.KeyDown:
				s.s = utils.EuclideanModule(s.s+1, len(status))
			}
		}
	}
	return s, cmd
}

func (s *SetStatusScreen) View() string {
	if s.awaitingMessage {
		return fmt.Sprintf("Your new status will be: %s. \nEnter your status message: %s\n\n(press Ctrl+Q to go back)",
			status[s.s], s.messageInput.View())
	}
	return utils.ViewMenu("What do you want your new status to be? ", s.s, &status, &s.selectedStyle, nil)
}
