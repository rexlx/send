<template>
    <div v-if="store.token !== ''" class="container">
        <div class="row">
            <div class="col-md">
                <p class="row text-center">welcome {{ store.user.email }}. click the admin dropdown to reveal more options</p>
            </div>
        </div>
    </div>
    <div v-else class="container">
        <div class="row">
            <h3>logged out</h3>
            <hr>
        </div>
        <div class="row">
            <div class="col">
                <p>please login to continue.</p>
            </div>
        </div>
        <div class="row">
            <div class="col">
                <p>access is <em>restricted</em>...</p>
            </div>
        </div>
    </div>
</template>

<script>
import {store} from "./store"

export default {
    data() {
        console.log(process.env.VUE_APP_IMG_URL)
    return {
      store,
    //   bucket: process.env.VUE_APP_IMG_URL,
    }
  }
}
</script>

<style>
body {
    background: black;
    color: bisque;
}
</style>
