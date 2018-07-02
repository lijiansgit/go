package libs

import (
	"os/exec"
	"os"
	"strings"
	"bytes"
	"errors"
)

func Cmd(cmds string, cmdDir ...string) (res string,err error) {
	if len(cmdDir) == 1 {
		if err := os.Chdir(cmdDir[0]); err != nil {
			return res, err
		}
	}

	var (
		cmd string
		cmdArgs, args []string
	)
	cmdArgs = strings.Split(cmds, " ")
	cmd = cmdArgs[0]
	args = cmdArgs[1:]
	if strings.Count(cmds, "'") == 2 {
		args = rsyncArgs(cmdArgs[1:])
	}
	command := exec.Command(cmd, args...)
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

	return	stdout.String(), nil
}

// 处理命令行参数中有带双引号子参数的参数,比如rsync中:'-e ssh -p 22'，只能处理一对双引号
func rsyncArgs(args []string) []string {
	var (
		n, m int
		str string
		args1, args2, args3 []string
	)
	for k, v := range args {
		if strings.HasPrefix(v, "'") == true {
			n = k
			args1 = args[:n]
		}
		if strings.HasSuffix(v, "'") == true {
			m = k + 1
			args3 = args[m:]
		}
	}
	args2 = args[n:m]
	str = strings.Join(args2, " ")	//'-e ssh -p 22'
	str = strings.Trim(str, "'") //-e ssh -p 22
	args = args1
	args = append(args, str)
	args = append(args, args3...)
	return args
}
