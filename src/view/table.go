package view

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/isayme/port-detector/src/connection"
	"github.com/isayme/port-detector/src/langs"
)

type model struct {
	table table.Model
	pause bool
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return MsgRefresh(0)
	})
}

func (m model) Init() tea.Cmd {
	// return nil
	return tea.Batch(
		// fetchTableData(),
		tickCmd(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MsgRefresh:
		ports, _ := connection.ListeningPorts()
		m.table.SetRows(convertToRows(ports))
	case tea.WindowSizeMsg:
		m.table.SetColumns(getColumns(msg.Width))
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "p":
			m.pause = !m.pause
		case "r":
			ports, _ := connection.ListeningPorts()
			m.table.SetRows(convertToRows(ports))
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.NewStyle().Padding(1, 2).Render(m.table.View()) +
		"\nPress q to quit, pause: " + strconv.FormatBool(m.pause)
}

func getColumns(width int) []table.Column {
	cmdLineWidth := width

	columns := []table.Column{
		{Title: "FAMILY", Width: 6},
		{Title: langs.Localize("PROTO"), Width: 5},
		{Title: langs.Localize("PORT"), Width: 6},
		{Title: langs.Localize("LOCAL_ADDR"), Width: 15},
		{Title: langs.Localize("PID"), Width: 6},
		{Title: langs.Localize("PROCESS"), Width: 14},
		{Title: "EXTRA", Width: 50},
	}

	for _, item := range columns {
		cmdLineWidth = cmdLineWidth - item.Width
	}
	columns = append(columns, table.Column{Title: langs.Localize("CMDLINE"), Width: cmdLineWidth})

	return columns
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
		Bold(true)

	style.Selected = style.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57"))

	t.SetStyles(style)

	m := model{table: t}

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			<-ticker.C

			if m.pause {
				// continue
			}
			// ports, _ := connection.ListeningPorts()
			// m.table.SetRows(convertToRows(ports))
			triggerRefreshTableData()
		}
	}()

	return m
}

func triggerRefreshTableData() tea.Cmd {
	return func() tea.Msg {
		// 模拟耗时操作（如 HTTP 请求、DB 查询）

		return MsgRefresh(0)
	}
}

func convertToRows(ports []connection.ListenPortInfo) []table.Row {
	rows := make([]table.Row, 0, len(ports))
	for _, p := range ports {
		rows = append(rows, table.Row{
			p.Family,
			p.Proto,
			strconv.Itoa(int(p.Port)),
			p.LocalAddr,
			strconv.Itoa(int(p.PID)),
			p.ProcName,
			p.Cmdline,
			fmt.Sprintf("user: %s, cpu: %0.2f%%, mem: %s", p.Username, p.CPUPercent,
				humanize.Bytes(p.MemInfo.RSS),
			),
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
