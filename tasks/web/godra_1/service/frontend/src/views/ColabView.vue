<template>
    <div class="d-flex flex-column">
        <div class="p-2 mx-auto">
            <h1> godra x surctf </h1>
            <p> В честь колабы godra и surctf устраиваем конкурс!
                Сделай 20 любых покупок на маркетплейсе и получи флаг!</p>

            <p> Твой флаг: {{ flag }} </p>
        </div>
    </div>
</template>

<script>
import axios from "axios";

export default {
    name: "ColabView.vue",
    props: ["tgInitData", "user"],
    emits: ["updateUser"],
    data() {
        return {
            flag: "loading..."
        }
    },

    methods: {
        getFlag() {
            axios.get("/api/colab", {
                params: {initData: this.tgInitData}
            }).then(resp => {
                this.flag = resp.data.flag
            }).catch(err => {
                if (err.response && err.response.data.error) {
                    this.flag = err.response.data.error
                } else {
                    this.flag = err
                }
            })
        }
    },
    mounted() {
        this.getFlag()
    }
}
</script>

<style scoped>

</style>