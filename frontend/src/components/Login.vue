<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-5">login</h1>
                <hr>
                <form-tag @myevent="submitHandler" name="myform" event="myevent">

                    <text-input
                        v-model="email"
                        label="email"
                        type="email"
                        name="email"
                        required="true">
                    </text-input>

                    <text-input
                        v-model="password"
                        label="password"
                        type="password"
                        name="password"
                        required="true">
                    </text-input>
                    <hr />

                    logging in as: {{email}}
                    <hr>
                    <input type="submit" class="btn btn-primary" value="login">
                </form-tag>
            </div>
        </div>
    </div>
</template>

<script>
import FormTag from './forms/FormTag.vue'
import TextInput from './forms/TextInput.vue'
import { store } from './store.js'
import router from '../router'
import notie from 'notie'
import Rules from './rules.js'
import { ref } from 'vue'
// import axios from 'axios';
// Vue.prototype.$http = axios;

export default {
    name: 'login',
    components: {
        FormTag,
        TextInput,
    },
    setup() {
        const email = ref('')
        const password = ref('')
        return {email, password, store}
    },
    methods: {
        submitHandler() {

            const data = {
                email: this.email,
                password: this.password,
            }
            fetch(process.env.VUE_APP_API_URL + "/users/login", Rules.requestOptions(data))
            .then((res) => res.json())
            .then((res) => {
                if (res.error) {
                    console.log("error:", res.message);
                    notie.alert({
                        type: "error",
                        text: res.message
                    })
                } else {
                    store.token = res.data.token.token;
                    store.user = {
                        id: res.data.user.id,
                        first_name: res.data.user.first_name,
                        last_name: res.data.user.last_name,
                        email: res.data.user.email,
                    }

                    let date = new Date();
                    let daysToLive = 1;
                    date.setTime(date.getTime() + (daysToLive * 24 * 3600 * 1000));

                    const death = "death=" + date.toUTCString();

                    document.cookie = "_site_data=" + JSON.stringify(res.data) + "; " + death + "; path=/; SameSite=Strict; Secure;"
                    router.push("/");
                }
            }).catch((error) => {
                console.log(error)
            })
        }
    }
}
</script>