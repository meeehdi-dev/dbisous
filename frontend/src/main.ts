import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import ui from "@nuxt/ui/vue-plugin";
import App from "./App.vue";
import AppWelcome from "./components/AppWelcome.vue";
import AppDatabase from "./components/database/AppDatabase.vue";
import AppSchema from "./components/database/AppSchema.vue";
import AppTable from "./components/database/AppTable.vue";
import "./style.css";

const app = createApp(App);

const router = createRouter({
  routes: [
    { path: "/", name: "welcome", component: AppWelcome },
    {
      path: "/database/:databaseId",
      children: [
        { path: "", name: "database", component: AppDatabase },
        {
          path: "schema/:schemaId",
          children: [
            { path: "", name: "schema", component: AppSchema },
            { path: "table/:tableId", name: "table", component: AppTable },
          ],
        },
      ],
    },
  ],
  history: createWebHistory(),
});

app.use(router);
app.use(ui);

app.mount("#app");
