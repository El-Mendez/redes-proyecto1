package loggedInMenu

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"strings"
)

type loggedInAction int

const (
	menu loggedInAction = iota
	messaging
	loggingOut
	addingContact
)

var loggedInOptions = [9]string{"Show all contacts", "Add a contact", "See a user details", "Send a message",
	"Send a message (group)", "Set a presence", "Send a file", "Log Out", "Delete Account"}

func handleIncoming(client *protocol.Client, p *tea.Program, m *LoggedInMenu) {
	client.Send <- &stanzas.Presence{}
	for stanza := range client.Receive {
		switch stanza := stanza.(type) {
		case *stanzas.Message:
			if stanza.From != "" && stanza.Body != "" {
				p.Send(notification{
					text: m.senderStyle.Render(stanza.From) +
						m.typeStyle.Render(" to you") +
						": " + stanza.Body,
				})
			}
		case *stanzas.Presence:
			switch stanza.Type {
			case "subscribed":
				p.Send(notification{m.alertStyle.Render(stanza.From + " accepted your friend request!")})
			case "unsubscribed":
				p.Send(notification{m.alertStyle.Render(stanza.From + " has stopped being friend.")})
			case "subscribe":
				p.Send(friendRequest{stanza.From})
			}
		}
	}
}

type LoggedInMenu struct {
	p      *tea.Program
	client *protocol.Client

	state loggedInAction

	selected      int
	selectedStyle lipgloss.Style

	alertStyle  lipgloss.Style
	typeStyle   lipgloss.Style
	senderStyle lipgloss.Style

	username      string
	usernameInput textinput.Model

	incomingFriend  string
	acceptingFriend bool

	contentInput textarea.Model

	messages []string
	viewport viewport.Model
}

func (m *LoggedInMenu) Start(client *protocol.Client, p *tea.Program) {
	m.state = menu
	m.client = client
	m.p = p
	m.incomingFriend = ""
	m.messages = make([]string, 0)

	go handleIncoming(client, p, m)
}

func InitialLoggedInMenu() *LoggedInMenu {
	ta := textarea.New()
	ta.Placeholder = "Ingresa aquí tu mensaje"
	ta.Reset()
	ta.ShowLineNumbers = false

	ti := textinput.New()
	ti.Placeholder = "testing@alumchat.fun"

	return &LoggedInMenu{
		selectedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true),
		alertStyle:    lipgloss.NewStyle().Faint(true),
		senderStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("86")).
			Foreground(lipgloss.Color("0")).
			Bold(true),
		typeStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("86")).
			Foreground(lipgloss.Color("9")).
			Bold(true),
		viewport:      viewport.New(30, 10),
		usernameInput: ti,
		contentInput:  ta,
	}
}

func (m *LoggedInMenu) Init() tea.Cmd {
	return nil
}

func (m *LoggedInMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case friendRequest:
		m.incomingFriend = msg.from
		return m, nil
	case notification:
		m.messages = append(m.messages, msg.text)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.incomingFriend != "" {
				var cmd tea.Cmd
				if m.acceptingFriend {
					cmd = m.confirmFriendship(m.incomingFriend)
				} else {
					cmd = m.denyFriendship(m.incomingFriend)
				}
				m.incomingFriend = ""
				return m, cmd
			}

			switch m.state {
			case menu:
				switch m.selected {
				case 1:
					m.state = addingContact
					m.usernameInput.Focus()
					return m, nil
				case 3:
					m.state = messaging
					m.usernameInput.Focus()
					return m, nil
				case 7:
					m.state = loggingOut
					client := m.client
					m.client = nil
					return m, m.logOut(client)
				case 8:
					m.state = loggingOut
					client := m.client
					m.client = nil
					return m, m.deleteAccount(client)
				}
			case messaging:
				if m.username == "" {
					name := m.usernameInput.Value()
					if _, ok := protocol.JIDFromString(name); ok {
						m.username = name
						m.usernameInput.Blur()
						m.usernameInput.Reset()

						m.contentInput.Focus()
						return m, nil
					}
				} else {
					content := m.contentInput.Value()
					if content != "" {
						m.contentInput.Reset()
						return m, m.sendMessage(m.client, m.username, content)
					}
				}
			case addingContact:
				name := m.usernameInput.Value()
				if _, ok := protocol.JIDFromString(name); ok {
					m.username = name
					m.usernameInput.Blur()
					m.usernameInput.Reset()
					m.state = menu
					return m, m.addContact(name)
				}
			}
		case tea.KeyDown:
			if m.incomingFriend != "" {
				m.acceptingFriend = !m.acceptingFriend
			} else if m.state == menu {
				m.selected = utils.EuclideanModule(m.selected+1, len(loggedInOptions))
				return m, nil
			}
		case tea.KeyUp:
			if m.incomingFriend != "" {
				m.acceptingFriend = !m.acceptingFriend
			} else if m.state == menu {
				m.selected = utils.EuclideanModule(m.selected-1, len(loggedInOptions))
				return m, nil
			}
		case tea.KeyEscape, tea.KeyCtrlQ:
			m.state = menu
			m.username = ""
			m.usernameInput.Blur()
			m.usernameInput.Reset()
			m.contentInput.Blur()
			m.contentInput.Reset()
			return m, nil
		}
	}

	// Update viewport
	var vpCmd, ciCmd, uiCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.contentInput, ciCmd = m.contentInput.Update(msg)
	m.usernameInput, uiCmd = m.usernameInput.Update(msg)

	return m, tea.Batch(vpCmd, ciCmd, uiCmd)
}

func (m *LoggedInMenu) View() string {
	if m.incomingFriend != "" {
		if m.acceptingFriend {
			return "You have a friend request from " + m.senderStyle.Render(m.incomingFriend) +
				"\nDo you want to accept it? \n\n" + m.selectedStyle.Render("[X] Yes") + "\n[ ] No \n\n"
		}
		return "You have a friend request from " + m.senderStyle.Render(m.incomingFriend) +
			"\nDo you want to accept it? \n\n[ ] Yes \n" + m.selectedStyle.Render("[X] No\n\n")
	}

	switch m.state {
	case messaging:
		if m.username == "" {
			return fmt.Sprintf("%s\nIngresa el usuario del destinatario \n%v \n\n (presiona Ctrl+Q para volver)",
				m.viewport.View(), m.usernameInput.View())
		} else {
			return fmt.Sprintf("%s\n\n%v \n\n (presiona Ctrl+Q para volver)",
				m.viewport.View(), m.contentInput.View())
		}
	case addingContact:
		return fmt.Sprintf("%s\nIngresa el usuario del usuario que quieres agregar\n%v \n\n (presiona Ctrl+Q para volver)",
			m.viewport.View(), m.usernameInput.View())
	case menu:
		return fmt.Sprintf("%s \n\n%s \n%s \n%s \n%s \n%s \n%s \n%s \n%s \n%s", m.viewport.View(),
			// Sé que podría hacer un loop, pero me rehúso xd
			utils.MenuOption(loggedInOptions[0], 0 == m.selected, m.selectedStyle),
			utils.MenuOption(loggedInOptions[1], 1 == m.selected, m.selectedStyle),
			utils.MenuOption(loggedInOptions[2], 2 == m.selected, m.selectedStyle),
			utils.MenuOption(loggedInOptions[3], 3 == m.selected, m.selectedStyle),
			utils.MenuOption(loggedInOptions[4], 4 == m.selected, m.selectedStyle),
			utils.MenuOption(loggedInOptions[5], 5 == m.selected, m.selectedStyle),
			utils.MenuOption(loggedInOptions[6], 6 == m.selected, m.selectedStyle),
			utils.MenuOption(loggedInOptions[7], 7 == m.selected, m.selectedStyle),
			utils.MenuOption(loggedInOptions[8], 8 == m.selected, m.selectedStyle),
		)
	}
	return "Loading..."
}
