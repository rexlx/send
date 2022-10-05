<template>
  <div class="container-fluid">
    <div v-if="response">
        <h5>{{ response.host }}</h5>
        <hr>
        <span>roundtrip time to run:{{ runtime }}</span>
        <p class="response">{{ cleanedRes.message }}</p>
    </div>
    <div v-else class="container-fluid">
        <div class="col">
            <p><small>no response selected...</small></p>
        </div>
    </div>
  </div>
</template>

<script>
import {store} from './store.js'
import Rules from './rules.js'
// import notie from 'notie'
import { ref, computed } from 'vue'

export default {
    beforeMount() {
        Rules.hasToken();

        if (parseInt(String(this.$route.params.responseId), 10) > 0) {
            //edit existing user
            fetch(process.env.VUE_APP_API_URL + "/admin/responses/get/" + this.$route.params.responseId, Rules.requestOptions(""))
            .then((res) => res.json())
            .then((data) => {
                if (data.error) {
                    console.log(data)
                } else {
                    this.response = data
                }
            })
        }
    },
    setup() {
        const response = ref("")
        const config = ref("")
        Rules.hasToken()
        const cleanedRes = computed (() => {
            return JSON.parse(response.value.reply)
            // return JSON.parse(JSON.parse(response.value.reply))
        })

        function msToTime(ms) {
            let seconds = (ms / 1000).toFixed(1);
            let minutes = (ms / (1000 * 60)).toFixed(1);
            let hours = (ms / (1000 * 60 * 60)).toFixed(1);
            let days = (ms / (1000 * 60 * 60 * 24)).toFixed(1);
            if (seconds < 60) return seconds + " Sec";
            else if (minutes < 60) return minutes + " Min";
            else if (hours < 24) return hours + " Hrs";
            else return days + " Days"
        }

        const runtime = computed(() => {
            let xmit = new Date(response.value.time_tx)
            let rcv = new Date(response.value.time_rx)
            return msToTime((rcv - xmit))
        })
            return { store, response, config, cleanedRes, runtime}
    }
}
</script>

<style scoped>

h5 {
    color: lightslategray;
}

.response {
    color: lightcyan;
    white-space: pre-wrap;
    width: 85%;
    font-family: monospace;
}

</style>