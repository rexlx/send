<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-3">configuration details</h1>
                <hr>
                <form-tag @editConfig="submitHandler" name="userform" event="editConfig">
                <text-input v-model="store.config.name"
                            type="text"
                            required="true"
                            label="name"
                            :value="store.config.name"
                            name="name"
                            ></text-input>
                <text-input v-model="store.config.user"
                            type="text"
                            label="username"
                            required="true"
                            :value="store.config.user"
                            name="username"
                            ></text-input>
                <text-input v-model="store.config.key"
                            type="text"
                            label="key"
                            :value="store.config.key"
                            name="keyfile"
                            ></text-input>
                <text-input v-model="store.config.port"
                            type="number"
                            required="true"
                            label="port"
                            :value="store.config.port"
                            name="port"
                            ></text-input>
                <text-input v-model="store.config.timeout"
                            type="number"
                            required="true"
                            label="timeout"
                            :value="store.config.timeout"
                            name="timeout"
                            ></text-input>
                <text-input v-model="store.config.hosts"
                            type="text"
                            label="hosts"
                            :value="store.config.hosts"
                            name="hosts"
                            ></text-input>
                <hr>
                <div class="float-start">
                    <input type="submit" class="btn btn-primary me-2" value="save">
                    <router-link to="/admin/configs" class="btn btn-outline-secondary">cancel</router-link>
                </div>
                <div class="clearfix"></div>
                </form-tag>
            </div>
        </div>
    </div>
</template>

<script>
import Rules from "./rules.js"
import FormTag from './forms/FormTag.vue'
import TextInput from './forms/TextInput.vue'
import notie from 'notie'
import {store} from "./store"
import router from "../router/index.js"
// import { ref } from 'vue'
import { useRoute } from 'vue-router'


export default {
    setup() {
        const route = useRoute()
        Rules.hasToken()
        fetch(process.env.VUE_APP_API_URL + "/admin/configs/get/" + route.params.configId, Rules.requestOptions(""))
            .then((res) => res.json())
            .then((data) => {
                if (data.error) {
                    this.$emit('error', data.message);
                } else {
                    store.config.value = data;
                }
            })

        return { store }
    },
    components: {
        'form-tag': FormTag,
        'text-input': TextInput,
    },
    methods: {
        submitHandler() {
            const data = {
                id: parseInt(String(this.$route.params.configId), 10),
                user: store.config.user,
                key: store.config.key,
                logpath: "disabled",
                hosts: store.config.hosts,
                command: store.config.command,
                timeout: store.config.timeout,
                port: store.config.port,
                fatal: false,
                ordered: false,
                name: store.config.name
            }
            fetch(`${process.env.VUE_APP_API_URL}/admin/configs/save`, Rules.requestOptions(data))
            .then((res) => res.json())
            .then((data) => {
                if (data.error) {
                    this.$emit('error', data.message);
                } else {
                    this.$emit('success', "changes saved");
                    router.push("/admin/configs")
                }
            })
            .catch((error) => {
                this.$emit('error', error);
            })
        },
        confirmDelete(id) {
            notie.confirm({
                text: "are you sure you want to delete this config?",
                submitText: "delete",
                submitCallback: function() {
                    console.log("goodnight...", id)

                    let payload = {
                        id: id,
                    }

                    fetch(process.env.VUE_APP_API_URL + "/admin/configs/delete", Rules.requestOptions(payload))
                    .then((response) => response.json())
                    .then((data) => {
                        if (data.error) {
                            this.$emit('error', data.message);
                        } else {
                            this.$emit('success', "config deleted");
                            router.push("/admin/configs")
                        }
                    })
                }
            })
        }
    }
}
</script>