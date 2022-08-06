package LoginScreen

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol"
	utils "github.com/el-mendez/redes-proyecto1/util"
)

type logInScreen struct {
	usernameInput    textinput.Model
	passwordInput    textinput.Model
	spin             spinner.Model
	loading          bool
	awaitingUsername bool
}

func (s *logInScreen) Init() tea.Cmd { return nil }

func New() *logInScreen {
	usernameInput := textinput.New()
	usernameInput.Prompt = ""

	spin := spinner.New()
	spin.Spinner = spinner.Dot

	passwordInput := textinput.New()
	passwordInput.Prompt = ""
	passwordInput.EchoMode = textinput.EchoPassword

	return &logInScreen{usernameInput, passwordInput, spin, false, false}
}

func (s *logInScreen) Focus() {
	s.awaitingUsername = true
	s.usernameInput.Focus()
}

func (s *logInScreen) Blur() {
	s.usernameInput.Reset()
	s.usernameInput.Blur()

	s.passwordInput.Reset()
	s.passwordInput.Blur()

	s.loading = false
	s.awaitingUsername = false
}

func (s *logInScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if s.loading {
		s.spin, cmd = s.spin.Update(msg)
		return s, cmd
	}

	if s.awaitingUsername {
		// it the user tries to submit the username
		if utils.IsEnter(msg) {
			if _, ok := protocol.JIDFromString(s.usernameInput.Value()); ok {
				s.awaitingUsername = false
				s.usernameInput.Blur()
				s.passwordInput.Focus()
			}
			return s, nil
		}

		// handle writing on the username field
		s.usernameInput, cmd = s.usernameInput.Update(msg)
		return s, cmd
	}

	if utils.IsEnter(msg) {
		if s.passwordInput.Value() != "" {

			username, _ := protocol.JIDFromString(s.usernameInput.Value())
			password := s.passwordInput.Value()

			s.passwordInput.Blur()
			s.loading = true
			return s, tea.Batch(logIn(username, password), s.spin.Tick)
		}

		return s, nil
	}

	s.passwordInput, cmd = s.passwordInput.Update(msg)
	return s, cmd
}

func (s *logInScreen) View() string {
	if s.loading {
		return fmt.Sprintf("Enter your account: %s\nEnter your password: %s\n\n%s Loading...",
			s.usernameInput.View(), s.passwordInput.View(), s.spin.View())
	}

	if s.awaitingUsername {
		return fmt.Sprintf("Enter your account: %s\n\n", s.usernameInput.View())
	}

	return fmt.Sprintf("Enter your account: %s\nEnter your password: %s\n\n",
		s.usernameInput.View(), s.passwordInput.View())
}
