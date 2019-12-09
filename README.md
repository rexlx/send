# send
send command to  remote machine with golang
```bash
send remote commands over ssh. works on MS or Linux (need to compile for both)
for now password authentication is NOT supported, ssh key only

usage:
POSIX:
$ send command host [args]

MS:
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
