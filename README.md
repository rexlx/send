# initial setup
## files to modify:
1. js/index.html:
div with id home, modify the address and port to the server hosting the front end ( you can use the provided       bin/serve)
2. js/main.js
const url and api need to reflect your setup
3. node/send.js
both const s = spawn() needs to know the actual path to your send binary

## additional config
1. send doesnt handle passwords currently, ssh keys must be traded with all target machines
2. i allow * in my cors policy, this may not be desirable to you, and if you understand what this message means, you know how to fix it...i think
3. while send.js should run fine, you should remove node_modules and npm install.

## troubleshooting
depending on whether or not you're using serve as the http server, you will have 2-3 logs to troubleshoot from
1. api.log -> logs what request it recieved and when
2. send.log -> logs what the send binary was sent, and whether or not it was able to connect
3. serve.log -> logs incoming http requests

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
