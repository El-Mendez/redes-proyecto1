package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"github.com/el-mendez/redes-proyecto1/views"
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
	mainMenu *mainMenu.MainMenu
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func initialModel() *Model {
	return &Model{
		mainMenu: mainMenu.New(),
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(views.LoggedInMsg); ok {
		msg.Client.Close()
		return m, tea.Quit
	}

	var cmd tea.Cmd
	_, cmd = m.mainMenu.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return m.mainMenu.View()
}
