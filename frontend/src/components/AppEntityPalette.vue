<script setup lang="ts">
import { useApp } from "@/composables/shared/useApp";
import { useConnections } from "@/composables/shared/useConnections";
import { Route } from "@/router";
import { CommandPaletteGroup, CommandPaletteItem } from "@nuxt/ui";
import { useMagicKeys } from "@vueuse/core";
import { computed, ref, watch } from "vue";
import { useRouter } from "vue-router";

const emit = defineEmits<{ close: [] }>();

const router = useRouter();
const { connection, schema, table } = useApp();
const { metadata } = useConnections();
const keys = useMagicKeys();

const open = ref(false);
const onConfirm = ref<() => void | Promise<void>>();

const groups = computed(() => {
  const groups: CommandPaletteGroup<CommandPaletteItem>[] = [];

  groups.push({
    id: "schema",
    label: "Go to schema",
    items: Object.keys(metadata.value[connection.value].columns).map((s) => ({
      prefix: "Go to",
      label: s,
      onSelect: async () => {
        await router.push({ name: Route.Schema });
        schema.value = s;
        table.value = "";
        emit("close");
      },
    })),
  });

  Object.keys(metadata.value[connection.value].columns).map((s) => {
    groups.push({
      id: "table",
      label: "Go to table",
      items: Object.keys(metadata.value[connection.value].columns[s]).map(
        (t) => ({
          prefix: "Go to",
          label: t,
          onSelect: async () => {
            await router.push({ name: Route.Table });
            schema.value = s;
            table.value = t;
            emit("close");
          },
        }),
      ),
    });
  });

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
