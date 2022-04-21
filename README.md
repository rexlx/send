Send is a an automation tool that works over ssh. It consists of three main components: 1. the send binary (add to your path) 2. the api (will need node and npm) and 3. (if serving out the frontend), an http server (the provided "serve" code / binary is for testing purposes only)

# update
I have added support for central logging to src/send.go but its still not implemented everywhere. support for reading in json configuration also added
changes to documentation coming.

# Packaging
as this project continues to grow, I aim to make packaging easier and easier. right now I assume you know a few things:
<br>
1. how to use git
2. how to run a binary and/or add it to your path
3. you know how to run a node app
4. install / configure an http server
<br>

# initial setup
## files to modify:
1. frontend/index.html:
div with id help, modify the address and port to the server hosting the front end, otherwise the help page wont work (you can use the provided bin/serve, an http server written in go.)
2. frontend/main.js
const url and api need to reflect your setup, refer to data/nodes.json for an example of how to define your target machines
3. api/send.js
both const s = spawn() needs to know the actual path to your send binary (this is the tool that sends commands over ssh)
you may want to specify const log (cuurently ./api.log)

## additional config
1. send is a passwordless application, ssh keys must be traded with all target machines
2. i allow * in my cors policy, this may not be desirable to you, and if you understand what this message means, you know how to fix it...i think
3. while api/send.js should run fine, you should remove node_modules and npm install.
4. you'll probably want to change the default user in frontend/main.js, as i doubt you also use the *rxlx* username :)

## troubleshooting
depending on whether or not you're using serve as the http server, you will have 2-3 logs to troubleshoot from
1. api.log -> logs what request it recieved and when, written by node/send.js
2. send.log -> logs what the send binary was sent, and whether or not it was able to connect, unless otherwise specified, it lives in the working directory of node/send.js
3. serve.log -> logs incoming http requests, unless specified otherwise, lives in the ./ dir it was called

# send (api)
the node app send.js is configured to listen on "http://ADDR_HERE:3000/send" for a post containing the follwing data (json):
```javascript
let command = {
  "cmd": "uptime",
  "user": "rxlx",
  "host": null,
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
