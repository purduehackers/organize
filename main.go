package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
)

type viewState int

const (
	host = "0.0.0.0"
	port = 23234
)

const (
	fileListView viewState = iota
	fileContentView
)

type model struct {
	cursor           int
	ready            bool
	viewport         viewport.Model
	fileNames        []string
	currentView      viewState
	selectedFileName string
	fileContent      string
	terminalHeight   int
	help             help.Model
	keys             keyMap
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Quit, k.Back}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Quit, k.Back},
	}
}

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Error("could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		return nil, nil
	}

	fileNames, err := readFiles("data")
	if err != nil {
		wish.Fatalln(s, "can't read directory")
		return nil, nil
	}

	m := model{
		fileNames:      fileNames,
		terminalHeight: pty.Window.Height,
		help:           help.New(),
		keys:           keys,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 && m.currentView == fileListView {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.fileNames) && m.currentView == fileListView {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Top):
			m.viewport.GotoTop()
		case key.Matches(msg, m.keys.Enter):
			if m.currentView == fileListView && m.cursor > 0 {
				selectedFile := m.fileNames[m.cursor-1]
				content, err := os.ReadFile("data/" + selectedFile)
				if err != nil {
					m.fileContent = "Error reading file"
				} else {
					m.fileContent = string(content)
					m.selectedFileName = selectedFile
				}
				parsedFileContent, err := glamour.Render(m.fileContent, "dark")
				if err != nil {
					m.viewport.SetContent("Error parsing markdown")
				}
				m.viewport.SetContent(parsedFileContent)
				m.currentView = fileContentView
			}
		case key.Matches(msg, m.keys.Back):
			if m.currentView == fileContentView {
				m.currentView = fileListView
				m.viewport.GotoTop()
			}
		}
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = false
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.currentView == fileListView {
		s := joinPurdueHackersView()
		s += introDescriptionView(m.viewport.Width)
		for i, fileName := range m.fileNames {
			selected := m.cursor == i+1
			styledFileName := positionListItemView(fileName, selected)
			s += styledFileName + "\n"
		}
		s += "\n"

		return fmt.Sprint(s)
	} else {
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	}
}
