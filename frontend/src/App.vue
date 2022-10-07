<template>
<div>
  <Header />
  <div>
    <router-view :key="componentKey" @success="success" @error="error" @warning="warning" @forceUpdate="forceUpdate" />
  </div>
  <Footer />
</div>
</template>

<script>
import Header from "./components/Header.vue"
import Footer from "./components/Footer.vue"
import { store } from './components/store.js'
import notie from 'notie'

const getCookie = (name) => {
  return document.cookie.split("; ").reduce((r, v) => {
    const values = v.split("=");
    return values[0] === name ? decodeURIComponent(values[1]) : r;
  }, "");
}

export default {
  name: 'App',
  components: {
    Header,
    Footer,
  },
  data() {
    return {
      store,
      componentKey: 0,
    }
  },
  beforeMount() {
    let data = getCookie("_site_data");
    if (data !== "") {
      let cookieData = JSON.parse(data);

      store.token = cookieData.token.token
      store.user = {
        id: cookieData.user.id,
        first_name: cookieData.user.first_name,
        last_name: cookieData.user.last_name,
        email: cookieData.user.email,
      }
    }
  },
  methods: {
    success(msg) {
      notie.alert({
        type: 'success',
        text: msg,
      })
    },
    error(msg) {
      notie.alert({
        type: 'error',
        text: msg,
      })
    },
    warning(msg) {
      notie.alert({
        type: 'warning',
        text: msg,
      })
    },
    forceUpdate() {
      this.componentKey += 1;
    }
  }
}
</script>

<style>

</style>