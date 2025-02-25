import { onMounted, ref } from "vue";
import { createSharedComposable } from "@vueuse/core";
import { useWails } from "@/composables/useWails";
import { useRouter } from "vue-router";
import { useUrlParams } from "@/composables/useUrlParams";
import { app } from "_/go/models";
import {
  Connect,
  CreateConnection,
  DeleteConnection,
  Disconnect,
  GetConnections,
  UpdateConnection,
} from "_/go/app/App";

export const useConnections = createSharedComposable(() => {
  const wails = useWails();
  const router = useRouter();
  const { databaseId } = useUrlParams();

  const connections = ref<Array<app.Connection>>([]);
  const activeConnections = ref<Array<string>>([]);

  async function fetchConnections() {
    const result = await wails(GetConnections);
    if (result instanceof Error) {
      // TODO: specific error handling
    } else {
      connections.value = result;
    }
  }

  async function addConnection(connection: app.Connection) {
    const result = await wails(() => CreateConnection(connection));
    if (result instanceof Error) {
      // TODO: specific error handling
    } else {
      await fetchConnections();
    }
  }

  async function updateConnectionInfo(connection: app.Connection) {
    const result = await wails(() => UpdateConnection(connection));
    if (result instanceof Error) {
      // TODO: specific error handling
    } else {
      await fetchConnections();
    }
  }

  async function removeConnection(id: string) {
    const result = await wails(() => DeleteConnection(id));
    if (result instanceof Error) {
      // TODO: specific error handling
    } else {
      await fetchConnections();
    }
  }

  function select(id: string) {
    if (
      activeConnections.value.some((c) => c === id) &&
      databaseId.value !== id
    ) {
      router.push(`/database/${id}`);
    }
  }

  async function connect(id: string) {
    const result = await wails(() => Connect(id));
    if (result instanceof Error) {
      // TODO: specific error handling
    } else {
      activeConnections.value.push(id);
      select(id);
    }
  }

  async function disconnect(id: string) {
    const result = await wails(() => Disconnect(id));
    if (result instanceof Error) {
      // TODO: specific error handling
    } else {
      activeConnections.value = activeConnections.value.filter(
        (connectionId) => connectionId !== id,
      );
      if (databaseId.value === id) {
        router.push("/");
      }
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
    updateConnectionInfo,
    removeConnection,
    connect,
    disconnect,
    select,
    getConnectionName,
    isConnected,
  };
});
