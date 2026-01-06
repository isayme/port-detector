package view

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/isayme/port-detector/src/connection"
	"github.com/isayme/port-detector/src/langs"
)

type model struct {
	table table.Model

	keys keyMap
	help help.Model

	pause bool
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	// 启动第一个定时器
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		ports, _ := connection.ListeningPorts()
		m.table.SetRows(convertToRows(ports))
		return m, tickCmd()
	case tea.WindowSizeMsg:
		m.table.SetColumns(getColumns(msg.Width))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			// toggle help
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			//	quit
			return m, tea.Quit
		case key.Matches(msg, m.keys.Pause):
			//	toggle refresh data
			m.pause = !m.pause
		case key.Matches(msg, m.keys.Refresh):
			// manual refresh data
			ports, _ := connection.ListeningPorts()
			m.table.SetRows(convertToRows(ports))
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	helpView := m.help.View(m.keys)

	return lipgloss.NewStyle().Padding(1, 2).Render(m.table.View()) +
		"\n" + fmt.Sprintf("auto refresh: %s", formatSwitch(!m.pause)) +
		"\n" + helpView
	//return lipgloss.NewStyle().Padding(1, 2).Render(m.table.View()) +
	//	"\nPress q to quit, pause: " + strconv.FormatBool(m.pause)
}

func newModel(rows []table.Row) model {
	t := table.New(
		table.WithColumns(getColumns(100)),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		BorderTop(true).
		//BorderLeft(true).
		//BorderRight(true).
		Bold(true)

	style.Selected = style.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57"))

	t.SetStyles(style)

	m := model{
		table: t,
		help:  help.New(),
		keys:  keys,
	}

	return m
}

func getColumns(width int) []table.Column {
	cmdLineWidth := width

	columns := []table.Column{
		//{Title: "FAMILY", Width: 6},
		{Title: langs.Localize("PROTO"), Width: 5},
		{Title: langs.Localize("PORT"), Width: 6},
		{Title: langs.Localize("LOCAL_ADDR"), Width: 15},
		{Title: langs.Localize("PID"), Width: 6},
		{Title: langs.Localize("USER"), Width: 10},
		{Title: langs.Localize("CPU"), Width: 10},
		{Title: langs.Localize("MEM"), Width: 6},
		{Title: langs.Localize("PROCESS"), Width: cmdLineWidth},
	}

	//for _, item := range columns {
	//	cmdLineWidth = cmdLineWidth - item.Width
	//}
	//columns = append(columns, table.Column{Title: langs.Localize("CMDLINE"), Width: cmdLineWidth})

	return columns
}

func convertToRows(ports []connection.ListenPortInfo) []table.Row {
	rows := make([]table.Row, 0, len(ports))
	for _, p := range ports {
		rows = append(rows, table.Row{
			formatProto(p.Proto),
			formatUint32(p.Port),
			p.LocalAddr,
			formatInt32(p.PID),
			p.Username,
			formatCPU(p.CPUPercent),
			formatMEM(p.MemInfo),
			p.ProcName,
		})
	}

	return rows
}

func Render() {
	ports, err := connection.ListeningPorts()
	if err != nil {
		fmt.Println("list listening ports failed, error:", err)
		os.Exit(1)
	}

	m := newModel(convertToRows(ports))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("failed to start:", err)
		os.Exit(1)
	}
}
