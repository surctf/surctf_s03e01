<template>
    <div v-if="product">
        <div class="d-flex flex-column">
            <div class="p-2 mx-auto">
                <div class="card text-center" style="width: 18rem;">
                    <img class="card-img-top" v-bind:src="product.image" alt="Card image cap">
                    <div class="card-body">
                        <h5 class="card-title">{{ product.name }}</h5>
                        <p class="card-text">{{ product.description }}</p>
                        <p class="card-footer"> {{ product.secret }} </p>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div v-else>
        Not found
    </div>
</template>

<script>
export default {
    name: "PurchaseView.vue",
    emits: ["updateUser"],
    props: ["tgInitData", "user"],
    computed: {
        product() {
            let prodId = this.$route.params.id
            for (const purchase of this.user.purchases.values()) {
                console.log(purchase.product.id, prodId)
                if (purchase.product.id === Number(prodId)) {
                    return purchase.product
                }
            }
            return undefined
        }
    },
    mounted() {
        this.$emit("updateUser")
    }
}
</script>

<style scoped>

</style>