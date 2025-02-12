import { onMounted, ref } from "vue";
import { createSharedComposable } from "@vueuse/core";
import { useWails } from "../wails";
import { Effect } from "effect";
import {
  ConnectToDatabase,
  CreateDatabase,
  DeleteDatabase,
  DisconnectFromDatabase,
  GetDatabases,
  UpdateDatabase,
} from "../../wailsjs/go/app/App";
import { database } from "../../wailsjs/go/models";
import { useRouter } from "vue-router";
import { useUrlParams } from "./useUrlParams";

export const useDatabase = createSharedComposable(() => {
  const wails = useWails();
  const router = useRouter();
  const { databaseId } = useUrlParams();

  const databases = ref<Array<database.Database>>([]);
  const activeConnections = ref<Array<string>>([]);

  const fetchDatabases = async () => {
    return Effect.runPromise(
      wails(GetDatabases).pipe(
        Effect.tap((c) => {
          databases.value = c;
        }),
      ),
    );
  };

  const addDatabase = async (dbInfo: database.Database) => {
    return Effect.runPromise(
      wails(() => CreateDatabase(dbInfo)).pipe(Effect.tap(fetchDatabases)),
    );
  };

  const updateDatabaseInfo = async (dbInfo: database.Database) => {
    return Effect.runPromise(
      wails(() => UpdateDatabase(dbInfo)).pipe(Effect.tap(fetchDatabases)),
    );
  };

  const removeDatabase = async (id: string) => {
    return Effect.runPromise(
      wails(() => DeleteDatabase(id)).pipe(Effect.tap(fetchDatabases)),
    );
  };

  const selectDatabase = async (id: string) => {
    router.push(`/database/${id}/schema`);
  };

  const connectDatabase = async (id: string) => {
    return Effect.runPromise(
      wails(() => ConnectToDatabase(id)).pipe(
        Effect.tap(() => {
          activeConnections.value.push(id);
          selectDatabase(id);
        }),
      ),
    );
  };

  const disconnectDatabase = async (id: string) => {
    return Effect.runPromise(
      wails(() => DisconnectFromDatabase(id)).pipe(
        Effect.tap(() => {
          activeConnections.value = activeConnections.value.filter(
            (connId) => connId !== id,
          );
          if (databaseId.value === id) {
            router.push("/");
          }
        }),
      ),
    );
  };

  const getDatabaseName = (database: database.Database) => {
    if (database.name) {
      return database.name;
    }
    const parts = database.connection_string.split("/");
    return parts[parts.length - 1];
  };

  const isConnected = (id: string) => {
    return activeConnections.value.includes(id);
  };

  onMounted(fetchDatabases);

  return {
    databases,
    activeConnections,
    fetchDatabases,
    addDatabase,
    updateDatabaseInfo,
    removeDatabase,
    connectDatabase,
    disconnectDatabase,
    selectDatabase,
    getDatabaseName,
    isConnected,
  };
});
