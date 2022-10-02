<template>
    <div v-if="store.token !== ''" class="container">
        <div class="row">
            <div class="col">
                <img alt="bookNshib" :src="`${bucket}/stock9.jpg`" class="img-fluid mt-5">
            </div>
        </div>
        <div class="row">
            <div class="col">
                <p>welcome {{ store.user.email }}</p>
            </div>
        </div>
    </div>
    <div v-else class="container">
        <div class="row">
            <h3>go away</h3>
            <hr>
        </div>
        <div class="row">
            <div class="col">
                <img alt="bookNshib" :src="`${bucket}/stock0.jpeg`" class="img-fluid mt-5">
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
      bucket: "https://storage.googleapis.com/rfitzhugh"
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
