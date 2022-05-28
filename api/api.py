import requests
import json

url = "http://send-svr:3000/send"

my_commands = ["df -h", "uptime", "hostname", "sleep 3", "last"]

command = {
  "cmd": ";".join(my_commands),
  "user": "rxlx",
  "host": "svr",
  "port": 22,
  "timeout": "120",
  "ordered": "false"
}

headers = {'content-type': 'application/json'}

try:
  cmd = json.dumps(command)
  res = requests.post(url, data=cmd, headers=headers)
  print(f"request sent...got a {res.status_code}")
except Exception as e:
  print(f"encountered an error! {e}")
