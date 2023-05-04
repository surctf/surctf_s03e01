<template>
    <template v-if="isLoading">
        Loading...
    </template>
    <template v-else>
        <template v-if="tgInitData">
            <template v-if="isRegistered">
                <nav class="navbar fixed-top navbar-light bg-light">
                    <a class="nav-brand">godra</a>
                    <router-link class="nav-item nav-link" to="/colab">godra x surctf</router-link>
                    <router-link class="nav-item nav-link" to="/colab">{{ user.user.username }} (Баланс: {{user.user.balance}})</router-link>
                </nav>

                <nav class="navbar fixed-bottom navbar-light bg-light">
                    <router-link class="nav-item nav-link" to="/market">Купить</router-link>
                    <router-link class="nav-item nav-link" to="/profile">Моё</router-link>
                    <router-link class="nav-item nav-link" to="/settings">Настройки</router-link>
                </nav>
            </template>
            <router-view class="text-center" style="margin-top: 60px; margin-bottom: 60px" :tgInitData="tgInitData" :user="user" @updateUser="updateUser" />
        </template>
        <p v-else>
            Пж заходи через тг
        </p>
    </template>

</template>

<script>
import axios from 'axios';

export default {
    data() {
        return {
            tgInitData: undefined,
            user: undefined,
            isLoading: true,
            isRegistered: false
        }
    },
    mounted() {
        this.tgInitData = window.Telegram.WebApp.initData
    },

    watch: {
        tgInitData(oldTgInitData, newTgInitData) {
            if (newTgInitData !== "") {
                this.getUser()
            }
        }
    },
    methods: {
        getUser() {
            axios
                .get("/api/user", {params: {initData: this.tgInitData}})
                .then((resp) => {
                    this.isRegistered = resp.status === 200;
                    this.isLoading = false
                    this.user = resp.data
                    this.$router.push("/profile")
                })
                .catch(err => {
                        this.$router.push("/signup")
                        this.isLoading = false
                        this.isRegistered = false
                        console.log(err)
                    }
                )
        },
        updateUser() {
            axios
                .get("/api/user", {params: {initData: this.tgInitData}})
                .then((resp) => {
                    this.isRegistered = resp.status === 200;
                    this.isLoading = false
                    this.user = resp.data
                })
                .catch(err => {
                        this.$router.push("/signup")
                        this.user = undefined
                        this.isLoading = false
                        this.isRegistered = false
                        console.log(err)
                    }
                )
        }
    }
}
</script>

<style>
.navbar > a{
    text-decoration-line: none;
    padding-left:10px;
    padding-right:10px;
    margin-left: 10px;
    margin-right: 10px;
    border: 1px solid black;
}
</style>
