package views

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/el-mendez/redes-proyecto1/protocol"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"strings"
	"time"
)

var mainOptions = []string{"Login", "Create Account", "Quit"}

type MainMenu struct {
	selectedStyle lipgloss.Style
	errorStyle    lipgloss.Style
	textArea      textinput.Model
	spin          spinner.Model

	selected int
	logging  bool
	signing  bool
	loading  bool

	username *protocol.JID
	password string
	err      string
}

func (m *MainMenu) Init() tea.Cmd {
	return textinput.Blink
}

func InitialMainMenu() *MainMenu {
	spin := spinner.New()
	spin.Spinner = spinner.Dot

	return &MainMenu{
		selected:      0,
		selectedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true),
		errorStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true),
		textArea:      textinput.New(),
		spin:          spin,
	}
}

func (m *MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Getting login response
	if result, ok := msg.(LoginResult); ok {
		m.loading = false
		m.username = nil
		m.password = ""

		if result.Err != nil && result.Client != nil {
			result.Client.Close()
		}

		if result.Err != nil && m.logging {
			m.err = m.errorStyle.Render(fmt.Sprintf("Could not Log In: %v", result.Err))
		} else if result.Err != nil && m.signing {
			m.err = m.errorStyle.Render(fmt.Sprintf("Could not Sign Up: %v", result.Err))
		} else {
			panic("TODO")
		}

		m.signing = false
		m.logging = false
	}

	if m.loading {
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		return m, cmd
	}

	if msg, ok := msg.(tea.KeyMsg); ok && !m.loading {
		// Clear all errors when you get any input
		m.err = ""

		// select anything in the main menu
		if !m.logging && !m.signing {
			switch msg.Type {
			case tea.KeyUp:
				m.selected = utils.EuclideanModule(m.selected-1, len(mainOptions))
			case tea.KeyDown:
				m.selected = utils.EuclideanModule(m.selected+1, len(mainOptions))
			case tea.KeyEnter:
				switch m.selected {
				case 0:
					m.logging = true
					m.textArea.Reset()
					m.textArea.Placeholder = "testing@alumchat.fun"
					m.textArea.EchoMode = textinput.EchoNormal
					m.textArea.Focus()
				case 1:
					m.signing = true
					m.textArea.Reset()
					m.textArea.Placeholder = "testing@alumchat.fun"
					m.textArea.EchoMode = textinput.EchoNormal
					m.textArea.Focus()
				case 2:
					return m, tea.Quit
				}
			}
			return m, nil
		}

		// On escape go back to main menu and reset everything
		if msg.Type == tea.KeyEsc {
			m.logging = false
			m.signing = false
			m.loading = false
			m.username = nil
			m.password = ""
		}

		if msg.Type == tea.KeyEnter {
			if m.username == nil {
				// Entering the username
				username, ok := protocol.JIDFromString(m.textArea.Value())
				if !ok {
					m.err = m.errorStyle.Render("Invalid account. Please make sure you enter an account in the form user@domain or user@domain/resource.")
				} else {
					// Prepare for entering the password
					m.username = &username
					m.textArea.Reset()
					m.textArea.Placeholder = ""
					m.textArea.EchoMode = textinput.EchoPassword
					m.textArea.Focus()
				}
				return m, nil
			} else {
				// Entering the password
				if strings.TrimSpace(m.textArea.Value()) == "" {
					return m, nil
				}
				m.password = m.textArea.Value()
				m.loading = true

				if m.signing {
					// Sing in
					return m, tea.Batch(m.spin.Tick, m.signup(*m.username, m.password))
				} else {
					// Log in
					return m, tea.Batch(m.spin.Tick, m.login(*m.username, m.password))
				}
			}
		}

	}
	if m.logging || m.signing {
		var cmd tea.Cmd
		m.textArea, cmd = m.textArea.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *MainMenu) View() string {
	if m.logging || m.signing {
		if m.username == nil {
			return fmt.Sprintf("Ingresa una cuenta: \n%s \n\n%s \n\n", m.textArea.View(), m.err)
		} else {
			if m.loading {
				return fmt.Sprintf("Ingresa una cuenta: %s \n\n%s Loading, please wait... \n\n", m.username.BaseJid(), m.spin.View())
			} else {
				return fmt.Sprintf("Ingresa una cuenta: %s \nIngresa tu contraseña: \n%s \n\n%s \n\n", m.username.BaseJid(), m.textArea.View(), m.err)
			}
		}
	}

	return fmt.Sprintf("Bienvenido al Chat de Méndez! \n\n%s \n%s \n%s \n\n%s \n\n",
		utils.MenuOption(mainOptions[0], 0 == m.selected, m.selectedStyle),
		utils.MenuOption(mainOptions[1], 1 == m.selected, m.selectedStyle),
		utils.MenuOption(mainOptions[2], 2 == m.selected, m.selectedStyle),
		m.err,
	)
}

type LoginResult struct {
	Client *protocol.Client
	Err    error
}

func (m *MainMenu) signup(jid protocol.JID, password string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		c, err := protocol.SignUp(&jid, password)
		return LoginResult{c, err}
	}
}

func (m *MainMenu) login(jid protocol.JID, password string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		c, err := protocol.LogIn(&jid, password)
		return LoginResult{c, err}
	}
}
