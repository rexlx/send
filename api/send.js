const express = require("express")
const fs = require("fs")
const cors = require("cors")
const path = require("path")
const bodyParser = require("body-parser")
const config = require("config")
const spawn = require("child_process").spawn

const log = config.get("api.log")
const sendLog = config.get("send.log")
const sendPath = config.get("api.sendBinary")
const key = config.get("send.key")
const uniq = config.get("send.unique")
const fatal = config.get("send.fatal")

const app = express()
app.use(bodyParser.urlencoded({ extended: false }))
app.use(bodyParser.json())
app.use(cors())
app.use("/data", express.static(path.join(__dirname, "data")))

// add our options method
// app.options("/send", (req, res) => {
//     // this is to stop that cors jazz
//     res.header("Access-Control-Allow-Origin", "*")
//     res.header("Access-Control-Allow-Method", "*")
//     res.header("Access-Control-Allow-Headers", "*")
//     res.end()
// })

// add our post method
app.post("/send", (req, res) => {
    // validate the incoming json
    if (Object.keys(req.body).length < 4) {
        res.status(400)
        res.send(`expected 4 keys in object, got ${Object.keys(req.body).length}`)
    }
    // test properties of the body
    command = req.body
    if (command.ordered) {
        const s = spawn(sendPath,
            [`-c "${req.body.cmd}"`,
            `-log ${sendLog}`,
            `-user ${req.body.user}`,
            `-hosts "${req.body.host}"`,
            `-port "${req.body.port}"`,
            `-timeout ${req.body.timeout}`,
            `-ordered`],
            { shell: true }
            )
        // stops a single stdout or err from ending the conn
        s.stdout.pipe(res, {end: false})
        s.stderr.pipe(res, {end: false})
        // its now safe to end
        s.on('close', () => {
            res.end()
        })
        console.log(command)
    } else {
        const s = spawn(sendPath,
            [`-c "${req.body.cmd}"`,
            `-log ${sendLog}`,
            `-user ${req.body.user}`,
            `-hosts "${req.body.host}"`,
            `-port "${req.body.port}"`
            `-timeout ${req.body.timeout}`],
            { shell: true }
            )
        // you may see other spawn examples using callback functions, this
        // stops a single stdout or err from ending the conn
        s.stdout.pipe(res, {end: false})
        s.stderr.pipe(res, {end: false})

        // its now safe to end
        s.on('close', () => {
            res.end()
        })
        }
        // console.log(command)

    let date = new Date()
    // log request
    let msg = `${date}-> running ${req.body.cmd} as ${req.body.user} on ${req.body.host}:"${req.body.port}`
    fs.appendFile(log, msg, (e) => {
        if (e) throw e
    })
})


app.listen(3000, () => {

})
