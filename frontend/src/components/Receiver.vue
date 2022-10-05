<template>
  <div class="container-fluid rx">
    <table v-if="responses" class="table table-dark table-striped">
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
            </table>
  </div>
</template>

<script>
import {store} from './store.js'
import Rules from './rules.js'
// import notie from 'notie'
// import { ref } from 'vue'

export default {
    props: ['responses'],
    setup() {
        Rules.hasToken()
            return { store }
    }
}
</script>

<style scoped>

.rx {
    width: 50%;                                                                                                 
    height: 600px;
    /* height: 500px;*/
    float: right;                                                                                               
    background-color: rgb(10, 10, 14);                                                                                  
    border: 1px;                                                                                                
    border-color: black;
    color: lightskyblue;
    max-height: 480;
    /* overflow: auto; */
    overflow-y: scroll;
}

.replies {
    list-style: none;
}

.mspace {
    font-family: monospace;
    color: whitesmoke;
}

</style>