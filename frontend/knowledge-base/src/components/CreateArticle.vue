<template>
    <v-form
            ref="form"
            v-model="valid"
            :lazy-validation="lazy"
            @submit="submit">
        <v-row>
            <v-col>
                <v-text-field
                        align="center"
                        id="title"
                        messages="Article title"
                        label="Title"
                        v-model="title"
                        :rules="titleRules"
                        required/>
            </v-col>
        </v-row>
        <v-row justify="center">
            <v-col>
                <v-text-field
                        id="URL"
                        messages="Article URL"
                        label="URL"
                        v-model="url"
                        :rules="urlRules"
                        required/>
            </v-col>
        </v-row>
        <v-row justify="center">
            <v-col>
                <v-btn @click="closeDialog" color="red darken-2">Abort</v-btn>
            </v-col>
            <v-col>
                <v-btn @click="submit" color="green"
                       :disabled="!valid">Create
                </v-btn>
            </v-col>
        </v-row>
        <v-row>
            <v-col>
                <v-alert border="top" color="red" v-if="error != null">{{ error }}</v-alert>
            </v-col>
        </v-row>
    </v-form>
</template>

<script>
    export default {
        name: "CreateArticle",
        data: function () {
            return {
                lazy: false,
                valid: true,
                title: "",
                url: "",
                error: null,
                lastTypeAction: null,
                titleRules: [
                    title => title.length > 3
                ],
                urlRules: [
                    url => !!(new RegExp('^(https?:\\/\\/)?' + // protocol
                        '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|' + // domain name
                        '((\\d{1,3}\\.){3}\\d{1,3}))' + // OR ip (v4) address
                        '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*' + // port and path
                        '(\\?[;&a-z\\d%_.~+=-]*)?' + // query string
                        '(\\#[-a-z\\d_]*)?$', 'i')).test(url),
                    url => {
                        if (!(new RegExp('^(https?:\\/\\/)?' + // protocol
                            '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|' + // domain name
                            '((\\d{1,3}\\.){3}\\d{1,3}))' + // OR ip (v4) address
                            '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*' + // port and path
                            '(\\?[;&a-z\\d%_.~+=-]*)?' + // query string
                            '(\\#[-a-z\\d_]*)?$', 'i')).test(url)) {
                            return false
                        }
                        this.lastTypeAction = (new Date()).getTime()
                        setTimeout(() => {
                                if ((new Date()).getTime() - this.lastTypeAction < 500) {
                                    return
                                }
                                this.axios.post('/extracts', {url: url})
                                    .then(response => this.title = response.data.title.substr(0, 30) + '...')
                                    .catch(error => console.log(error, error.response))
                            }, 500
                        )
                        return true
                    }
                ]
            }
        },
        computed: {},
        methods: {
            submit() {
                this.callCreateArticleEndpoint(this.title, this.url)
            },
            callCreateArticleEndpoint(title, url) {
                this.axios.post("/articles", {title: title, url: url})
                    .then(response => {
                        if (response.status !== 201) {
                            console.log("return status code should be 201", response);
                            return
                        }
                        this.reset()
                        this.$emit("articleCreated", response.data);
                    })
                    .catch(error => {
                        if (error.response && error.response.status === 409) {
                            this.error = "Error : Conflict";
                            return
                        }
                        console.log(error, error.response);
                    })
            },
            closeDialog() {
                this.reset()
                this.$emit('closeDialog')
            },
            reset() {
                this.title = ""
                this.url = ""
                this.error = null
            }
        },
    }
</script>

<style scoped>

</style>