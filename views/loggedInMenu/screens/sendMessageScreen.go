package screens

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/views"
)

type sendMessageScreen struct {
	to            string
	usernameInput textinput.Model
	messageInput  textarea.Model
}

func (s *sendMessageScreen) Init() tea.Cmd {
	return nil
}

func NewSendMessageScreen() *sendMessageScreen {
	input := textinput.New()

	message := textarea.New()
	message.ShowLineNumbers = false

	input.Placeholder = "usuario@dominio"
	return &sendMessageScreen{usernameInput: input, messageInput: message}
}

func (s *sendMessageScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// if still waiting for the user to specify who to send the message to
			if s.to == "" {
				name := s.usernameInput.Value()
				if _, ok := protocol.JIDFromString(name); ok {
					s.to = name
					s.usernameInput.Blur()
					s.messageInput.Focus()
					return s, nil
				}
			} else {
				// if we already know who to send the message to
				content := s.messageInput.Value()
				if content != "" {
					s.messageInput.Reset()
					return s, views.GlobalState.sendMessage(s.to, content)
				}
			}
		}
	}
	return s, nil
}

func (s *sendMessageScreen) Focus() {
	s.usernameInput.Focus()
}

func (s *sendMessageScreen) Blur() {
	s.usernameInput.Reset()
	s.usernameInput.Blur()
	s.messageInput.Reset()
	s.messageInput.Blur()
	s.to = ""
}

func (s *sendMessageScreen) View() string {
	if s.to == "" {
		return fmt.Sprintf("Ingresa el usuario del destinatario \n%v \n\n (presiona Ctrl+Q para volver)", s.usernameInput.View())
	}
	return fmt.Sprintf("%s\n\n%s \n\n (presiona Ctrl+Q para volver)", s.messageInput.View())
}
