<script setup lang="ts">
import { useConnections } from "@/composables/shared/useConnections";
import { useSidebar } from "@/composables/shared/useSidebar";
import { CommandPaletteGroup, CommandPaletteItem } from "@nuxt/ui";
import { app } from "_/go/models";
import { computed } from "vue";

const emit = defineEmits<{ close: [] }>();

const { connect, disconnect, connections, activeConnections } =
  useConnections();
const { slideoverOpen, editedConnection } = useSidebar();

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
      id: "delete_connection",
      label: "Delete connection",
      items: connections.value.map((connection) => ({
        prefix: "Delete",
        label: connection.name,
        onSelect: () => {
          // TODO: delete
          emit("close");
        },
      })),
    },
  );

  // eslint-disable-next-line @typescript-eslint/no-unsafe-return
  return groups;
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
    </template>
  </UModal>
</template>
