package mainMenu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"github.com/el-mendez/redes-proyecto1/views"
	"github.com/el-mendez/redes-proyecto1/views/mainMenu/screens/LoginScreen"
	"github.com/el-mendez/redes-proyecto1/views/mainMenu/screens/SignUpScreen"
)

var mainOptions = []string{"Login", "Create Account", "Quit"}
var screens = [2]views.Screen{
	LoginScreen.New(),
	SignUpScreen.New(),
}

type MainMenu struct {
	inMenu        bool
	selected      int
	selectedStyle lipgloss.Style
	err           string
}

func (m *MainMenu) Init() tea.Cmd {
	return nil
}

func (m *MainMenu) Focus() {
	m.inMenu = true
	m.selected = 0
}

func (m *MainMenu) Blur() {
	if !m.inMenu && m.selected < len(screens) {
		screens[m.selected].Blur()
	}
	m.selected = 0
	m.inMenu = true
}

func New() *MainMenu {
	return &MainMenu{
		selected: 0,
		selectedStyle: lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true),
	}
}

func (m *MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Force return to menu
	if utils.IsCtrlQ(msg) && !m.inMenu {
		screens[m.selected].Blur()
		m.inMenu = true
		return m, nil
	}

	if msg, ok := msg.(views.ErrorMsg); ok {
		m.inMenu = true
		screens[m.selected].Blur()
		m.err = msg.Err
		return m, nil
	}

	if m.inMenu {
		// Enter an option menu
		if msg, ok := msg.(tea.KeyMsg); ok {
			m.err = ""
			switch msg.Type {
			case tea.KeyEnter:
				if m.selected == 2 { // if quitting
					return m, tea.Quit
				}
				m.inMenu = false
				screens[m.selected].Focus()

			case tea.KeyUp:
				m.selected = utils.EuclideanModule(m.selected-1, len(mainOptions))
			case tea.KeyDown:
				m.selected = utils.EuclideanModule(m.selected+1, len(mainOptions))
			}
		}
		return m, nil
	}

	return screens[m.selected].Update(msg)
}

func (m *MainMenu) View() string {
	if m.inMenu {
		if m.err != "" {
			return utils.ViewMenu("Bienvenido al Chat de Méndez!", m.selected, &mainOptions, &m.selectedStyle, &m.err)
		}

		return utils.ViewMenu("Bienvenido al Chat de Méndez!", m.selected, &mainOptions, &m.selectedStyle, nil)
	}

	return screens[m.selected].View()
}
