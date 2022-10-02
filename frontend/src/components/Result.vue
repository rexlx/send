<template>
  <div class="container-fluid">
    <div v-if="response">
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
            return { store, response, config, cleanedRes}
    }
}
</script>

<style scoped>

.host {
    overflow-y: scroll;
    color: rgb(117, 106, 90);
}

.response {
    color: rgb(160, 205, 205);
    white-space: pre-wrap;
    width: 85%;
    font-family: monospace;
}

</style>