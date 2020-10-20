const express = require("express")             // express router
const fs = require("fs")                       // used to write the log
const bodyParser = require("body-parser")      // parses json
const spawn = require("child_process").spawn  // used to run system commands

// its common to call express app
const app = express()

// if you need to change the log name
const log = "api.log"

// tell express to use the parser middlewear
app.use(bodyParser.urlencoded({ extended: false }))
app.use(bodyParser.json())

// add our options method
app.options("/send", (req, res) => {
    // this is to stop that cors jazz
    res.header("Access-Control-Allow-Origin", "*")
    res.header("Access-Control-Allow-Method", "*")
    res.header("Access-Control-Allow-Headers", "*")
    res.end()
})

// add our post method
app.post("/send", (req, res) => {
    // again with the cors
    res.header("Access-Control-Allow-Origin", "*")
    res.header("Access-Control-Allow-Method", "*")
    res.header("Access-Control-Allow-Headers", "*")
    // validate the incoming json
    if (Object.keys(req.body).length < 4) {
        res.status(400)
        res.send(`expected 4 keys in object, got ${Object.keys(req.body).length}`)
    }
    // test properties of the body
    command = req.body
    if (command.ordered) {
        const s = spawn("/home/rxlx/bin/send", [`"${req.body.cmd}"`, `-u ${req.body.user}`, `-m "${req.body.host}"`,`-t ${req.body.timeout}`, `-o`], { shell: true })
                // stops a single stdout or err from ending the conn
        s.stdout.pipe(res, {end: false})
        s.stderr.pipe(res, {end: false})

        // its now safe to end
        s.on('close', () => {
            res.end()
        })
        console.log(command)
    } else {
        const s = spawn("/home/rxlx/bin/send", [`"${req.body.cmd}"`, `-u ${req.body.user}`, `-m "${req.body.host}"`,`-t ${req.body.timeout}`], { shell: true })
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
    let msg = `${date}\n  $ ${req.body.cmd} as ${req.body.user} on ${req.body.host}\n`
    fs.appendFile(log, msg, (e) => {
        if (e) throw e
    })
})


app.listen(3000, () => {

})
