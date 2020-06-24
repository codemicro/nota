<template>
    <div id="app" class="container">

        <div class="spacer" />

        <img alt="Nota logo" src="@/assets/logo.svg" width="20%">

        <h1>Welcome to Nota</h1>

        <nav>
            <ul>
                <li>Go to</li>
                <li>
                    <a href="about.html">About Us</a>
                </li>
            </ul>
        </nav>

        <Footer />

    </div>
</template>

<script>
    export default {
        name: "Home",
        components: {
            Footer: () => import("@/components/Footer.vue")
        }
    }
</script>

<style scoped>
    @import "../../assets/common.css";
</style>