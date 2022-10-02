import {store} from './store.js'
import router from '../router/index.js'

let Rules = {
    hasToken: function () {
        if (store.token === '') {
            router.push("/login");
            return false
        }
    },
    requestOptions: function(data) {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");
        headers.append("Authorization", "Bearer " + store.token);
        return {
            method: "POST",
            body: JSON.stringify(data),
            headers: headers,
        }
    },
    checkToken: function() {
        if (store.token !== "") {
            const data = {
                token: store.token,
            }
            const headers = new Headers();
            headers.append("Content-Type", "application/json");

            let requestOptions = {
                method: "POST",
                body: JSON.stringify(data),
                headers: headers,
            }
            fetch(process.env.VUE_APP_API_URL + '/vtk', requestOptions)
            .then((res) => res.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.error)
                } else {
                    if (!data.data) {
                        store.token = "";
                        store.user = {},
                        document.cookie = '_site_data=; Path=/; SameSite=strict; Secure; Expires=Thu, 01 Jan 19070 00:00:01 GMT;'
                    }
                }
            })
        }
    }
}

export default Rules;