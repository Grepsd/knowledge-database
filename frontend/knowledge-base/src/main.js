import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import ArticlesList from "./components/ArticlesList";
import '@mdi/font/css/materialdesignicons.css'
import axios from 'axios'
import VueAxios from "vue-axios";

Vue.config.productionTip = false
Vue.use(VueAxios, axios.create({
  baseURL: process.env.VUE_APP_BASE_URL
}))

new Vue({
  vuetify,
  render: h => h(App)
}).$mount('#app')
Vue.component('articles-list', ArticlesList)
