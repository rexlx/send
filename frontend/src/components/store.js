import { reactive } from "vue";

export const store = reactive({
    token: "",
    user: {},
    config: {},
    commandHistory: [],
})