package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
)

const (
	host = "localhost"
	port = 23234
)

type viewState int

const (
	fileListView viewState = iota
	fileContentView
)

type model struct {
	cursor int
	fileNames []string
	currentView viewState
	fileContent string
	scrollPosition int
	terminalHeight int
}

func readFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, len(files))
	for i := range files {
		fileNames[i] = files[i].Name()
	}

	return fileNames, nil
}

func renderEntry(str string, selected bool) string {
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	containerStyle := lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("63"))
	if (selected) {
		containerStyle = lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("226"))
	}

	textContent := textStyle.Render(str)
	containerContent := containerStyle.Render(textContent)

	return containerContent
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
		fileNames: fileNames,
		terminalHeight: pty.Window.Height,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.currentView == fileListView {
				if m.cursor > 0 {
					m.cursor--
				}
			} else {
				if (m.scrollPosition > 0) {
					m.scrollPosition--
				}
			}
		case "down":
			if m.currentView == fileListView {
				if m.cursor < len(m.fileNames) {
					m.cursor++
				}
			} else {
				maxScroll := len(strings.Split(m.fileContent, "\n")) - m.terminalHeight 
        if m.scrollPosition < maxScroll {
            m.scrollPosition++
        }
			}
		case "enter":
			if m.currentView == fileListView {
				selectedFile := m.fileNames[m.cursor - 1]
				content, err := os.ReadFile("data/" + selectedFile)
				if err != nil {
					m.fileContent = "Error reading file"
				} else {
					m.fileContent = string(content)
				}
				m.currentView = fileContentView
			} else {
				m.currentView = fileListView
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if (m.currentView == fileListView) {
		s, err := glamour.Render("# Files\n", "dark")
		for i, fileName := range m.fileNames {
			selected := m.cursor == i + 1
			styledFileName := renderEntry(fileName, selected)
			s += styledFileName + "\n"
		}
		s += "\n"
		s += "Press 'q' to quit\n"
	
		if (err != nil) {
			return "Error: Unable to parse markdown"
		}
		return fmt.Sprint(s)
	} else {
		parsedFileContent, err := glamour.Render(m.fileContent, "dark")
		if err != nil {
			return "Error: Unable to parse markdown"
		}

		lines := strings.Split(parsedFileContent, "\n")
		start := m.scrollPosition
		end := start + m.terminalHeight

		if end > len(lines) {
			end = len(lines)
		}

		displayLines := lines[start:end]
		displayContent := strings.Join(displayLines, "\n")

		return fmt.Sprint(displayContent)
	}
}
