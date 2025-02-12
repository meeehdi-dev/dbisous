import { onMounted, ref } from "vue";
import { createSharedComposable } from "@vueuse/core";
import { useWails } from "../wails";
import { Effect } from "effect";
import { useRouter } from "vue-router";
import { useUrlParams } from "./useUrlParams";
import { database } from "../../wailsjs/go/models";
import {
  Connect,
  CreateConnection,
  DeleteConnection,
  Disconnect,
  GetConnections,
  UpdateConnection,
} from "../../wailsjs/go/app/App";

export const useConnections = createSharedComposable(() => {
  const wails = useWails();
  const router = useRouter();
  const { databaseId } = useUrlParams();

  const connections = ref<Array<database.Connection>>([]);
  const activeConnections = ref<Array<string>>([]);

  const fetchConnections = async () => {
    return Effect.runPromise(
      wails(GetConnections).pipe(
        Effect.tap((c) => {
          connections.value = c;
        }),
      ),
    );
  };

  const addConnection = async (connection: database.Connection) => {
    return Effect.runPromise(
      wails(() => CreateConnection(connection)).pipe(
        Effect.tap(fetchConnections),
      ),
    );
  };

  const updateConnectionInfo = async (connection: database.Connection) => {
    return Effect.runPromise(
      wails(() => UpdateConnection(connection)).pipe(
        Effect.tap(fetchConnections),
      ),
    );
  };

  const removeConnection = async (id: string) => {
    return Effect.runPromise(
      wails(() => DeleteConnection(id)).pipe(Effect.tap(fetchConnections)),
    );
  };

  const select = async (id: string) => {
    if (activeConnections.value.some((c) => c === id)) {
      router.push(`/database/${id}`);
    }
  };

  const connect = async (id: string) => {
    return Effect.runPromise(
      wails(() => Connect(id)).pipe(
        Effect.tap(() => {
          activeConnections.value.push(id);
          select(id);
        }),
      ),
    );
  };

  const disconnect = async (id: string) => {
    return Effect.runPromise(
      wails(() => Disconnect(id)).pipe(
        Effect.tap(() => {
          activeConnections.value = activeConnections.value.filter(
            (connectionId) => connectionId !== id,
          );
          if (databaseId.value === id) {
            router.push("/");
          }
        }),
      ),
    );
  };

  const getConnectionName = (connection: database.Connection) => {
    if (connection.name) {
      return connection.name;
    }
    const parts = connection.connection_string.split("/");
    return parts[parts.length - 1];
  };

  const isConnected = (id: string) => {
    return activeConnections.value.includes(id);
  };

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
