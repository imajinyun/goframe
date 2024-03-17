package command

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/github"
	"github.com/spf13/cast"

	"github.com/imajinyun/goframe/cobra"
	"github.com/imajinyun/goframe/util"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noneStyle           = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = blurredStyle.Copy()
	focusedButton       = focusedStyle.Copy().Render("[Submit]")
	blurredButton       = fmt.Sprintf("[%s]", blurredStyle.Render("Submit"))
)

type model struct {
	err    error
	idx    int
	mode   cursor.Mode
	inputs []textinput.Model
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlR:
			m.mode++
			if m.mode > cursor.CursorHide {
				m.mode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.mode)
			}
			return m, tea.Batch(cmds...)
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown, tea.KeyEnter:
			s := msg.String()
			if s == tea.KeyEnter.String() && m.idx == len(m.inputs) {
				m.Make()
				return m, nil
			}

			if s == tea.KeyUp.String() || s == tea.KeyShiftTab.String() {
				m.idx--
			} else {
				m.idx++
			}

			if m.idx > len(m.inputs) {
				m.idx = 0
			} else if m.idx < 0 {
				m.idx = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.idx {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noneStyle
				m.inputs[i].TextStyle = noneStyle
			}
			return m, tea.Batch(cmds...)
		case tea.KeyCtrlE:
			log.Printf("%v\n", m.inputs)
			// m.Make()
			return m, nil
		}
	}

	cmd = m.updateInputs(msg)

	return m, cmd
}

func (m model) View() string {
	var sb strings.Builder

	for i := range m.inputs {
		sb.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			sb.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.idx == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&sb, "\n\n%s\n\n", *button)
	sb.WriteString(helpStyle.Render("cursor mode is "))
	sb.WriteString(cursorModeHelpStyle.Render(m.mode.String()))
	sb.WriteString(helpStyle.Render(" (ctrl+r to change style, esc to quit)"))

	return sb.String()
}

func (m model) Make() error {
	var name string
	var dir string
	var mod string
	var version string
	var release *github.RepositoryRelease
	var username string
	var password string

	// username = "imajinyun"
	// password = "L2rn/>GqfrW#2WTGt4JD"

	base := util.GetExecDir()
	name, mod = m.inputs[0].Value(), m.inputs[1].Value()
	isCurrDir := false

	{
		if util.Empty(name) {
			isCurrDir = true
		}

		dir = base
		if isCurrDir {
			name = filepath.Base(base)
		} else {
			dir = filepath.Join(base, name)
		}

		if isCurrDir {
			infos, err := os.ReadDir(dir)
			if err != nil {
				log.Printf("%v\n", err)
				return nil
			}

			cnt := 0
			for _, info := range infos {
				if info.Name()[0] != '.' {
					cnt++
				}
			}

			if cnt != 0 {
				log.Printf("%s is not empty\n", name)
				return nil
			}
		} else {
		}
	}

	{
		size := 1
		opts := &github.ListOptions{Page: 1, PerPage: size}
		client := github.NewClient(nil)
		releases, rsp, err := client.Repositories.ListReleases(context.Background(), "gohade", "hade", opts)
		log.Printf("%v\n", rsp.Rate.String())
		if err != nil {
		}

		for _, release := range releases {
			log.Printf("tag: %v\n", release.GetTagName())
			version = release.GetTagName()
		}
	}

	{
		var err error
		client := github.NewClient(&http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					req.SetBasicAuth(username, password)
					return nil, nil
				},
			},
		})

		if !util.Empty(version) {
			release, _, err = client.Repositories.GetReleaseByTag(context.Background(), "gohade", "hade", version)
			if err != nil || release == nil {
				log.Printf("release error: %v\n", err)
				return err
			}
			log.Printf("release tag: %v\n", release.GetTagName())
		} else {
			release, _, err = client.Repositories.GetLatestRelease(context.Background(), "gohade", "hade")
			if err != nil || release == nil {
				log.Printf("release error: %v\n", err)
				return err
			}
		}
	}

	tmp := filepath.Join(base, "gogin-"+version+"-"+cast.ToString(time.Now().Unix()))
	if err := os.Mkdir(tmp, os.ModePerm); err != nil {
		return err
	}

	url := release.GetZipballURL()
	err := util.DownloadFile(filepath.Join(tmp, "template.zip"), url)

	log.Printf("name: %v\n", name)
	log.Printf("dir: %v\n", dir)
	log.Printf("tmp: %v\n", tmp)
	log.Printf("mod: %v\n", mod)
	log.Printf("url: %v\n", url)
	log.Printf("tag: %v\n", release.GetTagName())
	log.Printf("version: %v\n", version)
	log.Printf("donwload: %v\n", "template.zip")

	if err != nil {
		log.Printf("download error: %v\n", err)
		return err
	}

	_, err = util.Unzip(filepath.Join(tmp, "template.zip"), tmp)
	if err != nil {
		return err
	}

	infos, err := os.ReadDir(tmp)
	if err != nil {
		return err
	}

	for _, info := range infos {
		log.Printf("info.name: %v\n", info.Name())
		if info.IsDir() && strings.Contains(info.Name(), "gohade-hade-") {
			if !isCurrDir {
				if err := os.Mkdir(dir, os.ModePerm); err != nil {
					return err
				}
			}

			if err := util.CopyDir(filepath.Join(tmp, info.Name()), dir); err != nil {
				return err
			}
		}
	}

	if err := os.RemoveAll(tmp); err != nil {
		return err
	}

	if err := os.RemoveAll(path.Join(dir, "github.com/imajinyun/goframe")); err != nil {
		return err
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		byt, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		if path == filepath.Join(dir, "go.mod") {
			byt = bytes.ReplaceAll(byt, []byte("module github.com/gohade/hade"), []byte("module "+mod))
			byt = bytes.ReplaceAll(byt, []byte("require ("), []byte("require (\n  github.com/gohade/hade "+version))
			if err := os.WriteFile(path, byt, 0o644); err != nil {
				return err
			}

			return nil
		}

		if bytes.Contains(byt, []byte("github.com/gohade/hade/app")) {
			byt = bytes.ReplaceAll(byt, []byte("github.com/gohade/hade/app"), []byte(mod+"/app"))
			if err := os.WriteFile(path, byt, 0o644); err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func (m model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func newModel() model {
	m := model{inputs: make([]textinput.Model, 3)}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Your project name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Your module name"
			t.CharLimit = 128
		case 2:
			t.Placeholder = "Select version"
			t.CharLimit = 1
		}
		m.inputs[i] = t
	}

	return m
}

var newCommand = &cobra.Command{
	Use:     "new",
	Short:   "Create a new Gogin application",
	Long:    "Create a new Gogin application",
	Aliases: []string{"init", "create"},
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := tea.LogToFile("./storage/log/gogin.log", "debug")
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}
		defer f.Close()

		p := tea.NewProgram(newModel())
		_, err = p.Run()
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}

		return nil
	},
}

func initNewCommand() *cobra.Command {
	return newCommand
}
