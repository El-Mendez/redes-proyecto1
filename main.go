package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"github.com/el-mendez/redes-proyecto1/views"
	"os"
)

func main() {
	utils.InitializeLogger()
	err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Model struct {
	logged   bool
	mainMenu *views.MainMenu
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func initialModel() *Model {
	return &Model{
		mainMenu: views.InitialMainMenu(),
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}

	if !m.logged {
		var cmd tea.Cmd
		_, cmd = m.mainMenu.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) View() string {
	if !m.logged {
		return m.mainMenu.View()
	}
	return ""
}
