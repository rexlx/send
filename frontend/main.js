// edit these ass needed
const url = "http://192.168.1.45:8080/nodes.json"
const api = "http://192.168.1.45:3000/send"

// init some vars to manipulate elsewhere
let nodes
let svr
let nodeList
let opnEl = document.createElement("option")

// technically this is an object but it will converted to json later
// consider these default vals that are tested for and manipulated as
// needed
let command = {
  "cmd": "uptime",
  "user": "rxlx",
  "host": null,
  "timeout": "120",
  "ordered": false
}

// this function simply gets the nodes.json which is passed to the drop
// down menu. these are effectively different groups of servers
const getNodes = async (url) => {
  // this is a common way to use the fetch api (http methods) with out
  // blocking the program
  let res = await fetch(url)
  if (res.status === 200) {
    const data = await res.json()
    return data
  }
  else {
    throw new Error(`! got status code: ${res.status}`)
  }
}

// here create the function that sends the command
const sendCmd = async (api, data) => {
  let res = await fetch(api, {
    method: "POST",
    mode: "cors", // we have to use cors mode to modify headers :/
    headers: {
      'Content-Type': 'application/json'
    },
    body: data
  })
  return res.text() // pass back text and not json, this is what ends up in the html
}

getNodes(url).then((data) => {
  nodes = data
  Object.keys(data).forEach((item) => {
    opnEl.textContent = item
    document.getElementById("svr-list").innerHTML += `<option>${item}</option>`
    // console.log(item)
  })
}).catch((err) => {
  console.log(err)
})

document.getElementById("ordered").addEventListener("click", (e) => {
  if (command.ordered) {
    command.ordered = false
    document.getElementById("ordered").innerHTML = "unordered"
  } else if (!command.ordered) {
      command.ordered = true
      document.getElementById("ordered").innerHTML = "ordered"
  }
  console.log(command.ordered)
})

document.getElementById("svr-list").addEventListener("change", (e) => {
  svr = e.target.value
  nodeList = nodes[svr]
})

document.getElementById("display").addEventListener("change", (e) => {
  // anytime we get a "change" on the checkbox, we want to first delete
  // so were not appending endlessly to thte html
  document.getElementById("servers").innerHTML = ""
  nodeList.forEach((item) => {
    document.getElementById("servers").innerHTML += `${item} <br>`
  })
})

document.getElementById("user").addEventListener("input", (e) => {
  command.user = e.target.value
  // console.log(command.user)
})

document.getElementById("timeout").addEventListener("input", (e) => {
  document.getElementById("timeout-label").innerHTML = "timeout"
  let int = parseInt(e.target.value)
  if (isNaN(int)) {
    document.getElementById("timeout-label").innerHTML = "must be an number in seconds"
  } else {
    command.timeout = int
  }
})

document.getElementById("command").addEventListener("keypress", (e) => {
  if (e.key === "Enter") {
    document.getElementById("output").innerHTML = ""
    // here we modify of json which is passed to the api
    command.host = nodeList.join(" ")
    command.cmd = e.target.value
    myJSON = JSON.stringify(command)
    sendCmd(api, myJSON).then((data) => {
      document.getElementById("output").innerHTML = data
    }).catch((err) => {
      console.log(err)
    })
  }
})
