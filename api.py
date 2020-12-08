import requests
import json

url = "http://send-svr:3000/send"

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
