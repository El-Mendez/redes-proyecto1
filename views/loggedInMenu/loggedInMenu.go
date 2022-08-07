package loggedInMenu

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"github.com/el-mendez/redes-proyecto1/views"
	"github.com/el-mendez/redes-proyecto1/views/loggedInMenu/screens/gotFriendRequest"
	"github.com/el-mendez/redes-proyecto1/views/loggedInMenu/screens/joinGroupScreen"
	"github.com/el-mendez/redes-proyecto1/views/loggedInMenu/screens/seeFriendsScreen"
	"github.com/el-mendez/redes-proyecto1/views/loggedInMenu/screens/sendFriendRequestScreen"
	"github.com/el-mendez/redes-proyecto1/views/loggedInMenu/screens/sendGroupMessageScreen"
	"github.com/el-mendez/redes-proyecto1/views/loggedInMenu/screens/sendMessageScreen"
	"github.com/el-mendez/redes-proyecto1/views/loggedInMenu/screens/setStatusScreen"
	"strings"
)

var options = []string{"Send a message", "Send a group message", "Send a friend request", "Join a group chat",
	"Change your status", "Send a file", "Show all your contacts", "Show a contact details",
	"Delete Account", "Logout"}

var screens = [8]views.Screen{
	sendMessageScreen.New(),
	sendGroupMessageScreen.New(),
	sendFriendRequestScreen.New(),
	joinGroupScreen.New(),
	setStatusScreen.New(),
	setStatusScreen.New(), // TODO replace with sending a file
	seeFriendsScreen.New(),
}

type LoggedInMenu struct {
	selected      int
	selectedStyle lipgloss.Style
	vp            viewport.Model
	content       *strings.Builder
	currentScreen views.Screen
}

func (m *LoggedInMenu) Init() tea.Cmd {
	return nil
}

func New() *LoggedInMenu {
	return &LoggedInMenu{
		vp: viewport.New(30, 10),
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true),
		content:       &strings.Builder{},
		currentScreen: nil,
	}
}

func (m *LoggedInMenu) Focus() {
	m.currentScreen = nil
	m.selected = 0
}

func (m *LoggedInMenu) Blur() {
	if m.currentScreen != nil {
		m.currentScreen.Blur()
		m.currentScreen = nil
	}
	m.selected = 0
	m.content.Reset()
	m.vp.SetContent("")
}

func (m *LoggedInMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Force return to menu
	if utils.IsCtrlQ(msg) && m.currentScreen != nil {
		m.currentScreen.Blur()
		m.currentScreen = nil
		return m, nil
	}

	// Enter an option menu
	if msg, ok := msg.(views.Notification); ok {
		m.content.WriteString("\n")
		m.content.WriteString(msg.Msg)
		m.vp.SetContent(m.content.String())
		m.vp.GotoBottom()
	}
	if msg, ok := msg.(views.NotificationAndBack); ok {
		m.content.WriteString("\n")
		m.content.WriteString(msg.Msg)
		m.vp.SetContent(m.content.String())
		m.vp.GotoBottom()

		if m.currentScreen != nil {
			m.currentScreen.Blur()
			m.currentScreen = nil
		}
	}

	if msg, ok := msg.(views.FriendRequest); ok {
		if m.currentScreen != nil {
			m.currentScreen.Blur()
		}

		m.currentScreen = gotFriendRequest.New(msg.From)
	}

	if m.currentScreen == nil {
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

				m.currentScreen = screens[m.selected]
				m.currentScreen.Focus()

			case tea.KeyUp:
				m.selected = utils.EuclideanModule(m.selected-1, len(options))
			case tea.KeyDown:
				m.selected = utils.EuclideanModule(m.selected+1, len(options))
			}
		}
		return m, nil
	}

	return m.currentScreen.Update(msg)
}

func (m *LoggedInMenu) View() string {
	if m.currentScreen == nil {
		return utils.ViewMenu(m.vp.View(), m.selected, &options, &m.selectedStyle, nil)
	}

	return m.vp.View() + "\n\n" + m.currentScreen.View()
}
