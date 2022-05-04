package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	conf     = flag.String("conf", "", "json config location")
	cmd      = flag.String("c", "", "specify command")
	user     = flag.String("user", "", "username")
	keyFile  = flag.String("key", "", "path to key if not default")
	hosts    = flag.String("hosts", "", "multiple hosts inside quotes")
	fromFile = flag.String("file", "", "specify a file separated by newline")
	logPath  = flag.String("log", "", "for flat file: /path/name.ext, syslog: `PROTO@ADDR:PORT`")
	timeout  = flag.Int("timeout", 90, "timeout in seconds")
	port     = flag.Int("port", 22, "port number")
	fatal    = flag.Bool("fatal", false, "return failed exit codes")
	ordered  = flag.Bool("ordered", false, "run in order instead of async")
	usr      string
	dest     string
	sshCfg   string
	cfg      Config
)

// Config gets datafilled by either 1. cli args, or 2. a specified json config
type Config struct {
	User     string `json:"user"`
	Key      string `json:"key"`
	LogPath  string `json:"logpath"`
	Hosts    string `json:"hosts"`
	Command  string `json:"cmd"`
	FromFile string `json:"fromfile"`
	Timeout  int    `json:"timeout"`
	Port     int    `json:"port"`
	Fatal    bool   `json:"fatal"`
	Ordered  bool   `json:"ordered"`
}

func main() {
	response := make(chan string)
	// parse the flags before anything
	flag.Parse()

	// if automated using a json config
	if *conf != "" {
		configFromJSON(*conf)
	} else {
		configFromArgs(&cfg)
	}

	logHandler(cfg.LogPath)
	// log for records should this be DEBUG?
	log.Printf("%v |-> using the following config |-> %v\n", "__SEND__", cfg)
	// hosts should be supplied like -host "host1 user2@host2 roooooot@topCkid"
	for _, host := range strings.Split(cfg.Hosts, " ") {
		// if they added a username@host, disregard c.User
		if strings.Contains(host, "@") {
			usr = strings.Split(host, "@")[0]
			dest = strings.Split(host, "@")[1]
		} else {
			usr = cfg.User
			dest = host
		}
		if cfg.Ordered {
			response := sendCmd(&cfg, usr, dest)
			fmt.Println(response)
		} else {
			go func(dest, usr string) {
				response <- sendCmd(&cfg, usr, dest)
			}(dest, usr)
		}
	}
	if !cfg.Ordered {
		for i := 0; i < len(strings.Split(cfg.Hosts, " ")); i++ {
			select {
			case results := <-response:
				fmt.Println(results)
			case <-time.After(time.Duration(cfg.Timeout) * time.Second):
				log.Printf("__TIMEOUT__ |-> %v timed out; timeout: %v", dest, cfg.Timeout)
				return
			}
		}
	}
}

func sendCmd(c *Config, usr, dest string) string {
	// init some vars to hold our stdour/err later
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// create an ssh conf using our loadKey function / user key
	sshConf := &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{loadKey(c.Key)}, HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%v", dest, cfg.Port), sshConf)
	if err != nil {
		fmt.Printf("error when attempting to dial %v as user: %v!\n%v\n", dest, usr, err)
		connErrMsg := fmt.Sprintf("error connecting to %v on port %v as %v\n", dest, cfg.Port, usr)
		log.Println(connErrMsg, err)
	}
	session, sessionErr := conn.NewSession()
	if sessionErr != nil {
		// TODO: need to add silent mode for procs running in BG
		fmt.Printf("session error! with %v on %v as user %v\n", dest, cfg.Port, usr)
		log.Printf("session error! with %v on %v as user %v\n", dest, cfg.Port, usr)
	}
	defer session.Close()

	// store the results of the command that was sent
	session.Stdout = &stdout
	session.Stderr = &stderr
	// log what we did for our records
	log.Printf("running %v on %v as user %v", cfg.Command, dest, usr)
	session.Run(cfg.Command)
	// if there was stderr, something went wrong
	if len(stderr.String()) > 0 {
		// whether we care or not is determined here. If set to fatal,
		// any stderr returned will terminate the program and return failed (os exit 1)
		if cfg.Fatal {
			fmt.Printf("Cant recover, fatal set to true. host affected: %v |->%v\n", dest, stderr.String())
			log.Printf("Cant recover, fatal set to true. host affected: %v |->%v\n", dest, stderr.String())
			os.Exit(1)
		}
		// otherwise tell us what failed and carry on with work
		log.Printf("__FAIL__ |-> command failed on %v->%v", dest, stderr.String())
		return fmt.Sprintf("%v", stderr.String())
	}
	return fmt.Sprintf("%v", stdout.String())
}

func logHandler(logPath string) {
	// if the user supplies (what we define as a) syslog path, unpack
	// -log tcp@hostname:port | -log udp@addr:port
	if strings.Contains(logPath, "@") {
		addr := strings.Split(logPath, "@")
		logger, e := syslog.Dial(addr[0], addr[1],
			syslog.LOG_WARNING|syslog.LOG_DAEMON, "__SEND__") // anything else here
		check(e)
		log.SetOutput((logger))
	} else {
		// otherwise user supplied a path (or fat fingered something)
		// -log /path/to/flatFile.txt
		logger, e := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		check(e)
		defer logger.Close()
		log.SetOutput(logger)
	}
}

func configFromJSON(conf string) {
	contents, e := ioutil.ReadFile(conf)
	check(e)
	// this is where we unmarshal the contents into cfg (Config)
	e = json.Unmarshal(contents, &cfg)
	check(e)
}

func configFromArgs(c *Config) {
	defaultDir, e := os.UserHomeDir()
	check(e)
	c.Port = *port
	c.Timeout = *timeout
	c.Ordered = *ordered
	c.Fatal = *fatal
	if *user == "" {
		c.User = strings.Split(defaultDir, "/")[2]
	} else {
		c.User = *user
	}
	if *keyFile == "" {
		// we could probably loop over an array and concat the algo name to `id_`
		// for now we only support RSA
		c.Key = fmt.Sprintf("%v/.ssh/id_rsa", defaultDir)
	} else {
		c.Key = *keyFile
	}
	if *logPath == "" {
		c.LogPath = "send.log"
	} else {
		c.LogPath = *logPath
	}
	if *hosts == "" {
		log.Printf("no hosts were provided in cli args! --> %v", *hosts)
		getHelp()
	} else {
		c.Hosts = *hosts
	}
	if *cmd != "" {
		c.Command = *cmd
	} else {
		c.Command = "whoami;hostname;w;uptime"
	}
	if *fromFile != "" {
		c.FromFile = *fromFile
	}
}

func check(e error) {
	if e != nil {
		log.Fatalf("send encountered an error!\t%+v\n", e)
	}
}

func getHelp() {
	fmt.Fprintf(os.Stderr, "usage of send:\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func loadKey(path string) ssh.AuthMethod {
	// contents is a bytes array
	contents, e := ioutil.ReadFile(path)
	check(e)
	key, err := ssh.ParsePrivateKey(contents)
	check(err)
	return ssh.PublicKeys(key)
	// c.Auth = ssh.PublicKeys(key)
}
