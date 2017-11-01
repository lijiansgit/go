package libs

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"errors"
	"strings"
	"time"
)

type Ssh struct {
	User     string
	Password string
	Addr     string
}

func (s *Ssh) Cmd(cmd string) (res string,err error) {
	passWord := []ssh.AuthMethod{ssh.Password(s.Password)}
	conf := &ssh.ClientConfig{User: s.User, Auth: passWord, HostKeyCallback: ssh.InsecureIgnoreHostKey(), Timeout: 20 * time.Minute}
	client, err := ssh.Dial("tcp", s.Addr, conf)
	if err != nil {
		return res, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return res, err
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	if err := session.Run(cmd); err != nil {
		if stderr.String() == "" {
			return res, err
		}
		return res, errors.New(stderr.String())
	}

	res = stdout.String()
	return res, nil
}

func NewSsh(user, password, addr string) *Ssh {
	return &Ssh{User: user, Password: password, Addr: addr}
}

// 去除字符串中的cutset
func StringTrim(str, cutset string) (res string) {
	for _, v := range cutset {
		str = strings.Replace(str, string(v), "", -1)
	}

	return str
}