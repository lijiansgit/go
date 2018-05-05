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
	if len(cmdDir) == 1 {
		if err := os.Chdir(cmdDir[0]); err != nil {
			return res, err
		}
	}

	in := bytes.NewBuffer(nil)
	var command *exec.Cmd
	if runtime.GOOS == `windows` {
		command = exec.Command("cmd")
	} else {
		command = exec.Command("sh")
	}
	command.Stdin = in
	// go func() {
	// 	in.WriteString(fmt.Sprintf("%s\n", cmds))
	// 	// in.WriteString("exit\n")
	// }()
	in.WriteString(fmt.Sprintf("%s\n", cmds))
	var stdout, stderr bytes.Buffer
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
