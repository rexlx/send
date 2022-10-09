<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-3">all accounts</h1>
            </div>
            <hr>
            <table v-if="this.ready" class="table table-dark table-striped">
            <thead>
                <tr>
                    <th>user</th>
                    <th>email</th>
                    <th>status</th>
                    <th>session</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="u in this.users" v-bind:key="u.id">
                <td>
                    <router-link :to="`/admin/users/${u.id}`">{{u.last_name}}, {{u.first_name}}</router-link>
                </td>
                <td>{{u.email}}</td>
                <td v-if="u.active === 1">
                <span class="badge bg-info">active</span>
                </td>
                <td v-else>
                <span class="badge bg-warning">inactive</span>
                </td>
                <td v-if="u.token.id > 0">
                <a href="javascripts:void(0);">
                    <span class="badge bg-info" @click="bootUser(u.id)">logged in</span>
                </a>
                </td>
                <td v-else>
                <span class="badge bg-warning">logged out</span>
                </td>
                </tr>
            </tbody>
            </table>
            <p v-else>loading...</p>
        </div>
    </div>
</template>

<script>
import Rules from "./rules.js"
import notie from 'notie'
import {store} from './store.js'

export default {
    data() {
        return {
            users: [],
            ready: false,
            store,
        }
    },
    beforeMount() {
        Rules.hasToken();
        fetch(process.env.VUE_APP_API_URL + "/admin/users", Rules.requestOptions(""))
        .then((res) => res.json())
        .then((res) => {
            if (res.error) {
                this.$emit('error', res.message);
            } else {
                this.users = res.data.users;
                this.ready = true;
            }
        })
        .catch((error) => {
            notie.alert({
                    type: "error",
                    text: "error 365 " + error,
                })
        });
    },
    methods: {
        bootUser(id) {
            if (id !== store.user.id) {
                notie.confirm({
                    text: "you sure..?",
                    submitText: "boot",
                    submitCallback: () => {
                        fetch(process.env.VUE_APP_API_URL + "/admin/boot/" + id, Rules.requestOptions(""))
                        .then((res) => res.json())
                        .then((data) => {
                            if (data.error) {
                                this.$emit('error', data.message)
                            } else {
                                this.$emit('success', data.message)
                                this.$emit('forceUpdate')
                            }
                        })
                    }
                })
            } else {
                this.$emit('error', "this action is forbidden")
            }
        }
    }
}
</script>

<style>
th {
    color: bisque;
}
</style>