package view

import (
	"fmt"
	"os"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/isayme/port-detector/src/connection"
	"github.com/isayme/port-detector/src/langs"
)

type model struct {
	table table.Model

	keys keyMap
	help help.Model

	pause bool
	nonRootWarning bool
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
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height - 4)
	case tea.KeyPressMsg:
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

func (m model) pageIndicator() string {
	totalRows := len(m.table.Rows())
	pageHeight := m.table.Height()
	if pageHeight < 1 {
		pageHeight = 1
	}
	if totalRows <= pageHeight {
		return ""
	}
	currentPage := m.table.Cursor()/pageHeight + 1
	totalPages := (totalRows + pageHeight - 1) / pageHeight
	return fmt.Sprintf("Page %d/%d (%d rows total)  ", currentPage, totalPages, totalRows)
}

func (m model) View() tea.View {
	helpView := m.help.View(m.keys)

	pageHint := m.pageIndicator()
	if pageHint != "" {
		pageHint = "\n" + pageHint
	}

	warningMsg := ""
	if m.nonRootWarning {
		warningMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Bold(true).
			Render("Warning: Running as non-root user. Only ports you have permission to see will be displayed.")
		warningMsg = "\n" + warningMsg + "\n"
	}

	return tea.NewView(
		lipgloss.NewStyle().Padding(1, 2).Render(m.table.View()) +
			warningMsg +
			"\n" + fmt.Sprintf("auto refresh: %s", formatSwitch(!m.pause)) +
			pageHint +
			"\n" + helpView)
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

func isRoot() bool {
	return os.Getuid() == 0
}

func Render() {
	ports, err := connection.ListeningPorts()
	if err != nil {
		fmt.Println("list listening ports failed, error:", err)
		os.Exit(1)
	}

	m := newModel(convertToRows(ports))
	m.nonRootWarning = !isRoot()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("failed to start:", err)
		os.Exit(1)
	}
}
