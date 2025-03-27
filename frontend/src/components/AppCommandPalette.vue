<script setup lang="ts">
import { useConnections } from "@/composables/shared/useConnections";
import { CommandPaletteGroup, CommandPaletteItem } from "@nuxt/ui";
import { computed, watch } from "vue";

const emit = defineEmits<{ close: [] }>();

const { connect, disconnect, connections, activeConnections } =
  useConnections();

watch([connections, activeConnections], ([a, b]) => {
  console.log({ a, b });
});

const groups = computed<CommandPaletteGroup<CommandPaletteItem>[]>(() => [
  {
    id: "connect",
    label: "Connect to database",
    items: connections.value
      .filter((connection) => !activeConnections.value.includes(connection.id))
      .map((connection) => ({
        prefix: "Connect to",
        label: connection.name || connection.connection_string,
        onSelect: async () => {
          await connect(connection.id);
          emit("close");
        },
      })),
  },
  {
    id: "disconnect",
    label: "Disconnect from database",
    items: connections.value
      .filter((connection) => activeConnections.value.includes(connection.id))
      .map((connection) => ({
        prefix: "Disconnect from",
        label: connection.name || connection.connection_string,
        onSelect: async () => {
          await disconnect(connection.id);
          emit("close");
        },
      })),
  },
  {
    id: "add_connection",
    label: "Add connection",
    items: [
      {
        prefix: "Add",
        label: "SQLite",
        suffix: "database",
        onSelect: () => {
          // TODO:
          emit("close");
        },
      },
    ],
  },
  {
    id: "edit_connection",
    label: "Edit connection",
    items: [
      {
        prefix: "Edit",
        label: "local pg",
        onSelect: () => {
          // TODO:
          emit("close");
        },
      },
    ],
  },
  {
    id: "delete_connection",
    label: "Delete connection",
    items: [
      {
        prefix: "Delete",
        label: "local pg",
        onSelect: () => {
          // TODO:
          emit("close");
        },
      },
    ],
  },
]);
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
