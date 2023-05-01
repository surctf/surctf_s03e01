<template>
    <div class="d-flex flex-column">
        <div class="p-2 mx-auto">

            <template v-if="isLoading">
                isLoading
            </template>
            <template v-else>
                <h1> Регистрация </h1>
                <p id="error" v-if="error"> {{ error }} </p>
                <form>
                    <label> Юзернейм </label><br>
                    <input v-model="username" placeholder="юзернейм сюда впиши"><br>
                    <button class="btn btn-primary" :disabled="isDisabled" v-on:click="submit">Готово</button>
                </form>
            </template>
        </div>
    </div>
</template>

<script>
import axios from "axios";

export default {
    name: "SignUpView.vue",
    props: ["tgInitData", "user"],
    emits: ["updateUser"],
    data() {
        return {
            username: "",
            error: undefined,
            isDisabled: false,
            isLoading: false
        }
    },
    watch: {
        username(oldUsername, newUsername) {
            if (newUsername.length > 24) {
                this.error = "Слишком длинный юзернейм"
                this.isDisabled = true
            } else {
                this.error = undefined
                this.isDisabled = false
            }
        }
    },
    methods: {
        submit: function () {
            this.isLoading = true
            axios.post("/api/user", {
                username: this.username
            }, {
                params: {
                    initData: this.tgInitData
                }
            }).then(resp => {
                this.$emit("updateUser")
                this.$router.push("/")
                this.isLoading = false
            }).catch(err => {
                if (err.response.data) {
                    this.error = err.response.data.error
                } else {
                    this.error = err
                }
                this.isLoading = false
            })
        }
    }
}
</script>

<style scoped>
#error {
    color: red;
}
</style>