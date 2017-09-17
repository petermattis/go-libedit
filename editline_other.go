// +build !darwin,!freebsd,!linux,!openbsd,!netbsd,!dragonfly

package libedit

import (
	"bufio"
	"os"
)

type state struct {
	reader                *bufio.Reader
	promptGenLeft         LeftPromptGenerator
	line                  string
	stdin, stdout, stderr *os.File
}

var editors []state

func Init(x string) (EditLine, error) {
	return InitFiles(x, os.Stdin, os.Stdout, os.Stderr)
}

func InitFiles(_ string, stdin, stdout, stderr *os.File) (EditLine, error) {
	st := state{
		reader: bufio.NewReader(stdin),
		stdin:  stdin, stdout: stdout, stderr: stderr,
	}
	editors = append(editors, st)
	return EditLine(len(editors) - 1), nil
}

func (el EditLine) Close()                                {}
func (el EditLine) SaveHistory() error                    { return nil }
func (el EditLine) AddHistory(_ string) error             { return nil }
func (el EditLine) LoadHistory(_ string) error            { return nil }
func (el EditLine) SetAutoSaveHistory(_ string, _ bool)   {}
func (el EditLine) UseHistory(_ int, _ bool) error        { return nil }
func (el EditLine) SetCompleter(_ CompletionGenerator)    {}
func (el EditLine) SetLeftPrompt(_ LeftPromptGenerator)   {}
func (el EditLine) SetRightPrompt(_ RightPromptGenerator) {}

func (el EditLine) Stdin() *os.File {
	return editors[el].stdin
}

func (el EditLine) Stdout() *os.File {
	return editors[el].stdout
}

func (el EditLine) Stderr() *os.File {
	return editors[el].stderr
}

func (el EditLine) GetLineInfo() (string, int) {
	return el.line, len(el.line)
}

func (el EditLine) GetLine() (string, error) {
	st := &editors[el]
	if st.promptGenLeft != nil {
		st.stdout.WriteString(st.promptGenLeft.GetLeftPrompt())
		st.stdout.Sync()
	}
	line, err := st.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	st.line = line
	return st.line, nil
}
