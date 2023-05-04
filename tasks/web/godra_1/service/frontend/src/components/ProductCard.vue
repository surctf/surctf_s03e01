<template>
    <div class="card text-center" style="width: 18rem;">
        <img class="card-img-top" v-bind:src="product.image" alt="Card image cap">
        <div class="card-body">
            <h5 class="card-title">{{ product.name }}</h5>
            <p class="card-text">{{ product.description }}</p>

            <button class="btn btn-primary" v-on:click="buyProduct(product.id)" v-if="!isPurchased">
                Купить за {{ product.cost }}
            </button>
            <button class="btn btn-primary" v-on:click="showProduct(product.id)" v-else>
                Посмотреть
            </button>
        </div>
    </div>
</template>

<script>
import axios from "axios";

export default {
    name: "ProductCard",
    props: ["product", "tgInitData", "isPurchased"],
    emits: ["updateProducts"],
    methods: {
        showProduct(productId) {
            this.$router.push({
                name: "purchase", params: {
                    id: productId,
                }
            })
        },
        buyProduct(productId) {
            this.$swal({
                text: `Купить "${this.product.name}" за ${this.product.cost}?`,
                showDenyButton: true,
                denyButtonText: "Нет",
                button: {
                    text: "Купить!",
                    closeModal: false
                },
                onAfterClose() {
                    this.$swal.showLoading()
                }
            }).then((result) => {
                if (!result.isConfirmed) {
                    return Promise.reject("not confirmed")
                }
                return axios.post("/api/products/" + this.product.id + "/buy", null,
                    {
                        params: {
                            initData: this.tgInitData
                        }
                    })
            }).then(resp => {
                this.$emit("updateProducts")
                return this.$swal.fire({
                    title: "Успешно! Посмотреть содержимое товара?",
                    showDenyButton: true,
                    confirmButtonText: 'Да',
                    denyButtonText: "Нет"
                })
            }).then((result) => {
                if (result.isConfirmed) {
                    this.$router.push({
                        name: "purchase", params: {
                            id: this.product.id
                        }
                    })
                }
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
img {
    max-height: 128px;
    width: auto;
}
</style>