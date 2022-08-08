package seeFriendsScreen

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"github.com/el-mendez/redes-proyecto1/views"
)

type seeFriendsScreen struct {
	list        list.Model
	selected    bool
	information string
}

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

func (s *seeFriendsScreen) Init() tea.Cmd { return nil }

func New() *seeFriendsScreen {
	return &seeFriendsScreen{}
}

func (s *seeFriendsScreen) Focus() {
	var items []list.Item
	views.State.FriendsMutex.Lock()
	defer views.State.FriendsMutex.Unlock()

	for friend, devices := range views.State.Friends {
		if len(devices) == 0 {
			items = append(items, item{title: friend, description: "disconnected"})
		}
		for deviceName := range devices {
			items = append(items, item{friend, deviceName})
		}
	}

	s.list = list.New(items, list.NewDefaultDelegate(), 60, 10)
	s.list.DisableQuitKeybindings()
	s.list.SetShowTitle(false)
	s.list.SetShowStatusBar(false)
	s.list.SetShowHelp(true)
}

func (s *seeFriendsScreen) Blur() {
	s.selected = false
	s.list.ResetFilter()
	s.list.ResetSelected()
}

func (s *seeFriendsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if !s.selected && utils.IsEnter(msg) && !s.list.SettingFilter() {
		s.selected = true
		views.State.FriendsMutex.Lock()
		username := s.list.SelectedItem().(item).title
		deviceName := s.list.SelectedItem().(item).description

		defer views.State.FriendsMutex.Unlock()

		if friend, ok := views.State.Friends[username]; ok {
			if state, ok := friend[deviceName]; ok {
				s.information = username
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
			s.information = "There is no user with that username."
		}
		s.information += "\n\n(Press Ctrl+Q to go back)"
	} else {
		s.list, cmd = s.list.Update(msg)
	}

	return s, cmd
}

func (s *seeFriendsScreen) View() string {
	if s.selected {
		return s.information
	}
	return s.list.View() + "\n(Press Ctrl+Q to go back)"
}
