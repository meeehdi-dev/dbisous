import { onMounted, ref } from "vue";
import { createSharedComposable } from "@vueuse/core";
import { useWails } from "@/composables/useWails";
import { useRouter } from "vue-router";
import { app } from "_/go/models";
import {
  Connect,
  CreateConnection,
  DeleteConnection,
  Disconnect,
  GetConnections,
  UpdateConnection,
} from "_/go/app/App";
import { useCompletions } from "@/composables/shared/useCompletions";
import { Route } from "@/router";
import { useApp } from "@/composables/shared/useApp";

type DatabaseMetadata = Record<string, Record<string, Array<string>>>;

export const useConnections = createSharedComposable(() => {
  const wails = useWails();
  const router = useRouter();
  const { connection: database } = useApp();

  const connections = ref<Array<app.Connection>>([]);
  const activeConnections = ref<Array<string>>([]);
  const metadata = ref<Record<string, { columns: DatabaseMetadata }>>({});

  async function fetchConnections() {
    const result = await wails(GetConnections);
    if (result instanceof Error) {
      return;
    }
    connections.value = result;
  }

  async function addConnection(connection: app.Connection) {
    const result = await wails(() => CreateConnection(connection));
    if (result instanceof Error) {
      return;
    }
    await fetchConnections();
  }

  async function updateConnection(connection: app.Connection) {
    const result = await wails(() => UpdateConnection(connection));
    if (result instanceof Error) {
      return;
    }
    await fetchConnections();
  }

  async function removeConnection(id: string) {
    const result = await wails(() => DeleteConnection(id));
    if (result instanceof Error) {
      return;
    }
    await fetchConnections();
  }

  const { register } = useCompletions();

  async function select(id: string) {
    if (
      activeConnections.value.some((c) => c === id) &&
      database.value !== id
    ) {
      register(metadata.value[id].columns);
      database.value = id;
      await router.push({ name: Route.Database });
    }
  }

  async function connect(id: string) {
    const result = await wails(() => Connect(id));
    if (result instanceof Error) {
      return;
    }
    metadata.value[id] = result;
    activeConnections.value.push(id);
    await select(id);
  }

  async function disconnect(id: string) {
    const result = await wails(() => Disconnect(id));
    if (result instanceof Error) {
      return;
    }
    activeConnections.value = activeConnections.value.filter(
      (connectionId) => connectionId !== id,
    );
    if (database.value === id) {
      await router.push({ name: Route.Welcome });
    }
  }

  function getConnectionName(connection: app.Connection) {
    if (connection.name) {
      return connection.name;
    }
    const parts = connection.connection_string.split("/");
    return parts[parts.length - 1];
  }

  function isConnected(id: string) {
    return activeConnections.value.includes(id);
  }

  onMounted(fetchConnections);

  return {
    connections,
    activeConnections,
    fetchConnections,
    addConnection,
    updateConnection,
    removeConnection,
    connect,
    disconnect,
    select,
    getConnectionName,
    isConnected,
    metadata,
  };
});
