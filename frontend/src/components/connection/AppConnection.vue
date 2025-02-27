<script setup lang="ts">
import { computed } from "vue";
import { useConnections } from "@/composables/useConnections";
import { useUrlParams } from "@/composables/useUrlParams";
import { app } from "_/go/models";
import { useWails } from "@/composables/useWails";
import { DeleteConnection } from "_/go/app/App";
import { useRouter } from "vue-router";

const { connection } = defineProps<{ connection: app.Connection }>();
const emit = defineEmits<{ connectionEdit: [app.Connection] }>();

const { isConnected, connect, disconnect, select, fetchConnections } =
  useConnections();
const { databaseId } = useUrlParams();
const wails = useWails();
const router = useRouter();

const connected = computed(() => isConnected(connection.id));

function getConnectionName(connection: app.Connection) {
  if (connection.name) {
    return connection.name;
  }
  const parts = connection.connection_string.replace(/\\/g, "/").split("/");
  return parts[parts.length - 1];
}

async function removeConnection(connection: app.Connection) {
  const result = await wails(() => DeleteConnection(connection.id));
  if (result instanceof Error) {
    // TODO: specific error handling
  } else {
    fetchConnections();
    if (connection.id === databaseId.value) {
      router.push("/");
    }
  }
}
</script>

<template>
  <UCard
    :ui="{
      root:
        connection.id === databaseId
          ? 'border-r-2 border-r-primary-400 transition-colors'
          : connected
            ? 'cursor-pointer border-r-2 border-r-primary-400/50 hover:border-r-primary-400 transition-colors'
            : 'border-r-2 border-r-transparent transition-colors',
      header: 'sm:p-4',
      body: 'sm:p-4',
      footer: 'sm:p-2',
    }"
    @click="select(connection.id)"
  >
    <div class="flex gap-4 items-center">
      <div class="flex flex-initial">
        <UIcon :name="`simple-icons:${connection.type}`" class="size-8" />
      </div>
      <div class="flex flex-auto flex-row gap-2 justify-between">
        <UTooltip :text="connection.name" :content="{ side: 'right' }">
          <span class="line-clamp-1 text-ellipsis">
            {{ getConnectionName(connection) }}
          </span>
        </UTooltip>
        <AppConnectionInfo :connection="connection" />
      </div>
    </div>

    <template #footer>
      <div class="flex gap-2 justify-end">
        <AppPopconfirm
          text="Are you sure?"
          @confirm="removeConnection(connection)"
        >
          <UTooltip text="Remove" :content="{ side: 'left' }">
            <UButton icon="lucide:trash" color="error" variant="soft" />
          </UTooltip>
        </AppPopconfirm>
        <UTooltip text="Edit" :content="{ side: 'top' }">
          <UButton
            icon="lucide:edit"
            color="neutral"
            variant="soft"
            @click="emit('connectionEdit', connection)"
          />
        </UTooltip>
        <UTooltip
          :text="connected ? 'Disconnect' : 'Connect'"
          :content="{ side: 'right' }"
        >
          <UButton
            :icon="connected ? 'lucide:unplug' : 'lucide:plug'"
            :color="connected ? 'warning' : 'primary'"
            variant="soft"
            @click.prevent="
              () => {
                connected ? disconnect(connection.id) : connect(connection.id);
              }
            "
          />
        </UTooltip>
      </div>
    </template>
  </UCard>
</template>
