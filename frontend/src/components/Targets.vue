<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-3">all accounts</h1>
            </div>
            <hr>
            <div class="jumbotron jumbotron-fluid">
                <div class="container">
                    <h1 class="display-4 errata">this feature incomplete</h1>
                    <p class="lead">you can add as many <em>targets</em> as you'd like.</p>
                    <p>right now we use a slice of strings for hosts. one day it should be a slice of type targets</p>
                </div>
            </div>
            <table v-if="this.ready" class="table table-dark table-striped">
            <thead>
                <tr>
                    <th>host</th>
                    <th>proto</th>
                    <th>status</th>
                    <th>presence</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="t in this.targets" :key="t.id">
                <td>
                    <router-link :to="`/admin/targets/${t.id}`">{{t.address}}</router-link>
                </td>
                <td>{{t.protocol}}</td>
                <td v-if="t.active === 1">
                <span class="badge bg-primary">listening</span>
                </td>
                <td v-else>
                <span class="badge bg-warning">inactive</span>
                </td>
                <td>
                    <span class="badge bg-primary">fix-me</span>
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
            targets: [],
            ready: false,
            store,
        }
    },
    beforeMount() {
        Rules.hasToken();
        fetch(process.env.VUE_APP_API_URL + "/admin/targets", Rules.requestOptions(""))
        .then((res) => res.json())
        .then((res) => {
            if (res.error) {
                this.$emit('error', res.message);
            } else {
                this.targets = res.data.targets;
                this.ready = true;
            }
        })
        .catch((error) => {
            notie.alert({
                    type: "error",
                    text: "error 765 " + error,
                })
        });
    },
    methods: {
        disableTarget(id) {
            if (id !== store.target.id) {
                notie.confirm({
                    text: "you sure..?",
                    submitText: "d e a c t i v a t e",
                    submitCallback: () => {
                        fetch(process.env.VUE_APP_API_URL + "/admin/ztarget/" + id, Rules.requestOptions(""))
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
h1 {
    color: bisque;
}

.errata {
    color:tomato
}
.errata:hover {
    color: skyblue;
}

p {
    color: azure;
}
</style>