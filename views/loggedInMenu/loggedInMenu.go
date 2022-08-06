package loggedInMenu

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"github.com/el-mendez/redes-proyecto1/views"
)

var options = []string{"Send a message", "Send a group message", "Send a friend request",
	"Change your status", "Send a file", "Join a group chat", "Mostrar todos los contactos",
	"Detalles de un contacto", "Delete Account", "Logout"}
var screens = [8]views.Screen{}

type LoggedInMenu struct {
	inMenu        bool
	selected      int
	selectedStyle lipgloss.Style
	vp            viewport.Model
}

func (m *LoggedInMenu) Init() tea.Cmd {
	return nil
}

func New() *LoggedInMenu {
	return &LoggedInMenu{
		inMenu: true,
		vp:     viewport.New(30, 10),
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true),
	}
}

func (m *LoggedInMenu) Focus() {
	m.inMenu = true
	m.selected = 0
}

func (m *LoggedInMenu) Blur() {
	m.inMenu = true
	m.selected = 0
}

func (m *LoggedInMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Force return to menu
	if utils.IsCtrlQ(msg) && !m.inMenu {
		screens[m.selected].Blur()
		m.inMenu = true
		return m, nil
	}

	if m.inMenu {
		// Enter an option menu
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.Type {
			case tea.KeyEnter:
				switch m.selected {
				case 8:
					return m, deleteAccount
				case 9:
					return m, logout
				}

				m.inMenu = false
				screens[m.selected].Focus()

			case tea.KeyUp:
				m.selected = utils.EuclideanModule(m.selected-1, len(options))
			case tea.KeyDown:
				m.selected = utils.EuclideanModule(m.selected+1, len(options))
			}
		}
		return m, nil
	}

	return screens[m.selected].Update(msg)
}

func (m *LoggedInMenu) View() string {
	if m.inMenu {
		return utils.ViewMenu(m.vp.View(), m.selected, &options, &m.selectedStyle, nil)
	}

	return m.vp.View() + "\n\n" + screens[m.selected].View()
}
