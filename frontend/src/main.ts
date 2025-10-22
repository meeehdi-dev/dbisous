import "./assets/main.css";

import { createApp, vaporInteropPlugin } from "vue";
import App from "./App.vue";
import router from "./router";
import ui from "@nuxt/ui/vue-plugin";

const app = createApp(App);

app.use(vaporInteropPlugin).use(router).use(ui);

app.mount("#app");
