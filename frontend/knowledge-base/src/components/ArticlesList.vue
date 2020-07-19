<template>
    <v-container>
        <v-row>
            <v-col>
                <v-icon large v-on:click="loadArticles" color="black">mdi-refresh</v-icon>
            </v-col>
            <v-col>
                <v-icon large v-on:click="addDialog = true" color="green">mdi-plus</v-icon>
            </v-col>
            <v-col>
                <v-icon large v-on:click="deleteMode = !deleteMode" :color="deleteMode ? 'red' : 'green'">mdi-delete
                </v-icon>
            </v-col>
        </v-row>
        <v-divider/>
        <TagsList></TagsList>
        <v-divider/>
        <ul>
            <v-row>
                <v-col class="column_wrapper">
                    <li v-for="(articles, index) in articlesSortedByFirstLetter" v-bind:key="index"
                        style="list-style: none">
                        <v-row class="font-weight-bold font-weight-black blue--text">
                            <h2>{{ index.toUpperCase()}}</h2>
                        </v-row>
                        <v-row v-for="(article) in articles" v-bind:key="article.id" align="center">
                            <a :href="article.url" target="_blank">{{ article.title }}</a>
                            <v-chip-group>
                                <v-chip v-for="tag in article.tags" :key="tag.id"
                                        class="ma-2"
                                        x-small
                                        color="blue lighten-3"
                                >
                                    {{ tag.name }}
                                </v-chip>
                            </v-chip-group>
                            <v-icon v-if="deleteMode" small class="padding" color="black"
                                    v-on:click="showDeleteAction(article.ID, article.title, article.url)">mdi-delete
                            </v-icon>
                        </v-row>
                    </li>
                </v-col>
            </v-row>
        </ul>
        <v-row justify="center">
            <v-dialog v-model="addDialog" persistent max-width="400">
                <v-card>
                    <v-card-title>Add an article</v-card-title>
                    <v-divider></v-divider>
                    <v-card-text>
                        <create-article v-on:articleCreated="addArticle"
                                        v-on:closeDialog="addDialog = false"></create-article>
                    </v-card-text>
                </v-card>
            </v-dialog>
        </v-row>
        <v-row justify="center">
            <v-dialog v-model="deleteDialog" persistent max-width="350">
                <v-card>
                    <v-card-title class="headline">Confirm article deletion</v-card-title>
                    <v-card-text>Are you certain you want to delete the article :
                        <v-row>
                            <v-col>
                                <v-chip color="light-blue lighten-5">{{ deleteArticle.title }}</v-chip>
                            </v-col>
                        </v-row>
                        <v-row>
                            linking to
                        </v-row>
                        <v-row>
                            <v-col>
                                <v-chip color="light-green lighten-5">
                                    <a :href="deleteArticle.URL" target="_blank">{{ deleteArticle.URL}}</a>
                                </v-chip>
                            </v-col>
                        </v-row>
                    </v-card-text>
                    <v-card-actions>
                        <v-spacer></v-spacer>
                        <v-container>
                            <v-row cols="8">
                                <v-col cols="4">
                                    <v-btn color="red darken-1" text @click="deleteDialog = false">No</v-btn>
                                </v-col>
                                <v-col cols="4">
                                    <v-btn color="green darken-1" text @click="sendDeleteArticleRequest()">Agree</v-btn>
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-card-actions>
                </v-card>
            </v-dialog>
        </v-row>
        <v-snackbar v-model="snackbar.show" :color="snackbar.color">
            {{ snackbar.text}}
        </v-snackbar>
    </v-container>
</template>

<script>
    import CreateArticle from "./CreateArticle";
    import TagsList from "./TagsList";
    import {mdiDelete} from '@mdi/js'

    export default {
        name: "ArticlesList",
        components: {
            CreateArticle,
            TagsList,
        },
        data: function () {
            return {
                articles: [],
                articlesSortedByFirstLetter: [],
                mdiDelete,
                deleteArticle: {title: null, id: null, URL: null},
                deleteDialog: false,
                addDialog: false,
                snackbar: {
                    show: false,
                    color: "white",
                    text: "",
                },
                deleteMode: false,
            }
        },
        methods: {
            loadArticles() {
                this.axios.get("/articles")
                    .then(response => {
                        if (response.status !== 200) {
                            return
                        }
                        this.articlesSortedByFirstLetter = {};
                        let articles = response.data.sort((prev, next) => prev.title.toLowerCase() > next.title.toLowerCase());
                        articles.forEach(item => {
                            if (this.articlesSortedByFirstLetter[item.title[0].toLowerCase()] === undefined) {
                                this.articlesSortedByFirstLetter[item.title[0].toLowerCase()] = []
                            }
                            this.articlesSortedByFirstLetter[item.title[0].toLowerCase()].push(item)
                        });
                        this.$forceUpdate()
                        this.snackbar = {
                            color: "green",
                            text: "Refresh successful",
                            show: true,
                        }
                    })
                    .catch(error =>
                        console.log(error)
                    )
            },
            addArticle(article) {
                if (!this.articlesSortedByFirstLetter[article.title[0].toLowerCase()]) {
                    this.articlesSortedByFirstLetter[article.title[0].toLowerCase()] = []
                }
                this.articlesSortedByFirstLetter[article.title[0].toLowerCase()].push(article)
                const ordered = {}
                Object.keys(this.articlesSortedByFirstLetter).sort().forEach(key => ordered[key] = this.articlesSortedByFirstLetter[key])
                this.articlesSortedByFirstLetter = ordered
                this.$forceUpdate()
                this.addDialog = false
                this.snackbar = {
                    color: "success",
                    text: "Article created",
                    show: true,
                }
            },
            showDeleteAction(id, title, URL) {
                this.deleteDialog = true;
                this.deleteArticle = {
                    id,
                    title,
                    URL
                }
            },
            sendDeleteArticleRequest() {
                this.axios.delete("/articles/" + this.deleteArticle.id)
                    .then(response => {
                        if (response.status !== 204) {
                            console.log("failed", response);
                            return
                        }
                    })
                    .catch(error => {
                        console.log(error, error.response)
                    });
                let firstLetter = this.deleteArticle.title[0].toLowerCase();
                Object.keys(this.articlesSortedByFirstLetter).forEach((letter) => letter !== firstLetter ? true : this.articlesSortedByFirstLetter[letter] = Object.values(this.articlesSortedByFirstLetter[letter]).filter(article => article.ID !== this.deleteArticle.id))
                Object.keys(this.articlesSortedByFirstLetter).forEach(letter => this.articlesSortedByFirstLetter[letter].length === 0 ? delete this.articlesSortedByFirstLetter[letter] : true)
                this.deleteDialog = false;
                this.$forceUpdate()
                this.snackbar = {
                    show: true,
                    color: "green",
                    text: "Article deleted",
                }
            },
        },
        mounted: function () {
            this.loadArticles()
        }
    }
</script>

<style scoped>
    .column_wrapper {
        max-height: 800px;
        display: flex;
        flex-flow: column wrap;
    }

    a:link, a:visited {
        color: black;
    }
</style>