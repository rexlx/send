<template>
  <div class="container-fluid rx">
    <div v-for="r in responses" :key=r class="card replies" :class="{ warn: r.good }">
        <div @click="doThing" class="card-body">{{ r.id }} | {{ r.host }} | {{ r.reply_to }}</div>
    </div>
    <!-- <table v-if="responses" class="table table-dark table-striped">
            <thead>
                <tr>
                    <th>id</th>
                    <th>host</th>
                    <th>time in</th>
                    <th>result</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="r in responses" v-bind:key="r.id">
                    <td>{{ r.id }}</td>
                    <td>
                        <router-link :to="`/admin/responses/${r.id}`">{{r.host}}</router-link>
                    </td>
                    <td>{{ r.time_rx }}</td>
                    <td v-if="r.good">
                        <span class="badge bg-success mspace">pass</span>
                    </td>
                    <td v-else>
                        <span class="badge bg-danger mspace">fail</span>
                    </td>
                </tr>
            </tbody>
            </table> -->
  </div>
</template>

<script>
import {store} from './store.js'
import Rules from './rules.js'
// import notie from 'notie'
// import { computed } from 'vue'

export default {
    props: ['responses'],
    setup() {
        Rules.hasToken()
        const doThing = async () => {
            console.log('i did a thing')
        }
            return { store, doThing }
    }
}
</script>

<style scoped>

.rx {
    width: 50%;                                                                                                 
    height: 600px;
    float: right;                                                                                               
    background-color: rgb(10, 10, 14);                                                                                  
    border: 1px;                                                                                                
    border-color: black;
    color: lightskyblue;
    max-height: 480;
    overflow-y: scroll;
}

.replies {
    background-color: rgb(20, 20, 40);
    cursor: pointer;
}

.err {
    background-color: olivedrab;
}

.mspace {
    font-family: monospace;
    color: whitesmoke;
}

</style>