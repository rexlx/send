<template>
  <div class="container-fluid rx">
    <div v-for="r in responses" :key=r.id class="card replies" :class="{ warn: r.good }">
        <div @click="$emit('focusDetails', r.reply)"
          class="card-body d-flex"
          :class="{ err: !r.good }">
          <div class="mr-auto p-2 mylink" @click="goTo(r)">{{ r.id }}</div>
          <div class="mr-auto p-2 greedy">{{ r.host }}</div>
          <div class="ml-auto p-2 trn" :title="r.reply_to" >{{ r.reply_to }}</div>
          <!-- {{ r.host }} | {{ r.reply_to }} <aside class="mylink" @click="goTo(r)">{{ r.id }}</aside> -->
        </div>
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
import router from "../router/index.js"
// import notie from 'notie'
// import { ref } from 'vue'

export default {
    props: ['responses', 'details'],
    setup() {
        Rules.hasToken()
        const goTo = async (r) => {
            router.push(`/admin/responses/${r.id}`)
        }
            return { store, goTo }
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
    color: honeydew;
    max-height: 480;
    overflow-y: scroll;
}

.trn {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.greedy {
    white-space: nowrap;
}

.replies {
    background-color: rgb(15, 17, 19);
    color: aliceblue;
}

.replies:hover {
    background: rgb(12, 36, 36);
}

.err {
    background-color: rgb(160, 122, 96);
}

.err:hover {
    background-color: rgb(180, 60, 30);
}

.mspace {
    font-family: monospace;
    color: whitesmoke;
}

.mylink {
    color: white;
    cursor: pointer;
}

</style>