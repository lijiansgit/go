package libs

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Cmd(cmds string, cmdDir ...string) (res string, err error) {
	var (
		command               *exec.Cmd
		stdout, stderr, stdin bytes.Buffer
	)

	if len(cmdDir) == 1 {
		if err := os.Chdir(cmdDir[0]); err != nil {
			return res, err
		}
	}

	if runtime.GOOS == `windows` {
		command = exec.Command("cmd")
	} else {
		command = exec.Command("sh")
	}

	// stdin = bytes.NewBuffer(nil)
	stdin.WriteString(fmt.Sprintf("%s\n", cmds))
	command.Stdin = &stdin
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		if stderr.String() == "" {
			return res, err
		}

		err = errors.New(stderr.String())
		return res, err
	}

	return stdout.String(), nil
}
