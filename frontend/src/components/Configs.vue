<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-3">configurations</h1>
            </div>
            <hr>
            <table v-if="this.ready" class="table table-dark table-striped">
            <thead>
                <tr>
                    <th>name</th>
                    <th>command</th>
                    <th>user</th>
                    <th>port</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="c in this.configs" :key="c.id">
                <td>
                    <router-link :to="`/admin/configs/${c.id}`">{{c.name}}</router-link>
                </td>
                <td>{{c.command}}</td>
                <td>{{c.user}}</td>
                <td>{{c.port}}</td>
                <!-- <td v-if="c.active === 1">
                <span class="badge bg-primary">can run</span>
                </td>
                <td v-else>
                <span class="badge bg-warning">cant run</span>
                </td>
                <td>
                    <span class="badge bg-primary">fix-me</span>
                </td> -->
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
            configs: [],
            ready: false,
            store,
        }
    },
    beforeMount() {
        Rules.hasToken();
        fetch(process.env.VUE_APP_API_URL + "/admin/configs", Rules.requestOptions(""))
        .then((res) => res.json())
        .then((res) => {
            if (res.error) {
                this.$emit('error', res.message);
            } else {
                this.configs = res.data.configs;
                this.ready = true;
            }
        })
        .catch((error) => {
            notie.alert({
                    type: "error",
                    text: "error 292 " + error,
                })
        });
    },
    methods: {
        disableConfig(id) {
            if (id !== store.configs.id) {
                notie.confirm({
                    text: "you sure..?",
                    submitText: "d e a c t i v a t e",
                    submitCallback: () => {
                        fetch(process.env.VUE_APP_API_URL + "/admin/zconfig/" + id, Rules.requestOptions(""))
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
