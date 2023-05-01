<template>
    <div class="d-flex flex-column">
        <div class="p-2 mx-auto">
            <button class="btn btn-danger" v-on:click="this.delete()"> Удалить мой аккаунт</button>
        </div>
    </div>
</template>

<script>

import axios from "axios";

export default {
    name: "SettingsView.vue",
    emits: ["updateUser"],
    props: ["tgInitData", "user"],
    methods: {
        delete() {
            this.$swal({
                text: `Уверен что хочешь удалить свой аккаунт?`,
                showDenyButton: true,
                denyButtonText: "Нет",
                button: {
                    text: "Да",
                    closeModal: false
                },
                onAfterClose() {
                    this.$swal.showLoading()
                }
            }).then((result) => {
                if (!result.isConfirmed) {
                    return Promise.reject("not confirmed")
                }
                return axios.delete("/api/user", {
                    params:
                        {initData: this.tgInitData}
                })
            }).then(resp => {
                this.$emit("updateUser")
                return this.$swal.fire({
                    title: "Успешно!",
                })
            }).catch((err) => {
                if (err === "not confirmed") {
                    return
                }

                let text
                if (err.response && err.response.data.error) {
                    text = err.response.data.error
                } else {
                    text = err
                }

                this.$swal.fire({
                    title: "Ошибка!",
                    text: text,
                })
            })


        }
    }
}
</script>

<style scoped>

</style>