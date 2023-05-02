<template>
    <template v-if="error">
        <p style="color: red"> {{ error }}</p>
    </template>
    <template v-else>
        <div class="d-flex flex-column" v-if="products && products.length">
            <div class="p-2 mx-auto" v-for="product in products">
                <ProductCard :tgInitData="tgInitData" :product="product" v-if="!product.purchased"/>
            </div>
        </div>
        <div v-else> Пока нет новых товаров </div>
    </template>
</template>

<script>
import axios from "axios";
import ProductCard from "@/components/ProductCard.vue";

export default {
    name: "MarketView",
    components: {ProductCard},
    props: ["tgInitData", "user"],
    emits: ["updateUser"],
    data() {
        return {
            products: undefined,
            error: "",
        }
    },
    methods: {
        getProducts() {
            axios.get("/api/products", {
                params: {
                    initData: this.tgInitData
                }
            }).then(resp => {
                this.products = resp.data.products
                console.log(resp.data)
            }).catch(err => {
                if (err.response.data) {
                    this.error = err.response.data.error
                } else {
                    this.error = err
                }
            })
        }
    },
    mounted() {
        this.getProducts()
    }
}
</script>

<style scoped>

</style>