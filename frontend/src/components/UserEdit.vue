<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-3">account details</h1>
                <hr>
                <form-tag @editUser="submitHandler" name="userform" event="editUser">
                <text-input v-model="user.first_name"
                            type="text"
                            required="true"
                            label="name"
                            :value="user.first_name"
                            name="first_name"
                            ></text-input>
                <text-input v-model="user.last_name"
                            type="text" required="true"
                            label="last name"
                            :value="user.last_name"
                            name="last_name"
                            ></text-input>
                <text-input v-model="user.email"
                            type="email"
                            required="true"
                            label="email"
                            :value="user.email"
                            name="email"
                            ></text-input>
                <text-input v-if="this.user.id === 0"
                            v-model="user.password"
                            type="password"
                            required="true"
                            label="password"
                            :value="user.password"
                            name="password"
                            ></text-input>
                <text-input v-else v-model="user.password"
                            type="password"
                            label="password"
                            help="changes here will effect the password"
                            :value="user.password"
                            name="password"
                            ></text-input>
                
                <div class="form-check">
                    <input v-model="user.active" class="form-check-input" type="radio" id="user-active" :value="1">
                    <label for="user-active" class="form-check-label">active</label>
                </div>
                <div class="form-check">
                    <input v-model="user.active" class="form-check-input" type="radio" id="user-active-2" :value="0">
                    <label for="user-active-2" class="form-check-label">inactive</label>
                </div>

                <hr>
                <div class="float-start">
                    <input type="submit" class="btn btn-primary me-2" value="save">
                    <router-link to="/admin/users" class="btn btn-outline-secondary">cancel</router-link>
                </div>
                <div class="float-end">
                    <a v-if="(this.$route.params.userId > 0) && (parseInt(String(this.$route.params.userId), 10) !== store.user.id)"
                            class="btn btn-danger" href="javascript:void(0);" @click="confirmDelete(this.user.id)">delete</a>
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

        if (parseInt(String(this.$route.params.userId), 10) > 0) {
            //edit existing user
            fetch(process.env.VUE_APP_API_URL + "/admin/users/get/" + this.$route.params.userId, Rules.requestOptions(""))
            .then((res) => res.json())
            .then((data) => {
                if (data.error) {
                    this.$emit('error', data.message);
                } else {
                    this.user = data;
                    this.user.password = "";
                }
            })
        }
    },
    data() {
        return {
            user: {
                id: 0,
                first_name: "",
                last_name: "",
                email: "",
                password: "",
                active: 0,
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
            id: parseInt(String(this.$route.params.userId), 10),
            first_name: this.user.first_name,
            last_name: this.user.last_name,
            email: this.user.email,
            password: this.user.password,
            active: this.user.active,
            }
            fetch(`${process.env.VUE_APP_API_URL}/admin/users/save`, Rules.requestOptions(data))
            .then((res) => res.json())
            .then((data) => {
                if (data.error) {
                    this.$emit('error', data.message);
                } else {
                    this.$emit('success', "changes saved");
                    router.push("/admin/users")
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

                    fetch(process.env.VUE_APP_API_URL + "/admin/users/delete", Rules.requestOptions(payload))
                    .then((response) => response.json())
                    .then((data) => {
                        if (data.error) {
                            this.$emit('error', data.message);
                        } else {
                            this.$emit('success', "user deleted");
                            router.push("/admin/users")
                        }
                    })
                }
            })
        }
    }
}
</script>