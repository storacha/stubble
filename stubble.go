package stubble

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/storacha/resteep"
)

func Run(stories []Story) error {
	err := resteep.Resteep(
		resteep.RunBubbleTea(
			func(data []byte) (Model, error) {
				m := Model{stories: stories}
				if len(data) > 0 {
					m.currentStoryIndex = int(data[0])
				}
				return m, nil
			},
			func(m Model) ([]byte, error) {
				return []byte{byte(m.currentStoryIndex)}, nil
			},
			tea.WithAltScreen(),
		),
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

var halfBlockBorder = lipgloss.Border{
	Top:         "▄",
	Bottom:      "▀",
	Left:        "▐",
	Right:       "▌",
	TopLeft:     "▗",
	TopRight:    "▖",
	BottomLeft:  "▝",
	BottomRight: "▘",
}

const selectedStoryForeground = lipgloss.Color("205")
const selectedStoryBackground = lipgloss.Color("#333333")

// noEnumeratorList creates a [list.List] without any enumerators or
// indentation. Only the items themselves render.
func noEnumeratorList() *list.List {
	return list.New().
		Indenter(func(items list.Items, index int) string { return "" }).
		Enumerator(func(items list.Items, i int) string { return "" }).
		EnumeratorStyle(lipgloss.NewStyle())
}

func (m Model) View() (result string) {
	l := noEnumeratorList().
		ItemStyleFunc(func(items list.Items, i int) lipgloss.Style {
			style := lipgloss.NewStyle().
				Width(storyListStyle.GetWidth() - 2)

			if i == m.currentStoryIndex {
				// Extend the border to the left side to surround highlight bar.
				border := halfBlockBorder
				border.TopLeft = border.Top
				border.BottomLeft = border.Bottom

				return style.
					Background(selectedStoryBackground).
					Foreground(selectedStoryForeground).
					Bold(true).
					Border(border, true, false, true, true).
					BorderForeground(selectedStoryBackground).
					BorderLeftBackground(selectedStoryForeground)
			}

			style = style.
				Faint(true).
				MarginLeft(1).
				MarginRight(1).
				MarginBottom(1)

			if i == 0 {
				style = style.
					MarginTop(1)
			}

			if i == m.currentStoryIndex-1 {
				style = style.
					MarginBottom(0)
			}

			return style
		})

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
		m.sidebarStyle().Render(
			lipgloss.JoinVertical(lipgloss.Center,
				logoStyle.Render(logo),
				storyListStyle.Render(l.String()),
			)),
		storyView,
	)
}

// https://patorjk.com/software/taag/#p=display&f=Small+Braille&t=Stubble (compressed)
const logo = `⢎⡑⣰⡀⡀⢀⣇⡀⣇⡀⡇⢀⡀
⠢⠜⠘⠤⠣⠼⠧⠜⠧⠜⠣⠣⠭`

var logoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00c3d4ff"))

var storyListStyle = lipgloss.NewStyle().
	Width(21)

func (m Model) fullHeight() lipgloss.Style {
	return lipgloss.NewStyle().Height(m.windowSize.Height)
}

func (m Model) sidebarStyle() lipgloss.Style {
	return m.fullHeight().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(lipgloss.Color("#555555")).
		MarginRight(3)
}
