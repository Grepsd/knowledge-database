<template>
    <v-container fluid>
        <v-chip-group :column="true">
            <v-layout align-center justify-center>
                <v-chip v-for="tag in tags" :key="tag.id">
                    {{ tag.name }}
                </v-chip>
            </v-layout>
        </v-chip-group>
    </v-container>
</template>

<script>
    export default {
        data: function () {
            return {
                tags: {},
            }
        },
        name: "TagsList",
        methods: {
            loadTags() {
                this.axios.get('/tags?has_articles=true')
                    .then(response => this.tags = response.data)
            }
        },
        mounted() {
            this.loadTags()
        }
    }
</script>

<style scoped>

</style>
