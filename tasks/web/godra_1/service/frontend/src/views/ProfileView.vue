<template>
    <div class="home">
        <div class="d-flex flex-column">
            <h1> Твоё </h1>
            <div class="p-2 mx-auto" v-for="purchase in user.purchases">
                <ProductCard :isPurchased="true" :product="purchase.product"/>
            </div>
        </div>
    </div>
</template>

<script>

import ProductCard from "@/components/ProductCard.vue";

export default {
    components: {ProductCard},
    emits: ["updateUser"],
    props: ["tgInitData", "user"],
    mounted() {
            this.$emit("updateUser")
    }
}
</script>
