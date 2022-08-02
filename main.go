package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"github.com/el-mendez/redes-proyecto1/views/loggedInMenu"
	"github.com/el-mendez/redes-proyecto1/views/mainMenu"
	"os"
)

var program *tea.Program

func main() {
	utils.InitializeLogger("./log.conf.json")
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	program = p

	if err := p.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Model struct {
	logged       bool
	mainMenu     *mainMenu.MainMenu
	loggedInMenu *loggedInMenu.LoggedInMenu
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func initialModel() *Model {
	return &Model{
		mainMenu:     mainMenu.InitialMainMenu(),
		loggedInMenu: loggedInMenu.InitialLoggedInMenu(),
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}

	if msg, ok := msg.(mainMenu.LoginResult); ok && msg.Err == nil {
		m.logged = true
		m.loggedInMenu.Start(msg.Client, program)
		return m, nil
	}

	if _, ok := msg.(loggedInMenu.LogOutResult); ok {
		m.logged = false
		m.mainMenu.Start()
		return m, nil
	}

	if m.logged {
		var cmd tea.Cmd
		_, cmd = m.loggedInMenu.Update(msg)
		return m, cmd
	} else {
		var cmd tea.Cmd
		_, cmd = m.mainMenu.Update(msg)
		return m, cmd
	}
}

func (m *Model) View() string {
	if m.logged {
		return m.loggedInMenu.View()
	}
	return m.mainMenu.View()
}
