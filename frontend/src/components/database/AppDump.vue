<script setup lang="ts">
import { useConnections } from "@/composables/useConnections";
import { useUrlParams } from "@/composables/useUrlParams";
import { FormSubmitEvent } from "@nuxt/ui/dist/module";
import * as v from "valibot";
import { computed, reactive, ref, watch } from "vue";

enum ExportType {
  SQL = "sql",
}

const { databaseId } = useUrlParams();
const { metadata } = useConnections();

const schema = v.object({
  type: v.enum(ExportType),
});
const parser = v.safeParser(schema);
type Schema = v.InferOutput<typeof schema>;

const state = reactive<Partial<Schema>>({ type: ExportType.SQL });

const types = ref(
  Object.entries(ExportType).map(([label, value]) => ({ label, value })),
);

watch(
  databaseId,
  () => {
    console.log(metadata.value[databaseId.value]);
  },
  { immediate: true },
);

const schemas = computed(() => {
  const md = metadata.value[databaseId.value].columns;
  console.log(md);
  return Object.keys(md);
});

function submitConnection(event: FormSubmitEvent<Schema>) {
  console.log(event);
}

const activeSchema = ref("");
</script>

<template>
  <div class="p-2 w-full">
    <div class="flex gap-4 h-48 w-full">
      <div class="flex flex-auto flex-col">
        <div
          v-for="s of schemas"
          :key="s"
          :class="['flex gap-2 items-center cursor-pointer']"
        >
          <UCheckbox />
          <div
            :class="[
              'flex flex-auto px-2 py-1 rounded justify-between items-center transition-colors',
              activeSchema === s ? 'bg-primary-500/50' : 'bg-transparent',
            ]"
            @click="activeSchema = activeSchema === s ? '' : s"
          >
            <span>{{ s }}</span>
            <UIcon
              name="lucide:chevron-right"
              :class="[
                'transition-opacity text-primary size-5',
                activeSchema === s ? 'opacity-100' : 'opacity-0',
              ]"
            />
          </div>
        </div>
      </div>
      <USeparator orientation="vertical" />
      <div class="flex flex-auto">tables</div>
      <USeparator orientation="vertical" />
      <div class="flex flex-auto">columns</div>
    </div>
    <UForm :schema="parser" :state="state" @submit="submitConnection">
      <UFormField label="Type" name="type">
        <USelect v-model="state.type" :items="types" :ui="{ base: 'w-36' }" />
      </UFormField>
    </UForm>
  </div>
</template>
