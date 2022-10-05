<template>
  <div class="container-fluid log">
    <div v-if="store.commandHistory.length > 0">
        <div v-for="cfg in store.commandHistory" :key="cfg" >
            <div>
                {{ cfg.command }}
            </div>
        </div>
    </div>
    <div v-else class="container-fluid">
        <div class="col">
            <p><small>command history empty</small></p>
        </div>
    </div>
  </div>
</template>

<script>
import {store} from './store.js'
import Rules from './rules.js'
import notie from 'notie'
import { ref } from 'vue'

export default {
    setup() {
        const reply = ref("")
        Rules.hasToken()

        const fallback = async () => {
            fetch(process.env.VUE_APP_API_URL + "/admin/config", Rules.requestOptions(""))
            .then((res) => res.json())
            .then((res) => {
                if (res.error) {
                    notie.alert({
                        type: "error",
                        text: res.message
                    })
                } else {
                    //--:REX
                    store.config = res.data.config
                }
                }).catch((error) => {
                console.log(error)
            })
        }
            return { store, reply, fallback }
    }
}
</script>

<style scoped>

.host {
    overflow-y: scroll;
    color: rgb(117, 106, 90);
}

.log {
    font-family: monospace;
    width: 50%;                                                                                                 
    height: 500px;                                                                                              
    float: left;                                                                                               
    background-color: rgb(10, 10, 14);                                                                                  
    border: 1px;                                                                                                
    border-color: black;
    /* color: rgb(148, 122, 180); */
    color: rgb(104, 121, 120);
    overflow-y: scroll;
}


</style>