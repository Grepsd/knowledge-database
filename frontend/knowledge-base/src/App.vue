<template>
    <v-app id="inspire">
        <v-navigation-drawer
                v-model="drawer"
                app
                clipped
        >
            <v-list dense>
                <v-list-item link>
                    <v-list-item-action>
                        <v-icon>mdi-view-dashboard</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Dashboard</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
                <v-list-item link>
                    <v-list-item-action>
                        <v-icon>mdi-cog</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Settings</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
            </v-list>
        </v-navigation-drawer>

        <v-main>
            <v-container class="fill-height" fluid style="align-items: start">
                <v-row>
                    <v-col>
                        <articles-list>
                        </articles-list>
                    </v-col>
                </v-row>
            </v-container>
        </v-main>

        <v-footer app>
            <span>&copy; {{ new Date().getFullYear() }}</span>
        </v-footer>
    </v-app>
</template>

<script>
    import ArticlesList from "./components/ArticlesList";

    export default {
        components: {
            ArticlesList
        },
        props: {
            source: String,
        },
        data: () => ({
            drawer: null,
        }),
        created() {
            // this.$vuetify.theme.dark = true
        },
    }
</script>

<style type="text/css">

</style>