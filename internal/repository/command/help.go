package command

import (
	"fmt"
	"strings"
)

type HelpRepositoryInput struct {
	SlashCommand string
}

type IHelpRepository interface {
	HelpMessage(input HelpRepositoryInput) (string, error)
}

type HelpRepository struct {
}

func NewHelpRepository() *HelpRepository {
	return &HelpRepository{}
}

func (h *HelpRepository) HelpMessage(input HelpRepositoryInput) (string, error) {
	return strings.Join([]string{
		fmt.Sprintf("usage: %s <command> <subcommand> [parameters]", input.SlashCommand),
		"",
		"These are common commands:",
		"",
		"help: show this message",
		"",
		"To see more help text, you can run:",
		"",
		fmt.Sprintf("%s <command> help", input.SlashCommand),
		fmt.Sprintf("%s <command> <subcommand> help", input.SlashCommand),
	}, "\n"), nil
}
