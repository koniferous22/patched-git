package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/koniferous22/patched-git/git-wrapper/utils"
)

type YesNoModel struct {
	Question          string
	InvalidInput      bool
	Entered           bool
	Result            bool
	ShouldExit        bool
	ShouldDisplayHelp bool
}

const YesNoHelpText = "Press\n" +
	"* 'y' for Yes\n" +
	"* 'n' for No\n" +
	"* 'q'/'esc'/'ctrl+c' to Quit\n" +
	"* 'h'/'t' to toggle visiblity of help\n"

func CreateYesNoModel(question string) YesNoModel {
	return YesNoModel{
		Question:          question,
		ShouldDisplayHelp: true,
	}
}

func (m YesNoModel) Init() tea.Cmd {
	return nil
}

func (m YesNoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			m.InvalidInput = false
			m.Result = true
			m.Entered = true
			return m, tea.Quit
		case "n", "N":
			m.InvalidInput = false
			m.Result = false
			m.Entered = true
			return m, tea.Quit
		case "h", "t":
			m.ShouldDisplayHelp = !m.ShouldDisplayHelp
			return m, nil
		case "q", "esc", "ctrl+c":
			m.ShouldExit = true
			return m, tea.Quit
		default:
			m.InvalidInput = true
			return m, nil
		}
	}
	return m, nil
}
func (m YesNoModel) View() string {
	s := fmt.Sprintf("%s%s%s\n", utils.FontBold, m.Question, utils.Reset)
	if m.ShouldDisplayHelp {
		s += YesNoHelpText
	}

	if m.InvalidInput {
		s += fmt.Sprintf("%sInvalid Input%s\n", utils.ColorRed, utils.Reset)
	}
	if m.Entered {
		if m.Result {
			s += fmt.Sprintf("%sEntered Yes%s\n", utils.ColorGreen, utils.Reset)
		} else {
			s += fmt.Sprintf("%sEntered No%s\n", utils.ColorRed, utils.Reset)
		}
	}
	return s
}

func runYesNo(question string) (bool, error) {
	gitignorePromptModel := CreateYesNoModel(question)
	program := tea.NewProgram(gitignorePromptModel)
	result, err := program.Run()
	if err != nil {
		return false, fmt.Errorf("error during yesno prompt:\n%w", err)
	}
	if result, ok := result.(YesNoModel); ok {
		return result.Result, nil
	}
	return false, fmt.Errorf("error retrieving prompt results:\n%w", err)
}
