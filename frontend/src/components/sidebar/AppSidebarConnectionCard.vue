<script setup lang="ts">
import { computed, ref } from "vue";
import { useConnections } from "@/composables/shared/useConnections";
import { app } from "_/go/models";
import { useWails } from "@/composables/useWails";
import { DeleteConnection } from "_/go/app/App";
import { useRouter } from "vue-router";
import { useApp } from "@/composables/shared/useApp";
import { Route } from "@/router";

const { value } = defineProps<{ value: app.Connection }>();
const emit = defineEmits<{ edit: [string] }>();

const { isConnected, connect, disconnect, select, fetchConnections } =
  useConnections();
const { connection } = useApp();
const wails = useWails();
const router = useRouter();

const connected = computed(() => isConnected(value.id));

const connecting = ref(false);

function getConnectionName(connection: app.Connection) {
  if (connection.name) {
    return connection.name;
  }
  const parts = connection.connection_string.replace(/\\/g, "/").split("/");
  return parts[parts.length - 1];
}

async function removeConnection(id: string) {
  const result = await wails(() => DeleteConnection(id));
  if (result instanceof Error) {
    return;
  }
  await fetchConnections();
  if (id === connection.value) {
    await router.push({ name: Route.Welcome });
  }
}

async function onConnect(id: string) {
  connecting.value = true;
  await connect(id);
  connecting.value = false;
}

async function onDisconnect(id: string) {
  connecting.value = true;
  await disconnect(id);
  connection.value = "";
  connecting.value = false;
}
</script>

<template>
  <UCard
    :ui="{
      root:
        value.id === connection
          ? 'border-r-primary-400 border-r-2 transition-colors'
          : connected
            ? 'border-r-primary-400/50 hover:border-r-primary-400 cursor-pointer border-r-2 transition-colors'
            : 'border-r-2 border-r-transparent transition-colors',
      header: 'sm:p-2',
      body: 'sm:p-2',
      footer: 'sm:p-2',
    }"
    @click="select(value.id)"
  >
    <div class="flex items-center gap-2">
      <div class="flex flex-initial">
        <UIcon :name="`simple-icons:${value.type}`" class="size-8" />
      </div>
      <div class="flex flex-auto flex-row justify-between gap-2">
        <UTooltip :text="value.name" :content="{ side: 'right' }">
          <span class="line-clamp-1 text-ellipsis">
            {{ getConnectionName(value) }}
          </span>
        </UTooltip>
        <AppSidebarConnectionInfo :connection="value" />
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2">
        <AppPopconfirm
          text="Are you sure?"
          @confirm="removeConnection(value.id)"
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
            @click="emit('edit', value.id)"
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
                connected ? onDisconnect(value.id) : onConnect(value.id);
              }
            "
          />
        </UTooltip>
      </div>
    </template>
  </UCard>
</template>
