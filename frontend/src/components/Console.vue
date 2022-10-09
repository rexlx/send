<template>
    <div class="container-fluid tx">
        <div class="row">
            <div class="col">
                <h3 class="mt-3 user-theme">console</h3>
            </div>
        </div>
        <hr>
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
                    <button @click="saveCommand()" class="btn btn-outline-light" type="button">save</button>
                    <button @click="getCommands" class="btn btn-outline-light" type="button">load</button>
                </div>
            </div>
            <div class="input-group mb-3 terminal col">
                <div class="input-group-prepend">
                    <button @click="addToQueue" class="btn btn-outline-light" type="button">queue</button>
                    <button @click="send" class="btn btn-outline-light" type="button">send</button>
                </div>
                <input type="text"
                   class="form-control name"
                   ref="prompt"
                   placeholder="command"
                   v-model="command"
                   @keyup.enter="addToQueue"
                   @keyup.enter.shift.exact="send"
                   >
            </div>
        </div>
        <div class="row">
            <div class="col">
                <select v-model="selectedConfig" class="config">
                    <option disabled value="">configs</option>
                    <option v-for="c in allConfigs" :value="c.id" :key="c.id">
                    {{c.name}}
                    </option>
                </select>
            <div class="col">
                <select v-model="selectedCommand" class="config">
                    <option disabled value="">commands</option>
                    <option v-for="c in savedCommands" :value="c.command" :key="c.id">
                    {{c.command_name}}
                    </option>
                </select>
            </div>
            </div>
        </div>
        <div class="row">
            <div class="col">
                <div v-if="command" class="single">sending: {{ command }}</div>
                <div v-else-if="q.length > 0" class="single">
                    sending:
                    <div v-for="i in q" :key="i" >
                        <span class="single">{{ i }}</span>
                    </div>  
                </div>
            </div>
        </div>
        <Receiver @focus-details="doThing" :responses="responses" />
        <Focus :details="details" />
        <div class="row tools text-center">
            <div class="col">
                <select v-model="resPerFetch" class="config">
                    <option disabled value="25">25</option>
                    <option value="50">50</option>
                    <option value="50">100</option>
                    <option value="50">150</option>
                    <option value="50">200</option>
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

        watch(selectedCommand, (currentValue) => {
            command.value = currentValue
        })

        watch(resPerFetch, (currentValue) => {
            resPerFetch.value = currentValue
        })

        watch(selectedConfig, (currentValue) => {
            changeConfig(currentValue)
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
            }, 2250);
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
            if (command.value !== "") {
                q.value.push(command.value)
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
                reply_to: store.config.command.slice(0, 120)
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
            })
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
        }

        const wsConnect = async () => {
            let sck = new WebSocket("ws://storage.nullferatu.com:8888/wsc")
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
            wsConnect,
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
    padding: 20px;
    margin: 1px;
    width: 50%;
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

.active {
    background: rgba(70, 140, 210, 0.600);
}

.tools:hover {
    background: rgb(12, 12, 36);
}

input[type=text] {
  background-color: rgb(29, 29, 39);
  color: rgb(196, 200, 202);
}
</style>