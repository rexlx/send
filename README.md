Send is a an automation tool that works over ssh. It consists of three main components:
<hr>

1. the send binary (gives you access to the send command if added to your path)
2. the api (acts as a broker between the http interface and the send binary. requires node/npm)
3. the http interface 

# Packaging
as this project continues to grow, I aim to make packaging easier and easier. right now I assume you know a few things:
<br>
1. how to use git to clone the repository
2. how to install packages on your operating system (go, node, and npm in this case)
<br>

# initial setup
## files to modify:
1. send/api/config/default.json contains information that node will use during runtime. the api section is used for api configuration, such as logging path, and the location of the send binary. the send section contains information pertinent the send configuration which the end user cant control. Like which key they can use and where send logs to.
2. send/api/data/nodes.json is where the http interface looks when it populates the dropdown
3. send/frontend/* is where you can modify the stylesheet/html and make any required configuration changes to index.html and main.js

## additional config
1. send is a passwordless application, ssh keys must be traded with all target machines
2. once compiled, send will need to be added to your path.
3. if you're running the api, it needs to know the path as well

# send (api)
this is what the api expects. the values serve as examples, but the keys matter.
```javascript
let command = {
  "cmd": "uptime",
  "user": "sadmin",
  "host": null,
  "port": 22,
  "timeout": "120",
  "ordered": false
}
```
please see the provided "api.py" for a working example
<br>

# send (cli)
send commands to remote machines with golang
```bash
rxlx ~ $ send -h
Usage of send:
  -c string
    	specify command
  -conf string
    	json config location
  -fatal
    	return failed exit codes
  -file string
    	specify a file separated by newline
  -hosts string
    	multiple hosts inside quotes
  -key string
    	path to key if not default
  -log PROTO@ADDR:PORT
    	for flat file: /path/name.ext, syslog: PROTO@ADDR:PORT
  -ordered
    	run in order instead of async
  -port int
    	port number (default 22)
  -timeout int
    	timeout in seconds (default 90)
  -user string
    	username
```
