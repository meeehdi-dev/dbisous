import AppWelcome from "@/components/AppWelcome.vue";
import AppConnection from "@/components/connection/AppConnection.vue";
import AppDatabase from "@/components/connection/AppDatabase.vue";
import AppSchema from "@/components/connection/AppSchema.vue";
import AppTable from "@/components/connection/AppTable.vue";

export enum Route {
  Welcome = "welcome",
  Connection = "connection",
  Database = "database",
  Schema = "schema",
  Table = "table",
}

export const routes = [
  { path: "/", name: Route.Welcome, component: AppWelcome },
  { path: "/connection", name: Route.Connection, component: AppConnection },
  { path: "/database", name: Route.Database, component: AppDatabase },
  { path: "/schema", name: Route.Schema, component: AppSchema },
  { path: "/table", name: Route.Table, component: AppTable },
];
