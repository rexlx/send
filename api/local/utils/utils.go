package utils

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/rexlx/vapi/local/data"
	"golang.org/x/crypto/ssh"
)

type GenericResponse struct {
	Good    bool        `json:"error"`
	Message string      `json:"message"`
	TimeTX  time.Time   `json:"time_tx"`
	TimeRX  time.Time   `json:"time_rx"`
	Host    string      `json:"host"`
	Config  data.Config `json:"data"`
}

func SendSSHCmd(dest, user string, c *data.Config, t time.Time) GenericResponse {
	// init some vars to hold our stdour/err later
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// create an ssh conf using our loadKey function / user key
	sshConf := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{loadKey(c.Key)}, HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%v", dest, c.Port), sshConf)
	if err != nil {
		return GenericResponse{
			Good:    false,
			Message: fmt.Sprintf("error connecting to %v on port %v as %v: %v", dest, c.Port, user, err),
			Host:    dest,
		}

	}
	session, sessionErr := conn.NewSession()
	if sessionErr != nil {
		return GenericResponse{
			Good:    false,
			Message: fmt.Sprintf("session error! with %v on %v as user %v\n", dest, c.Port, user),
			Host:    dest,
		}
	}
	//--:REX this might close when called...
	defer session.Close()

	// store the results of the command that was sent
	session.Stdout = &stdout
	session.Stderr = &stderr
	session.Run(c.Command)
	// if there was stderr, something went wrong
	if len(stderr.String()) > 0 {
		// whether we care or not is determined here. If set to fatal,
		// any stderr returned will terminate the program and return failed (os exit 1)
		if c.Fatal {
			return GenericResponse{
				Good:    false,
				Message: "fatal feature not fully implemented",
				Host:    dest,
			}
		}
		return GenericResponse{
			Good:    false,
			Message: stderr.String(),
			TimeTX:  t,
			TimeRX:  time.Now(),
			Host:    dest,
			Config:  *c,
		}
	}
	return GenericResponse{
		Good: true,
		// Message: strings.Join(strings.Fields(stdout.String()), " "),
		Message: stdout.String(),
		TimeTX:  t,
		TimeRX:  time.Now(),
		Host:    dest,
		Config:  *c,
	}
}

func loadKey(path string) ssh.AuthMethod {
	// contents is a bytes array
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	key, err := ssh.ParsePrivateKey(contents)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}
