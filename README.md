# What is this?
Send proxies requests to an endpoint that houses keys to remote nodes. The idea is that a trusted user can send requests to perform tasks without passwords.

## who is this for?
everyone.
<br>
the cli tool is for the engineer who loves the sound of successive keyboard strokes
<br>
the api is for the programmer who knows how to use the tool better than me
<br>
the web app is for the person who just needs to perform some tasks as needed

# Packaging
as this project continues to grow, I aim to make packaging easier and easier. right now I assume you know a few things:
<br>
1. how to use git
2. how to run a binary and/or add it to your path
3. you know how to run a node app
4. you have a browser that can render the html or have an http server to host it out.
5. you need to access nodes.json over http. this can be at localhost or served out.
<br>

![example](https://storage.googleapis.com/rfitzhugh/send01.png)

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
<br>
for python:


<br>
```
import requests
import json

url = "http://192.168.1.42:3000/send"

my_commands = ["df -h", "uptime", "hostname", "sleep 3", "last"]

command = {
  "cmd": ";".join(my_commands),
  "user": "rxlx",
  "host": "svr",
  "timeout": "120",
  "ordered": "false"
}

headers = {'content-type': 'application/json'}

cmd = json.dumps(command)
res = requests.post(url, data=cmd, headers=headers)

if res.status_code == 200:
    print(res.text)
else:
    print(f"requst sent, but i got the response: {res.status_code}")
    print(res.text)
```
<br>

# send (cli)
send commands to  remote machines with golang
```bash
send remote commands over ssh. works on MS or Linux (need to compile for both)
for now password authentication is NOT supported, ssh key only

usage:
POSIX:
$ send command host [args]

MS: **NOTE** if you add send.exe to the path, you can use just "send"
> send.exe command host [args]

specifcy a log name:
$ send "sudo updatedb" host1 -l logs/some.log

send a commnad to multiple hosts:
$ send "df -h" -m "host1 host2 host3"

specify different host names:
send --list-python -m "rxlx@rxlx rfitz@surx"

send a command to a specific user using a specific key:
$ send "locate special.xml" user@host -k /path/to/key

arguments:
-u    user name
-s    supress stdout
-k    specify key path
-p    ssh port
-m    multiple hosts: -m "host1 host2 host3"
-f    read hosts from file separated by new line
-t    command timeout in seconds (default is 120)
-l    logfile name (default is send.log)
-o    execute in order instead of asynchronously

optional commands:
--list-python  show cpu usage of all python processes
--list-perl    show cpu usage of all perl processes
--top-ten      show top ten processes by cpu
```
