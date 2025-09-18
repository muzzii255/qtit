package dashboard

import (
	"fmt"
	"net/http"
	"net/url"

	// "net/url"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg struct{}

type Model struct {
	httpClient  *http.Client
	torrents    []Torrent
	selected    int
	paginator   paginator.Model
	prevPage    int
	host        string
	commandMode bool
	command     textinput.Model
	statusMessage string
	statusStyle lipgloss.Style
	width int
}



func New(qb Qbit) Model {
	ti := textinput.New()
	ti.Placeholder = ":command"
	ti.Prompt = ":"
	ti.Width = 30
	ti.PromptStyle = CommandPromptStyle
	ti.TextStyle = CommandTextStyle

	client, err := LoginToQbit(qb.Url, qb.Username, qb.Password)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	torrents, err := FetchTorrents(client, qb.Url)

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 8
	p.ActiveDot = PaginatorActiveDot.Render("•")
	p.InactiveDot = PaginatorInactiveDot.Render("·")
	p.SetTotalPages(len(torrents))

	return Model{
		httpClient:  client,
		torrents:    torrents,
		selected:    0,
		paginator:   p,
		prevPage:    p.Page,
		host:        qb.Url,
		command:     ti,
		commandMode: false,
		statusMessage: "no activity yet",
		width: 500,
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	currentPage := m.paginator.Page

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.commandMode {
			var cmd tea.Cmd
			m.command, cmd = m.command.Update(msg)

			if msg.Type == tea.KeyEnter {
				input := strings.TrimSpace(m.command.Value())
				hash := m.pageItems()[m.selected].Hash
				name := m.pageItems()[m.selected].Name

				if strings.HasPrefix(input, "add ") {
					magnet := strings.TrimSpace(strings.TrimPrefix(input, "add "))
					if strings.HasPrefix(magnet, "magnet:?xt=") {
						err := AddMagnet(m.httpClient, m.host, magnet)
						if err != nil {
							m.statusMessage = fmt.Sprintf("Error adding magnet: %s",err.Error())
							m.statusStyle = StatusError
						} else {
							m.statusMessage = "magnet added!"
							m.statusStyle = StatusSuccess
							m.torrents, _ = FetchTorrents(m.httpClient, m.host)
						}
					}
				}

				switch input {
				case "stop":
					err := PostTorrentAction(m.httpClient, m.host, "stop", url.Values{"hashes": {hash}})
					if err != nil{
						m.statusMessage = fmt.Sprintf("%s error stopping!, error: %s", name,err.Error())
						m.statusStyle = StatusError
					}else{
						m.statusMessage = fmt.Sprintf("%s torrent stopped!", name)
						m.statusStyle = StatusSuccess
						
					}
				case "start":
					err := PostTorrentAction(m.httpClient, m.host, "start", url.Values{"hashes": {hash}})
					if err != nil{
						m.statusMessage = fmt.Sprintf("%s error starting!, error: %s", name,err.Error())
						m.statusStyle = StatusError
					}else{
						m.statusMessage = fmt.Sprintf("%s torrent started!", name)
						m.statusStyle = StatusSuccess
						
					}
				case "delete":
					err := PostTorrentAction(m.httpClient, m.host, "delete", url.Values{"hashes": {hash}, "deleteFiles": {"false"}})
					
					if err != nil{
						m.statusMessage = fmt.Sprintf("%s error deleting!, error: %s", name,err.Error())
						m.statusStyle = StatusError
					}else{
						m.statusMessage = fmt.Sprintf("%s torrent deleted!", name)
						m.statusStyle = StatusSuccess
						
					}
				case "recheck":
					err := PostTorrentAction(m.httpClient, m.host, "recheck", url.Values{"hashes": {hash}})
					if err != nil{
						m.statusMessage = fmt.Sprintf("%s error rechecking!, error: %s", name,err.Error())
						m.statusStyle = StatusError
					}else{
						m.statusMessage = fmt.Sprintf("%s torrent rechecked!", name)
						m.statusStyle = StatusSuccess
						
					}
				}
				updatedTorrents, err := FetchTorrents(m.httpClient, m.host)
				if err == nil {
					m.torrents = updatedTorrents
				}
				
				m.commandMode = false
				m.command.SetValue("")
			}

			
			if msg.Type == tea.KeyEsc {
				m.commandMode = false
				m.command.SetValue("")
			}

			return m, cmd
		}

		switch msg.String() {
		case "up":
			if m.selected > 0 {
				m.selected--
			}
		case "down":
			if m.selected < len(m.pageItems())-1 {
				m.selected++
			}
		case "left":
			m.paginator, cmd = m.paginator.Update(msg)
		case "right":
			m.paginator, cmd = m.paginator.Update(msg)
		case "q":
			Logout(m.httpClient, m.host)
			return m, tea.Quit
		case ":":
			m.commandMode = true
			m.command.Focus()
			return m, nil

		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	
		
	case tickMsg:
		torrents, err := FetchTorrents(m.httpClient, m.host)
		if err == nil {
			m.torrents = torrents
		}
		return m, tick()

	}

	_, isKeyMsg := msg.(tea.KeyMsg)
	if msg != nil && !isKeyMsg {
		m.paginator, cmd = m.paginator.Update(msg)
	}

	if currentPage != m.paginator.Page {
		m.selected = 0
		m.prevPage = m.paginator.Page
	}

	return m, cmd
}

func (m Model) View() string {
	var b strings.Builder

	centered := lipgloss.Place(10, 1, lipgloss.Center, lipgloss.Center,
		HeaderStyle.Render("Qtit v1.0 — Connected to: "+m.host),
	)
	b.WriteString(centered + "\n")
	
	headers := fmt.Sprintf("  %-50s  %-9s  %-10s  %-10s  %-8s  %-6s  %-10s  %-10s  %-16s  %-6s  %-6s  %-6s  %-6s  %-16s",
		"NAME", "%", "↓SPEED", "↑SPEED", "PEERS", "ETA", "STATE", "SIZE", "ADDED", "SEEDS", "LEECH", "PRIV", "F-START", "S-SEED")
	b.WriteString(headers + "\n")

	items := m.pageItems()
	for i, t := range items {
		progressText := fmt.Sprintf("%.2f%%", t.Progress*100)
		speedInfo := formatSpeed(t.Speed)
		peersInfo := fmt.Sprintf(" %-3d", t.Peers)
		etaInfo := formatETA(t.ETA)
		name := TruncateName(t.Name, 50)
		size := FormatSize(t.Size)
		added := FormatAddedOn(t.AddedOn)
		upseedInfo := formatSpeed(t.UpSpeed)

		line := fmt.Sprintf("%-50s  %-9s  %-10s  %-10s  %-8s  %-6s  %-10s  %-10s  %-16s  %-6d  %-6d  %-6t  %-7t  %-16t",
			name, progressText, speedInfo, upseedInfo, peersInfo, etaInfo, t.State, size, added, t.Seeds,
			t.Leech, t.Private, t.ForceStart, t.SuperSeeding)

		if i == m.selected {
			line = "› " + SelectedStyle.Render(line)
		} else {
			line = "  " + line
		}

		b.WriteString(line + "\n")
	}

	b.WriteString("\n")

	pagination := "  " + m.paginator.View()
	b.WriteString(pagination + "\n\n")
	if m.commandMode {
		box := PanelStyle

		prompt := CommandPromptStyle.Render("COMMAND MODE")

		cmdInput := m.command.View()

		// Combine and style
		cmdBox := box.Render(prompt + "\n" + cmdInput)
		b.WriteString(cmdBox + "\n")
	}
	controls := SubtleStyle.Render("↑↓ Navigate  ←→ Pages  : Commands  Q Quit")
	b.WriteString(controls)
	
	if m.statusMessage != "" {
		b.WriteString("\n" + m.statusStyle.Render(m.statusMessage) + "\n")
	}
	
	return b.String()
}

func (m Model) pageItems() []Torrent {
	start, end := m.paginator.GetSliceBounds(len(m.torrents))

	if start >= len(m.torrents) {
		return []Torrent{}
	}
	if end > len(m.torrents) {
		end = len(m.torrents)
	}

	return m.torrents[start:end]
}
