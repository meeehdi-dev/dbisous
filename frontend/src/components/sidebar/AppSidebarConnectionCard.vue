<script setup lang="ts">
import { computed, ref } from "vue";
import { useConnections } from "@/composables/shared/useConnections";
import { app } from "_/go/models";
import { useWails } from "@/composables/useWails";
import { DeleteConnection, TestConnection } from "_/go/app/App";
import { useRouter } from "vue-router";
import { useApp } from "@/composables/shared/useApp";
import { Route } from "@/router";

const { value } = defineProps<{ value: app.Connection }>();
const emit = defineEmits<{ edit: [string]; duplicate: [string] }>();

const { isConnected, connect, disconnect, select, fetchConnections } =
  useConnections();
const { connection } = useApp();
const wails = useWails();
const router = useRouter();
// eslint-disable-next-line no-undef
const toast = useToast();

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
  connecting.value = false;
}

async function testConnection(connection: app.Connection) {
  const result = await wails(() =>
    TestConnection(connection.type, connection.connection_string),
  );
  if (result instanceof Error) {
    return;
  }
  toast.add({
    title: "Test connection",
    description: "Successfully pinged database!",
  });
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
        <UTooltip
          :text="connected ? 'Disconnect' : 'Connect'"
          :content="{ side: 'right' }"
        >
          <UButton
            :icon="connected ? 'lucide:unplug' : 'lucide:plug'"
            :color="connected ? 'warning' : 'primary'"
            :loading="connecting"
            variant="soft"
            @click.stop="
              () => {
                connected ? onDisconnect(value.id) : onConnect(value.id);
              }
            "
          />
        </UTooltip>
        <UPopover
          :content="{ side: 'right', align: 'start' }"
          :ui="{ content: 'flex flex-col gap-2 p-2' }"
          arrow
        >
          <UButton
            icon="lucide:ellipsis-vertical"
            variant="soft"
            color="neutral"
          />
          <template #content>
            <UButton
              icon="lucide:plug-zap"
              variant="soft"
              color="success"
              label="Test"
              @click="testConnection(value)"
            />
            <UButton
              icon="lucide:edit"
              color="neutral"
              variant="soft"
              label="Edit"
              @click="emit('edit', value.id)"
            />
            <UButton
              icon="lucide:copy"
              color="secondary"
              variant="soft"
              label="Duplicate"
              @click="emit('duplicate', value.id)"
            />
            <AppPopconfirm
              text="Are you sure?"
              @confirm="removeConnection(value.id)"
            >
              <UButton
                icon="lucide:trash"
                color="error"
                variant="soft"
                label="Delete"
              />
            </AppPopconfirm>
          </template>
        </UPopover>
      </div>
    </template>
  </UCard>
</template>
