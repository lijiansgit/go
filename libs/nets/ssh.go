// ssh执行命令，支持密码/密钥登陆执行

package nets

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Ssh struct {
	PrivateKey string
	User       string
	Password   string
	Addr       string
	Port       string
}

// NewSsh 创建新的ssh
// 优先使用privateKey免密码登陆，如果为空，则使用密码登陆
// TODO ssh.Dial 和 cmd 分开
func NewSsh(privateKey, user, password, addr, port string) *Ssh {
	return &Ssh{
		PrivateKey: privateKey,
		User:       user,
		Password:   password,
		Addr:       addr,
		Port:       port,
	}
}

func (s *Ssh) Cmd(cmd string) (res string, err error) {
	var (
		config *ssh.ClientConfig
	)

	if s.PrivateKey != "" {
		var hostKey ssh.PublicKey
		key, err := ioutil.ReadFile(s.PrivateKey)
		if err != nil {
			return res, err
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return res, err
		}

		hostKey, err = s.getHostKey()
		if err != nil {
			return res, err
		}

		config = &ssh.ClientConfig{
			User: s.User,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.FixedHostKey(hostKey),
			Timeout:         10 * time.Minute,
		}
	} else {
		passWord := []ssh.AuthMethod{ssh.Password(s.Password)}
		config = &ssh.ClientConfig{
			User:            s.User,
			Auth:            passWord,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         10 * time.Minute,
		}
	}

	addr := s.Addr + ":" + s.Port
	client, err := ssh.Dial("tcp", addr, config)
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

func (s *Ssh) getHostKey() (ssh.PublicKey, error) {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	// file, err := os.Open("/root/.ssh/known_hosts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], s.Addr) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error parsing %q: %v", fields[2], err))
			}
			break
		}
	}

	if hostKey == nil {
		return nil, errors.New(fmt.Sprintf("no hostkey for %s", s.Addr))
	}
	return hostKey, nil
}

// 去除字符串中的cutset
func StringTrim(str, cutset string) (res string) {
	for _, v := range cutset {
		str = strings.Replace(str, string(v), "", -1)
	}

	return str
}
