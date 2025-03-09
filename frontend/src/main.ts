import "@/style.css";

import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import ui from "@nuxt/ui/vue-plugin";
import App from "@/App.vue";
import { routes } from "./router";

const app = createApp(App);

const router = createRouter({
  routes,
  history: createWebHistory(),
});

app.use(router);
app.use(ui);

app.mount("#app");
