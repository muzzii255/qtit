package main

import (
	"fmt"
	"strings"
	"os"
	"qtit/dashboard"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))

)




type model struct {
	focusIndex     int
	inputs         []textinput.Model
	cursorMode     cursor.Mode
	width   int
}



func (m model) Init() tea.Cmd {
	return textinput.Blink
}


func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 3),
		width: 300,
	}
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.Width = 30
	
		switch i {
		case 0:
			t.Placeholder = "Qbittorent Host UrL"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Username"
		case 2:
			t.Placeholder = "Password"
		}
	
		m.inputs[i] = t
	}
	return m
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit


		case "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	
	
	}

	cmd := m.updateInputs(msg)
	

	return m, cmd
}


func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)


	return b.String()
}


func main() {
	p := tea.NewProgram(initialModel())
	fmodel, err := p.Run()
	if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
	}
	finalModel, ok := fmodel.(model)
	if !ok {
			fmt.Println("Unexpected error: invalid model type.")
			os.Exit(1)
	}

	qbitCreds := dashboard.Qbit{
		Url: finalModel.inputs[0].Value(),
		Username: finalModel.inputs[1].Value(),
		Password: finalModel.inputs[2].Value(),
	}
	
	
	
	dash := dashboard.New(qbitCreds)
	tea.NewProgram(dash).Run()
	
	
	
}