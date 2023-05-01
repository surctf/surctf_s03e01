import {createRouter, createWebHashHistory} from 'vue-router'
import ProfileView from '../views/ProfileView.vue'
import MarketView from '../views/MarketView.vue'
import SignUpView from '../views/SignUpView.vue'
import PurchaseView from "@/views/PurchaseView.vue";
import SettingsView from "@/views/SettingsView.vue";
import ColabView from "@/views/ColabView.vue";

const routes = [
    {
        path: '/profile',
        name: 'profile',
        props: true,
        component: ProfileView
    },
    {
        path: '/market',
        name: 'market',
        props: true,
        component: MarketView
    },
    {
        path: '/purchase/:id',
        name: 'purchase',
        props: true,
        component: PurchaseView
    },
    {
        path: '/signup',
        name: 'signup',
        props: true,
        component: SignUpView
    },
    {
        path: '/settings',
        name: 'settings',
        props: true,
        component: SettingsView
    },
    {
        path: '/colab',
        name: 'colab',
        props: true,
        component: ColabView
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router
