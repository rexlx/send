<template>
    <div class="container-fluid tx">
        <div class="row">
            <div class="col">
                <h3 class="mt-3 user-theme">console</h3>
            </div>
        </div>
        <hr>
        <div class="row">
            <div class="col">
                <div class="col">
                <select v-model="selectedConfig" class="config">
                    <option disabled value="">configs</option>
                    <option v-for="c in allConfigs" :value="c.id" :key="c.id">
                    {{c.name}}
                    </option>
                </select>
                <select v-model="selectedCommand" class="config">
                    <option disabled value="">commands</option>
                    <option v-for="c in savedCommands" :value="c.command" :key="c.id">
                    {{c.command_name}}
                    </option>
                </select>
            </div>
            </div>
        </div>
        <div class="row config-details">
            <div class="col details">
                config: {{ store.config.name }}
                <br>
                hosts: {{ store.config.hosts }}
            </div>
            <div class="col">
                timeout: {{ store.config.timeout }}
            </div>
            <div class="col">
                identity: {{ store.config.key }}
            </div>
        </div>
        <hr>
        <div class="row">
            <div class="input-group mb-3 terminal col">
                <input type="text" class="form-control name" v-model="commandName" placeholder="name">
                <div class="input-group-prepend">
                    <button @click="saveCommand" class="btn btn-dark" type="button">save</button>
                    <button @click="getCommands" class="btn btn-dark" type="button">load</button>
                </div>
            </div>
            <div class="input-group mb-3 terminal col">
                <div class="input-group-prepend">
                    <button @click="addToQueue" class="btn btn-dark" type="button">queue</button>
                    <button @click="send" class="btn btn-dark" type="button">send</button>
                </div>
                <input type="text"
                   class="form-control name"
                   ref="mainFocus"
                   placeholder="command"
                   v-model="command"
                   @keyup.enter="addToQueue"
                   @keyup.enter.shift.exact="send"
                   >
            </div>
        </div>
        <div class="row">
            <ul class="nav nav-tabs">
                <li class="nav-item">
                    <a class="nav-link" @click="tab=1" :class="{active : tab == 1}">queue</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" @click="tab=0" :class="{active : tab == 0}">history</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link disabled" :class="{active : tab == 2}" aria-disabled="true">details</a>
                </li>
            </ul>
            <!-- <div class="col">
                <button @click="tab=0" class="btn btn-dark" type="button">history</button>
                <button @click="tab=2" class="btn btn-dark" type="button">queue</button>
                <button @click="clearQueue" class="btn btn-dark" type="button">clear queue</button>
            </div> -->
        </div>
        <Receiver @focus-details="doThing" :responses="responses" />
        <Focus :tab="tab" :details="details" :q="q" />
        <div class="row tools text-center">
            <div class="col">
                <button @click="refreshDB" class="btn btn-outline-dark" type="button">bonk</button>
                <select v-model="resPerFetch" class="config">
                    <option disabled value="25">results/page</option>
                    <option value="50">50</option>
                    <option value="100">100</option>
                    <option value="150">150</option>
                    <option value="200">200</option>
                </select>
            </div>
        </div>
    </div>
</template>

<script>
import {store} from './store.js'
import { ref, watch, onMounted } from 'vue'
import notie from 'notie'
import Rules from './rules.js'
import Receiver from './Receiver.vue'
import Focus from './Focus.vue'

export default {
    components: { Focus, Receiver },
    setup() {
        const q = ref([])
        const socket = ref(null)
        const command = ref('')
        const commandName = ref('')
        const responses = ref([])
        const resPerFetch = ref(25)
        const savedCommands = ref([])
        const selectedCommand = ref('')
        const selectedConfig = ref('')
        const allConfigs = ref([])
        const details = ref('')
        const tab = ref(1)
        const mainFocus = ref(null)

        watch(selectedCommand, (currentValue) => {
            command.value = currentValue
            mainFocus.value.focus()
        })

        watch(resPerFetch, (currentValue) => {
            if (currentValue > 25) {
               resPerFetch.value = currentValue 
            }
        })

        watch(selectedConfig, (currentValue) => {
            changeConfig(currentValue)
            mainFocus.value.focus()
        })

        onMounted(() => {
            fetch(process.env.VUE_APP_API_URL + "/admin/responses/num/" + resPerFetch.value, Rules.requestOptions(""))
            .then((res) => res.json())
            .then((res) => {
                if (res.error) {
                    notie.alert({
                        type: "error",
                        text: res.message
                    })
                } else {
                    responses.value = res.data.data
                }
                }).catch((error) => {
                console.log(error)
            })

        fetch(process.env.VUE_APP_API_URL + "/admin/commands/get", Rules.requestOptions({user_id: store.user.id}))
            .then((res) => res.json())
            .then((res) => {
                if (res.error) {
                    notie.alert({
                        type: "error",
                        text: res.message
                    })
                } else {
                    savedCommands.value = res.data.data
                }
                }).catch((error) => {
                console.log(error)
            })

        fetch(process.env.VUE_APP_API_URL + "/admin/configs", Rules.requestOptions(""))
            .then((res) => res.json())
            .then((res) => {
                if (res.error) {
                    notie.alert({
                        type: "error",
                        text: res.message
                    })
                } else {
                    // console.log(res.data.configs)
                    allConfigs.value = res.data.configs
                }
                }).catch((error) => {
                console.log(error)
            })
        })

        const refreshDB = async () => {
            fetch(process.env.VUE_APP_API_URL + "/admin/responses/num/" + resPerFetch.value, Rules.requestOptions(""))
                .then((res) => res.json())
                .then((res) => {
                    if (res.error) {
                        notie.alert({
                            type: "error",
                            text: res.message
                        })
                    } else {
                        responses.value = res.data.data
                    }
                    }).catch((error) => {
                    console.log(error)
                })
        }

        const delayedRefresh = async () => {
            setTimeout(() => {
                fetch(process.env.VUE_APP_API_URL + "/admin/responses/num/" + resPerFetch.value, Rules.requestOptions(""))
                .then((res) => res.json())
                .then((res) => {
                    if (res.error) {
                        notie.alert({
                            type: "error",
                            text: res.message
                        })
                    } else {
                        responses.value = res.data.data
                    }
                    }).catch((error) => {
                    console.log(error)
                })
            }, 1250);
        }

        const clearQueue = async () => {
            q.value = []
            command.value = ""
        }

        const addToQueue = async () => {
            if (command.value !== "") {
                q.value.push(command.value)
            }
            command.value = ""
        }

        const changeConfig = async (num) => {
            fetch(process.env.VUE_APP_API_URL + "/admin/configs/get/" + num, Rules.requestOptions(""))
                .then((res) => res.json())
                .then((res) => {
                    if (res.error) {
                        notie.alert({
                            type: "error",
                            text: res.message
                        })
                    } else {
                        
                        store.config = res
                    }
                    }).catch((error) => {
                    console.log(error)
                })
        }

        const send = async () => {
            if (!store.config.hosts) {
                notie.alert({
                        type: "error",
                        text: "no hosts to send to (likely no config selected)"
                    })
                    return
            }
            if (command.value !== "") {
                q.value.push(command.value)
            } else {
                console.log(command.value)
                if (q.value.length === 0) {
                    notie.alert({
                        type: "error",
                        text: "no command to send"
                    })
                    return
                }
            }

            store.config.command = q.value.join(';')

            const data = {
                user: store.config.user,
                name: store.config.name,
                key: store.config.key,
                logpath: store.config.logpath,
                hosts: store.config.hosts,
                command: store.config.command,
                fromfile: store.config.fromfile,
                timeout: store.config.timeout,
                port: store.config.port,
                fatal: store.config.fatal,
                ordered: store.config.ordered,
            }
            q.value = []
            command.value = ""

            fetch(`${process.env.VUE_APP_API_URL}/admin/send`, Rules.requestOptions(data))
            .then((res) => res.json())
            .then((res) => {
                if (res.error) {
                    notie.alert({
                        type: "error",
                        text: res.message
                    })
                    // console.log(res.error)
                } else {
                    store.commandHistory.unshift(res.data.config)
                    // store.commandHistory.push(res.data.config)
                    delayedRefresh()
                }
                }).catch((error) => {
                notie.alert({
                    type: "error",
                    text: error
                })
                return
            })
            mainFocus.value.focus()
            notie.alert({
                    type: "success",
                    text: "command sent"
                })
        }

        const getConfigs = async () => {
            fetch(process.env.VUE_APP_API_URL + "/admin/configs", Rules.requestOptions(""))
            .then((res) => res.json())
            .then((res) => {
                if (res.error) {
                    notie.alert({
                        type: "error",
                        text: res.message
                    })
                } else {
                    allConfigs.value = res.data.data
                }
                }).catch((error) => {
                console.log(error)
            })
        }

        const saveCommand = async () => {
            const data = {
                command_name: commandName.value,
                user_id: store.user.id,
                command: q.value.join(';')
            }
            commandName.value = ""
            fetch(process.env.VUE_APP_API_URL + "/admin/save/command", Rules.requestOptions(data))
            .then((res) => res.json())
            .then((res) => {
                if (res.error) {
                    notie.alert({
                        type: "error",
                        text: res.message
                    })
                } else {
                    notie.alert({
                        type: "success",
                        text: "command was saved :)"
                    })
                }
                }).catch((error) => {
                console.log(error)
            })
        }

        const getCommands = async () => {
            const data = {
                user_id: store.user.id
            }
            fetch(process.env.VUE_APP_API_URL + "/admin/commands/get", Rules.requestOptions(data))
            .then((res) => res.json())
            .then((res) => {
                if (res.error) {
                    notie.alert({
                        type: "error",
                        text: res.message
                    })
                } else {
                    savedCommands.value = res.data.data
                }
                }).catch((error) => {
                console.log(error)
            })
        }

        const doThing = async (data) => {
            details.value = JSON.parse(data).message
            tab.value = 2
        }

        const wsConnect = async () => {
            let sck = new WebSocket("ws://localhost:8888/wsc")
            sck.onopen = () => {
                console.log("connected nice")
            }
        }

        return { 
            store,
            command,
            commandName,
            getCommands,
            addToQueue,
            clearQueue,
            allConfigs,
            send,
            savedCommands,
            q,
            details,
            getConfigs,
            responses,
            doThing,
            delayedRefresh,
            changeConfig,
            refreshDB,
            resPerFetch,
            saveCommand,
            selectedCommand,
            selectedConfig,
            tab,
            wsConnect,
            mainFocus,
            socket
            }
    }

}
</script>

<style scoped>

select {
    background: rgb(30, 30, 40);
}

.terminal {
    color: rgb(70, 140, 210);
    padding: 10px;
    margin: 1px;
    /* width: 50%; */
    float: right;
    border: 1px;
    border-color: black;
}

.config {
    color: rgb(70, 140, 210);
    padding: 10px;
    margin: 1px;
    float: left;
    border: 1px;
    border-color: rgb(30, 30, 40);
}

.single {
    color: rgb(100, 125, 150);
    margin: 4px;
    font-family: monospace;
}

.config-details {
    color: rgb(100, 125, 150);
    margin: 4px;
}

.name {
    padding: 4px;
    border: 1px;
    border-color: rgb(12, 14, 12);
    width: 30%;
    float: left;
}

.tools {
    float: right;
    margin: 2px;
    padding: 10px;
    border-radius: 2px;
    border: 1px solid black;
}

input[type=text] {
  background-color: rgb(29, 29, 39);
  color: rgb(196, 200, 202);
}
</style>