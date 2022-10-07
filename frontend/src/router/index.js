import { createRouter, createWebHistory } from 'vue-router'
import Body from './../components/Body.vue'
import Login from './../components/Login.vue'
import Users from "./../components/Users.vue"
import User from "./../components/UserEdit.vue"
import Targets from "./../components/Targets.vue"
import Target from "./../components/TargetEdit.vue"
import Console from "../components/Console.vue"
import Receiver from "../components/Receiver.vue"
import Result from "../components/Result.vue"
import Configs from "../components/Configs.vue"
import Config from "../components/ConfigEdit.vue"
import Rules from '../components/rules'

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Body,
    },
    {
        path: '/login',
        name: 'Login',
        component: Login,
    },
    {
        path: '/admin/users',
        name: 'Users',
        component: Users,
    },
    {
        path: '/admin/users/:userId',
        name: 'User',
        component: User,
    },
    {
        path: '/admin/targets',
        name: 'Targets',
        component: Targets,
    },
    {
        path: '/admin/targets/:targetId',
        name: 'Target',
        component: Target,
    },
    {
        path: '/admin/responses/:responseId',
        name: 'Response',
        component: Result,
    },
    {
        path: '/admin/configs',
        name: 'Configs',
        component: Configs,
    },
    {
        path: '/admin/configs/:configId',
        name: 'Config',
        component: Config,
    },
    {
        path: '/admin/console',
        name: 'Console',
        component: Console,
    },
    {
        path: '/admin/config',
        name: 'Receiver',
        component: Receiver,
    },
]

const router = createRouter({history: createWebHistory(), routes})
router.beforeEach(() => {
    Rules.checkToken();
})
export default router