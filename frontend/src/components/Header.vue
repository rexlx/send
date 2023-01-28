<template>
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
  <div class="container-fluid">
    <a class="navbar-brand" href="#">site</a>
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarNav">
      <ul class="navbar-nav me-auto mb-2 mb-lg-0">
        <li class="nav-item">
          <router-link class="nav-link active" aria-current="page" to="/">home</router-link>
        </li>

        <li v-if="store.token !== ''" class="nav-item dropdown">
        <a class="nav-link dropdown-toggle" href="#" id="navdrop" role="button" data-bs-toggle="dropdown"
        aria-expanded="false">
        admin
        </a>
        <ul class="dropdown-menu" aria-labelledby="navdrop">
          <li>
          <router-link class="dropdown-item" to="/admin/console">console</router-link>
        </li>
          <li>
            <router-link class="dropdown-item" to="/admin/users">users</router-link>
          </li>
          <li>
            <router-link class="dropdown-item" to="/admin/users/0">add user</router-link>
          </li>
        <!-- <li>
          <router-link class="dropdown-item" to="/admin/targets">targets</router-link>
        </li> -->
        <!-- <li>
          <router-link class="dropdown-item" to="/admin/targets/0">add target</router-link>
        </li> -->
        <li>
          <router-link class="dropdown-item" to="/admin/configs">configs</router-link>
        </li>
        <li>
          <router-link class="dropdown-item" to="/admin/configs/0">add config</router-link>
        </li>
        </ul>
        </li>
        <li class="nav-item">
          <router-link v-if="store.token == ''" class="nav-link" to="/login">login</router-link>
          <a href="javascript:void(0);" v-else class="nav-link" @click="logout">logout</a>
        </li>
      </ul>

      <span class="navbar-text">
        {{  store.user.first_name ?? 'not logged in...'}}
      </span>
    </div>
  </div>
</nav>
</template>

<script>
import { store } from './store.js'
import router from '../router/index.js'
import Rules from './rules.js'

export default {
  data() {
    return {
      store
    }
  },
  methods: {
    logout() {
      const data = {
        token: store.token,
        user: store.user.email
      }
      fetch(process.env.VUE_APP_API_URL +"/users/logout", Rules.requestOptions(data))
      .then((res) => res.json())
      .then((res) => {
        if (res.error) {
          console.log(res.message)
        } else {
          store.token = "";
          store.user = {};
          document.cookie = '_site_data=; Path=/; SameSite=Strict; secure; Expires=Thu, 01 Jan 1970 00:00:01 GMT;'
          router.push("/login")
        }
      })
      
    }
  }
}
</script>