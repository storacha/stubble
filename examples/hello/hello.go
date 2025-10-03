package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/storacha/stubble"
)

func main() {
	err := stubble.Run([]stubble.Story{
		{
			Title: "Hello World",
			NewModel: func() tea.Model {
				return model{}
			},
		},
		{
			Title: "Hello World with text",
			NewModel: func() tea.Model {
				return model{
					text: "oh, hi.",
				}
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
