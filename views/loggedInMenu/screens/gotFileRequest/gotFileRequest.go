package gotFileRequest

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	utils "github.com/el-mendez/redes-proyecto1/util"
)

var options = []string{"Accept", "Reject"}

type FileRequestScreen struct {
	accepted      int
	username      string
	sid           string
	id            string
	selectedStyle lipgloss.Style
	locationInput textinput.Model
}

func (s *FileRequestScreen) Init() tea.Cmd {
	return nil
}

func New(username string, sid string, id string) *FileRequestScreen {
	input := textinput.New()
	input.Placeholder = "example.txt"

	return &FileRequestScreen{
		username:      username,
		sid:           sid,
		id:            id,
		locationInput: input,
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true),
	}
}

func (s *FileRequestScreen) Focus() {
	s.accepted = 0
	s.locationInput.Blur()
}

func (s *FileRequestScreen) Blur() {
	s.accepted = 0
	s.locationInput.Blur()
}

func (s *FileRequestScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if s.locationInput.Focused() {
		if utils.IsEnter(msg) && s.locationInput.Value() != "" {
			return s, acceptFileRequest(s.username, s.id, s.sid, s.locationInput.Value())
		} else {
			s.locationInput, cmd = s.locationInput.Update(msg)
			return s, cmd
		}

	} else if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyEnter:
			if s.accepted == 0 {
				return s, s.locationInput.Focus()
			} else {
				return s, rejectFileRequest(s.username, s.id)
			}
		case tea.KeyUp, tea.KeyDown:
			s.accepted = utils.EuclideanModule(s.accepted+1, 2)
		}
	}
	return s, nil
}

func (s *FileRequestScreen) View() string {
	if s.locationInput.Focused() {
		return fmt.Sprintf("Where do you want to save the file? \n%s", s.locationInput.View())
	} else {
		return utils.ViewMenu("You got a new file from "+s.username+". Do you want to accept it?", s.accepted, &options, &s.selectedStyle, nil)
	}
}
