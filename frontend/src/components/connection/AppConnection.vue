<script setup lang="ts">
import { database } from "../../../wailsjs/go/models";
import { computed } from "vue";
import { useDatabase } from "../../composables/useDatabase";
import { useUrlParams } from "../../composables/useUrlParams";

const { connection } = defineProps<{ connection: database.Database }>();

const { isConnected, connectDatabase, disconnectDatabase, selectDatabase } =
  useDatabase();
const { databaseId } = useUrlParams();

const connected = computed(() => isConnected(connection.id));

const getDatabaseName = (database: database.Database) => {
  if (database.name) {
    return database.name;
  }
  const parts = database.connection_string.split("/");
  return parts[parts.length - 1];
};
</script>

<template>
  <UCard
    class="w-full"
    :ui="{ root: connection.id === databaseId ? 'ring-primary-400/50' : '' }"
  >
    <div class="flex gap-4 items-center">
      <UIcon :name="`simple-icons:${connection.type}`" class="size-8" />
      <UTooltip
        :text="connection.connection_string"
        :content="{
          side: 'right',
        }"
      >
        <span>
          {{ getDatabaseName(connection) }}
        </span>
      </UTooltip>
    </div>

    <template #footer>
      <div class="flex justify-between">
        <UTooltip
          :text="connected ? 'Disconnect' : 'Connect'"
          :content="{ side: 'right' }"
        >
          <UChip
            :color="connected ? 'success' : 'neutral'"
            :ui="{ base: connected ? 'animate-ping' : undefined }"
          >
            <UButton
              :icon="connected ? 'lucide:unplug' : 'lucide:plug'"
              :color="connected ? 'warning' : 'success'"
              variant="soft"
              @click="
                () => {
                  connected
                    ? disconnectDatabase(connection.id)
                    : connectDatabase(connection.id);
                }
              "
            />
          </UChip>
        </UTooltip>

        <div class="flex gap-2 justify-end">
          <AppConnectionRemoveButton :connection="connection" />
          <UTooltip text="Select" :content="{ side: 'right' }">
            <UButton
              icon="lucide:play"
              @click="selectDatabase(connection.id)"
              :color="connected ? 'success' : 'neutral'"
              variant="soft"
              :disabled="!connected"
            />
          </UTooltip>
        </div>
      </div>
    </template>
  </UCard>
</template>
