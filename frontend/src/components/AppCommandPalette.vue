<script setup lang="ts">
import { useConnections } from "@/composables/shared/useConnections";
import { useSidebar } from "@/composables/shared/useSidebar";
import { CommandPaletteGroup, CommandPaletteItem } from "@nuxt/ui";
import { useMagicKeys } from "@vueuse/core";
import { app } from "_/go/models";
import { computed, ref, watch } from "vue";

const emit = defineEmits<{ close: [] }>();

const {
  connect,
  disconnect,
  removeConnection,
  connections,
  activeConnections,
} = useConnections();
const { slideoverOpen, editedConnection } = useSidebar();
const keys = useMagicKeys();

const open = ref(false);
const onConfirm = ref<() => void | Promise<void>>();

const connectable = connections.value.filter(
  (connection) => !activeConnections.value.includes(connection.id),
);
const disconnectable = connections.value.filter((connection) =>
  activeConnections.value.includes(connection.id),
);
const groups = computed(() => {
  const groups: CommandPaletteGroup<CommandPaletteItem>[] = [];

  if (disconnectable.length > 0) {
    groups.push({
      id: "disconnect",
      label: "Disconnect from database",
      items: disconnectable.map((connection) => ({
        prefix: "Disconnect from",
        label: connection.name || connection.connection_string,
        onSelect: async () => {
          await disconnect(connection.id);
          emit("close");
        },
      })),
    });
  }
  if (connectable.length > 0) {
    groups.push({
      id: "connect",
      label: "Connect to database",
      items: connectable.map((connection) => ({
        prefix: "Connect to",
        label: connection.name || connection.connection_string,
        onSelect: async () => {
          await connect(connection.id);
          emit("close");
        },
      })),
    });
  }
  groups.push(
    {
      id: "add_connection",
      label: "Add connection",
      items: Object.entries(app.ConnectionType).map(([key, value]) => ({
        prefix: "Add",
        label: key,
        suffix: "database",
        onSelect: () => {
          editedConnection.value = {
            type: value,
            name: "",
            connection_string: "",
          };
          slideoverOpen.value = true;
          emit("close");
        },
      })),
    },
    {
      id: "edit_connection",
      label: "Edit connection",
      items: connections.value.map((connection) => ({
        prefix: "Edit",
        label: connection.name,
        onSelect: () => {
          editedConnection.value = connection;
          slideoverOpen.value = true;
          emit("close");
        },
      })),
    },
    {
      id: "remove_connection",
      label: "Remove connection",
      items: connections.value.map((connection) => ({
        prefix: "Remove",
        label: connection.name,
        onSelect: () => {
          onConfirm.value = async () => {
            await removeConnection(connection.id);
            emit("close");
          };
          open.value = true;
        },
      })),
    },
  );

  // eslint-disable-next-line @typescript-eslint/no-unsafe-return
  return groups;
});

watch(keys["enter"], async (enter) => {
  if (!enter || !open.value || !onConfirm.value) {
    return;
  }

  await onConfirm.value();
});
</script>

<template>
  <UModal>
    <template #content>
      <UCommandPalette
        :fuse="{ fuseOptions: { keys: ['prefix', 'label', 'suffix'] } }"
        :groups="groups"
        placeholder="Search command..."
        class="h-80"
        :ui="{
          itemLabelPrefix:
            'text-(--ui-text-dimmed) [&>mark]:text-(--ui-bg) [&>mark]:bg-(--ui-primary)',
          itemLabelSuffix:
            'text-(--ui-text-dimmed) [&>mark]:text-(--ui-bg) [&>mark]:bg-(--ui-primary)',
        }"
      />
      <UModal
        v-model:open="open"
        title="Are you sure?"
        :ui="{ footer: 'justify-end' }"
      >
        <template #footer>
          <UButton
            icon="lucide:x"
            color="error"
            variant="soft"
            label="Cancel"
            @click="open = false"
          />
          <UButton
            icon="lucide:check"
            variant="soft"
            label="Confirm"
            @click="onConfirm"
          />
        </template>
      </UModal>
    </template>
  </UModal>
</template>
