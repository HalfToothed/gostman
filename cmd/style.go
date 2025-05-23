package cmd

import "github.com/charmbracelet/lipgloss"

const maxWidth = 80

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	HeaderDecoration,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

var (
	borderStyle = lipgloss.NewStyle().
			Padding(0, 0).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("62"))

	focusedBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("205"))
)

var titleStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("62")).
	Foreground(lipgloss.Color("230"))

var headingStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("11")).
	Foreground(lipgloss.Color("0")).
	MarginLeft(2)

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(0, 0, 0, 0)
	s.HeaderText = lg.NewStyle().
		Foreground(green).
		Bold(true).
		Padding(0, 1, 0, 0)
	s.HeaderDecoration = lg.NewStyle().
		Foreground(indigo).
		Padding(0, 1, 0, 0)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.Copy().
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))
)

var (
	inactiveTabStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, false).BorderForeground(green).Padding(0, 2).Margin(0, 1)
	activeTabStyle   = inactiveTabStyle.Border(lipgloss.DoubleBorder(), false, false, true, false)
)

func (m help) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, m.styles.HeaderText.Render("+-- "+text), lipgloss.WithWhitespaceChars("/"), lipgloss.WithWhitespaceForeground(indigo))
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, m.styles.HeaderText.Render("+-- "+text))
}

func (m Model) appBoundaryMessage(text string) string {
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, m.styles.ErrorHeaderText.Render(text))
}

func (m board) appBoundaryMessage(text string) string {
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, m.styles.ErrorHeaderText.Render(text))
}

func (m env) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, m.styles.HeaderText.Render("+-- "+text))
}

func (m env) appBoundaryMessage(text string) string {
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, m.styles.ErrorHeaderText.Render(text))
}
