<script setup lang="ts">
import { useApp } from "@/composables/shared/useApp";
import { useConnections } from "@/composables/shared/useConnections";
import { useSidebar } from "@/composables/shared/useSidebar";
import { CommandPaletteGroup, CommandPaletteItem } from "@nuxt/ui";
import { useMagicKeys } from "@vueuse/core";
import { app } from "_/go/models";
import { computed, ref, watch } from "vue";

const emit = defineEmits<{ close: [] }>();

const { connection } = useApp();
const {
  connect,
  disconnect,
  removeConnection,
  connections,
  activeConnections,
  select,
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
      id: "select",
      label: "Select database",
      items: disconnectable
        .filter((c) => c.id !== connection.value)
        .map((c) => ({
          prefix: "Select",
          label: c.name || c.connection_string,
          onSelect: async () => {
            await select(c.id);
            emit("close");
          },
        })),
    });
    groups.push({
      id: "disconnect",
      label: "Disconnect from database",
      items: disconnectable.map((c) => ({
        prefix: "Disconnect from",
        label: c.name || c.connection_string,
        onSelect: async () => {
          await disconnect(c.id);
          emit("close");
        },
      })),
    });
  }
  if (connectable.length > 0) {
    groups.push({
      id: "connect",
      label: "Connect to database",
      items: connectable.map((c) => ({
        prefix: "Connect to",
        label: c.name || c.connection_string,
        onSelect: async () => {
          await connect(c.id);
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
      items: connections.value.map((c) => ({
        prefix: "Edit",
        label: c.name,
        onSelect: () => {
          editedConnection.value = c;
          slideoverOpen.value = true;
          emit("close");
        },
      })),
    },
    {
      id: "remove_connection",
      label: "Remove connection",
      items: connections.value.map((c) => ({
        prefix: "Remove",
        label: c.name,
        onSelect: () => {
          onConfirm.value = async () => {
            await removeConnection(c.id);
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
