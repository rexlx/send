<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-3">target details</h1>
                <hr>
                <form-tag @editTarget="submitHandler" name="userform" event="editTarget">
                <text-input v-model="target.address"
                            type="text"
                            required="true"
                            label="address"
                            :value="target.address"
                            name="address"
                            ></text-input>
                <text-input v-model="target.username"
                            type="text"
                            label="username"
                            required="true"
                            :value="target.username"
                            name="username"
                            ></text-input>
                <text-input v-model="target.password"
                            type="password"
                            label="password"
                            :value="target.password"
                            name="password"
                            ></text-input>
                <text-input v-model="target.port"
                            type="text"
                            required="true"
                            label="port"
                            :value="target.port"
                            name="port"
                            ></text-input>
                <text-input v-model="target.protocol"
                            type="text"
                            required="true"
                            label="protocol"
                            :value="target.protocol"
                            name="protocol"
                            ></text-input>
                <text-input v-model="target.token"
                            type="password"
                            label="token"
                            :value="target.token"
                            name="token"
                            ></text-input>
                
                <div class="form-check">
                    <input v-model="target.active" class="form-check-input" type="radio" id="target-active" :value="1">
                    <label for="target-active" class="form-check-label">listen</label>
                </div>
                <div class="form-check">
                    <input v-model="target.active" class="form-check-input" type="radio" id="target-active-2" :value="0">
                    <label for="target-active-2" class="form-check-label">sleep</label>
                </div>

                <hr>
                <div class="float-start">
                    <input type="submit" class="btn btn-primary me-2" value="save">
                    <router-link to="/admin/targets" class="btn btn-outline-secondary">cancel</router-link>
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


export default {
    beforeMount() {
        Rules.hasToken();

        if (parseInt(String(this.$route.params.targetId), 10) > 0) {
            //edit existing user
            fetch(process.env.VUE_APP_API_URL + "/admin/targets/get/" + this.$route.params.targetId, Rules.requestOptions(""))
            .then((res) => res.json())
            .then((data) => {
                if (data.error) {
                    this.$emit('YOOOOO', data.message);
                } else {
                    this.target = data;
                    this.target.password = "";
                }
            })
        }
    },
    data() {
        return {
            target: {
                id: 0,
                address: "",
                protocol: 0,
                username: "",
                port: 0,
                password: "",
                active: 0,
                token: "",
            },
            store,
        }
    },
    components: {
        'form-tag': FormTag,
        'text-input': TextInput,
    },
    methods: {
        submitHandler() {
            const data = {
            id: parseInt(String(this.$route.params.targetId), 10),
            address: this.target.address,
            username: this.target.username,
            port: parseInt(String(this.target.port)),
            protocol: parseInt(String(this.target.protocol)),
            password: this.target.password,
            active: parseInt(String(this.target.active)),
            token: this.target.token,
            }
            fetch(`${process.env.VUE_APP_API_URL}/admin/targets/save`, Rules.requestOptions(data))
            .then((res) => res.json())
            .then((data) => {
                if (data.error) {
                    this.$emit('error', data.message);
                } else {
                    this.$emit('success', "changes saved");
                    router.push("/admin/targets")
                }
            })
            .catch((error) => {
                this.$emit('error', error);
            })
        },
        confirmDelete(id) {
            notie.confirm({
                text: "are you sure you want to delete this user?",
                submitText: "delete",
                submitCallback: function() {
                    console.log("goodnight...", id)

                    let payload = {
                        id: id,
                    }

                    fetch(process.env.VUE_APP_API_URL + "/admin/targets/delete", Rules.requestOptions(payload))
                    .then((response) => response.json())
                    .then((data) => {
                        if (data.error) {
                            this.$emit('error', data.message);
                        } else {
                            this.$emit('success', "target deleted");
                            router.push("/admin/targets")
                        }
                    })
                }
            })
        }
    }
}
</script>