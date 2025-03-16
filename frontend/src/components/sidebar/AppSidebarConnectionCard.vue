<script setup lang="ts">
import { computed, ref } from "vue";
import { useConnections } from "@/composables/shared/useConnections";
import { app } from "_/go/models";
import { useWails } from "@/composables/useWails";
import { DeleteConnection } from "_/go/app/App";
import { useRouter } from "vue-router";
import { useApp } from "@/composables/shared/useApp";
import { Route } from "@/router";

const { connection } = defineProps<{ connection: app.Connection }>();
const emit = defineEmits<{ connectionEdit: [app.Connection] }>();

const { isConnected, connect, disconnect, select, fetchConnections } =
  useConnections();
const { database } = useApp();
const wails = useWails();
const router = useRouter();

const connected = computed(() => isConnected(connection.id));

const connecting = ref(false);

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
    return;
  }
  await fetchConnections();
  if (connection.id === database.value) {
    await router.push({ name: Route.Welcome });
  }
}

async function onConnect(connection: app.Connection) {
  connecting.value = true;
  await connect(connection.id);
  connecting.value = false;
}

async function onDisconnect(connection: app.Connection) {
  connecting.value = true;
  await disconnect(connection.id);
  connecting.value = false;
}
</script>

<template>
  <UCard
    :ui="{
      root:
        connection.id === database
          ? 'border-r-2 border-r-primary-400 transition-colors'
          : connected
            ? 'cursor-pointer border-r-2 border-r-primary-400/50 hover:border-r-primary-400 transition-colors'
            : 'border-r-2 border-r-transparent transition-colors',
      header: 'sm:p-2',
      body: 'sm:p-2',
      footer: 'sm:p-2',
    }"
    @click="select(connection.id)"
  >
    <div class="flex items-center gap-2">
      <div class="flex flex-initial">
        <UIcon :name="`simple-icons:${connection.type}`" class="size-8" />
      </div>
      <div class="flex flex-auto flex-row justify-between gap-2">
        <UTooltip :text="connection.name" :content="{ side: 'right' }">
          <span class="line-clamp-1 text-ellipsis">
            {{ getConnectionName(connection) }}
          </span>
        </UTooltip>
        <AppSidebarConnectionInfo :connection="connection" />
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2">
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
            :loading="connecting"
            variant="soft"
            @click.prevent="
              () => {
                connected ? onDisconnect(connection) : onConnect(connection);
              }
            "
          />
        </UTooltip>
      </div>
    </template>
  </UCard>
</template>
