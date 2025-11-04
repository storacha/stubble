package stubble

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/storacha/resteep"
)

func Run(stories []Story) error {
	_, err := resteep.Resteep(
		func(data []byte) resteep.ResteepableModel {
			if data == nil {
				return Model{
					stories: stories,
				}
			} else {
				m := Model{
					stories:           stories,
					currentStoryIndex: int(data[0]),
				}
				return m
			}
		},
		tea.WithAltScreen(),
	)
	return err
}

type Model struct {
	stories           []Story
	currentStoryIndex int
	currentStoryModel tea.Model
	windowSize        tea.WindowSizeMsg
}

type Story struct {
	Title    string
	NewModel func() tea.Model
}

func (m Model) Init() tea.Cmd {
	// We haven't "switched to" the current story to initialize it yet, so do that
	// now.
	return m.switchStory(m.currentStoryIndex)
}

type switchStoryMsg int

func (m Model) switchStory(index int) tea.Cmd {
	return func() tea.Msg {
		return switchStoryMsg(index)
	}
}

type storyMsg struct {
	msg tea.Msg
}

func storyCmd(cmd tea.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}

	return func() tea.Msg {
		return storyMsg{msg: cmd()}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			return m, m.switchStory(m.currentStoryIndex + 1)
		case "k", "up":
			return m, m.switchStory(m.currentStoryIndex - 1)
		}

	case tea.WindowSizeMsg:
		m.windowSize = msg

	case switchStoryMsg:
		index := int(msg)
		if index < 0 || index >= len(m.stories) {
			return m, nil
		}
		m.currentStoryIndex = index
		m.currentStoryModel = m.stories[index].NewModel()
		return m, storyCmd(m.currentStoryModel.Init())

	case storyMsg:
		if m.currentStoryModel == nil {
			return m, nil
		}

		newModel, cmd := m.currentStoryModel.Update(msg.msg)
		m.currentStoryModel = newModel
		return m, storyCmd(cmd)
	}

	return m, nil
}

func (m Model) View() (result string) {
	l := list.New().
		Enumerator(m.storyEnumerator).
		EnumeratorStyle(storyEnumeratorStyle).
		ItemStyle(lipgloss.NewStyle().Width(storyListStyle.GetWidth() - 2))

	for _, s := range m.stories {
		l.Item(s.Title)
	}

	var storyView string
	if m.currentStoryModel == nil {
		storyView = ""
	} else {
		storyView = m.currentStoryModel.View()
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.storyListStyle().Render(l.String()),
		storyView,
	)
}

func (m Model) Marshal() ([]byte, error) {
	return []byte{byte(m.currentStoryIndex)}, nil
}

var storyEnumeratorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#888888")).
	MarginRight(1)

var storyListStyle = lipgloss.NewStyle().
	Width(20).
	Border(lipgloss.NormalBorder(), false, true, false, false).
	MarginRight(3)

func (m Model) storyListStyle() lipgloss.Style {
	return storyListStyle.Height(m.windowSize.Height)
}

func (m Model) storyEnumerator(items list.Items, i int) string {
	if i == m.currentStoryIndex {
		return "➤"
	}
	return " "
}
