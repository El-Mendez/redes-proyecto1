package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"github.com/el-mendez/redes-proyecto1/views"
	"github.com/el-mendez/redes-proyecto1/views/loggedInMenu"
	"github.com/el-mendez/redes-proyecto1/views/mainMenu"
	"os"
)

func main() {
	utils.InitializeLogger("./log.conf.json")

	m := initialModel()
	m.mainMenu.Focus()

	p := tea.NewProgram(m, tea.WithAltScreen())
	views.State.P = p

	if err := p.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Model struct {
	mainMenu     *mainMenu.MainMenu
	loggedInMenu *loggedInMenu.LoggedInMenu
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func initialModel() *Model {
	return &Model{
		mainMenu:     mainMenu.New(),
		loggedInMenu: loggedInMenu.New(),
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}

	// The user logs in
	if msg, ok := msg.(views.LoggedInMsg); ok {
		views.State.Client = msg.Client
		go views.HandleIncoming(msg.Client)
		m.mainMenu.Blur()
		m.loggedInMenu.Focus()
		return m, nil
	}

	// The user logs out
	if _, ok := msg.(views.LoggedOutMsg); ok {
		views.State.Client = nil
		m.loggedInMenu.Blur()
		m.mainMenu.Focus()
		return m, nil
	}

	if views.State.Client == nil {
		_, cmd = m.mainMenu.Update(msg)
	} else {
		_, cmd = m.loggedInMenu.Update(msg)
	}

	return m, cmd
}

func (m *Model) View() string {
	if views.State.Client == nil {
		return m.mainMenu.View()
	} else {
		return m.loggedInMenu.View()
	}
}
