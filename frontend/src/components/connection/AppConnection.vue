<script setup lang="ts">
import { database } from "../../../wailsjs/go/models";
import { computed } from "vue";
import { useConnections } from "../../composables/useConnections";
import { useUrlParams } from "../../composables/useUrlParams";

const { connection } = defineProps<{ connection: database.Connection }>();

const { isConnected, connect, disconnect, select, activeConnections } =
  useConnections();
const { databaseId } = useUrlParams();
const toast = useToast();

const connected = computed(() => isConnected(connection.id));

function getConnectionName(connection: database.Connection) {
  if (connection.name) {
    return connection.name;
  }
  const parts = connection.connection_string.split("/");
  return parts[parts.length - 1];
}

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text);
  toast.add({
    title: "Successfully copied to clipboard!",
    description: text,
  });
}
</script>

<template>
  <UCard
    :ui="{
      root:
        activeConnections.length > 0
          ? connection.id === databaseId
            ? 'opacity-100'
            : connected
              ? 'opacity-80'
              : 'opacity-60'
          : 'opacity-100',
    }"
  >
    <div class="flex gap-4 items-center">
      <UIcon
        :name="`simple-icons:${connection.type}`"
        class="size-8 text-primary-400"
      />
      <div class="w-full flex flex-col">
        <div class="flex flex-1 flex-row gap-2 items-center justify-between">
          <span>
            {{ getConnectionName(connection) }}
          </span>
          <UPopover mode="hover" :content="{ side: 'right' }">
            <UIcon
              name="lucide:info"
              class="size-6 text-secondary-400/50 hover:text-secondary-400 transition-colors"
            />
            <template #content>
              <div class="p-2 flex flex-col gap-2 text-gray-400">
                <UTooltip text="Connection string" :content="{ side: 'left' }">
                  <div class="flex flex-row gap-2 items-center">
                    <UIcon name="lucide:link" class="text-secondary-400" />
                    <UButton
                      color="neutral"
                      variant="ghost"
                      trailing-icon="lucide:copy"
                      :ui="{ base: 'px-1' }"
                      :label="connection.connection_string"
                      @click="copyToClipboard(connection.connection_string)"
                    />
                  </div>
                </UTooltip>
                <UTooltip text="Creation date" :content="{ side: 'left' }">
                  <div class="flex flex-row gap-2 items-center">
                    <UIcon name="lucide:calendar" class="text-primary-400/50" />
                    <span class="text-sm">
                      {{ new Date(connection.created_at).toLocaleString() }}
                    </span>
                  </div>
                </UTooltip>
                <UTooltip text="Last update" :content="{ side: 'left' }">
                  <div
                    v-if="connection.created_at !== connection.updated_at"
                    class="flex flex-row gap-2 items-center"
                  >
                    <UIcon
                      name="lucide:calendar-sync"
                      class="text-secondary-400"
                    />
                    <span class="text-sm">
                      {{ new Date(connection.created_at).toLocaleString() }}
                    </span>
                  </div>
                </UTooltip>
              </div>
            </template>
          </UPopover>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-between">
        <UTooltip
          :text="connected ? 'Disconnect' : 'Connect'"
          :content="{ side: 'right' }"
        >
          <UChip
            :show="connected"
            color="primary"
            :ui="{ base: 'animate-ping' }"
          >
            <UButton
              :icon="connected ? 'lucide:unplug' : 'lucide:plug'"
              :color="connected ? 'warning' : 'primary'"
              variant="soft"
              @click="
                () => {
                  connected
                    ? disconnect(connection.id)
                    : connect(connection.id);
                }
              "
            />
          </UChip>
        </UTooltip>

        <div class="flex gap-2 justify-end">
          <AppConnectionRemoveButton :connection="connection" />
          <UTooltip text="Edit" :content="{ side: 'top' }">
            <UButton icon="lucide:edit" color="neutral" variant="soft" />
          </UTooltip>
          <UTooltip text="Select" :content="{ side: 'right' }">
            <UButton
              icon="lucide:play"
              @click="select(connection.id)"
              :color="connected ? 'primary' : 'neutral'"
              :variant="
                !connected || connection.id === databaseId ? 'soft' : 'solid'
              "
              :disabled="!connected || connection.id === databaseId"
            />
          </UTooltip>
        </div>
      </div>
    </template>
  </UCard>
</template>
