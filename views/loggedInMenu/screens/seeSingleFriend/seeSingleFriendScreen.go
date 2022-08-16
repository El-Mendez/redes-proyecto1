package seeSingleFriend

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"github.com/el-mendez/redes-proyecto1/views"
)

type seeSingleFriend struct {
	usernameInput textinput.Model
	list          list.Model
	selected      bool
	information   string
}

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

func (s *seeSingleFriend) Init() tea.Cmd { return nil }

func New() *seeSingleFriend {
	usernameInput := textinput.New()
	usernameInput.Prompt = ""
	usernameInput.Placeholder = "test@alumchat.fun"
	return &seeSingleFriend{usernameInput: usernameInput}
}

func (s *seeSingleFriend) Focus() {
	s.usernameInput.Focus()
}

func (s *seeSingleFriend) Blur() {
	s.usernameInput.Reset()
	s.usernameInput.Blur()
	s.selected = false
	s.list.ResetFilter()
	s.list.ResetSelected()
}

func (s *seeSingleFriend) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if s.usernameInput.Focused() {
		if utils.IsEnter(msg) {
			s.usernameInput.Blur()

			var items []list.Item

			views.State.FriendsMutex.Lock()
			if friend, ok := views.State.Friends[s.usernameInput.Value()]; ok {
				for name, devices := range friend {
					items = append(items, item{name, devices.Show})
				}
				if len(items) == 0 {
					items = append(items, item{"Disconnected", ""})
				}
			}

			views.State.FriendsMutex.Unlock()

			s.list = list.New(items, list.NewDefaultDelegate(), 60, 10)
			s.list.DisableQuitKeybindings()
			s.list.SetShowTitle(false)
			s.list.SetShowStatusBar(false)
			s.list.SetShowHelp(true)
		}
		s.usernameInput, cmd = s.usernameInput.Update(msg)
		return s, cmd
	}

	if !s.selected && utils.IsEnter(msg) && !s.list.SettingFilter() {
		s.selected = true
		deviceName := s.list.SelectedItem().(item).title

		views.State.FriendsMutex.Lock()

		defer views.State.FriendsMutex.Unlock()

		if friend, ok := views.State.Friends[s.usernameInput.Value()]; ok {
			if state, ok := friend[deviceName]; ok {
				s.information = s.usernameInput.Value() + "/" + deviceName
				s.information += "\nShow: " + state.Show

				if len(state.Status) > 0 {
					s.information += "\nStatus: "
					for _, str := range state.Status {
						s.information += "\n - " + str
					}
				}

			} else {
				s.information = "The user is currently not logged in."
			}
		} else {
			s.information = "You have no friends with that username."
		}
		s.information += "\n\n(Press Ctrl+Q to go back)"
	} else {
		s.list, cmd = s.list.Update(msg)
	}

	return s, cmd
}

func (s *seeSingleFriend) View() string {
	if s.usernameInput.Focused() {
		return "Whose status do you want to see? " + s.usernameInput.View()
	} else if s.selected {
		return s.information
	}
	return s.list.View() + "\n(Press Ctrl+Q to go back)"
}
