<script setup lang="ts">
import { useApp } from "@/composables/shared/useApp";
import { useConnections } from "@/composables/shared/useConnections";
import { useWails } from "@/composables/useWails";
import type { FormSubmitEvent } from "@nuxt/ui";
import { Export } from "_/go/app/App";
import { client } from "_/go/models";
import * as v from "valibot";
import { computed, reactive, ref } from "vue";

const { connection } = useApp();
const { metadata } = useConnections();
const wails = useWails();
// eslint-disable-next-line no-undef
const toast = useToast();

const exportSchema = v.object({
  type: v.enum(client.ExportType),
  schema_only: v.boolean(),
  data_only: v.boolean(),
  drop_schema: v.boolean(),
  ignore_constraints: v.boolean(),
  wrap_in_transaction: v.boolean(),
  drop_table: v.enum(client.ExportDrop),
  selected: v.record(
    v.string(),
    v.union([v.boolean(), v.literal("indeterminate")]),
  ),
});
const parser = v.safeParser(exportSchema);
type ExportSchema = v.InferOutput<typeof exportSchema>;

const state = reactive<ExportSchema>({
  type: client.ExportType.SQL,
  schema_only: false,
  data_only: false,
  drop_schema: false,
  ignore_constraints: false,
  wrap_in_transaction: true,
  drop_table: client.ExportDrop.Drop_and_create,
  selected: {},
});

const types = ref(
  Object.entries(client.ExportType).map(([label, value]) => ({ label, value })),
);
const drop = computed(() =>
  Object.entries(client.ExportDrop).map(([label, value]) => ({
    label: label.replace(/_/g, " "),
    value,
  })),
);

const schemas = computed(() => {
  const md = metadata.value[connection.value].columns;
  return Object.keys(md);
});

const tables = computed(() => {
  if (activeSchema.value === "") {
    return [];
  }
  const md = metadata.value[connection.value].columns;
  return Object.keys(md[activeSchema.value]);
});

const columns = computed(() => {
  if (activeSchema.value === "" || activeTable.value === "") {
    return [];
  }
  const md = metadata.value[connection.value].columns;
  return md[activeSchema.value][activeTable.value];
});

async function submit(event: FormSubmitEvent<ExportSchema>) {
  const result = await wails(() =>
    Export(connection.value, {
      ...event.data,
      selected: Object.entries(state.selected)
        .filter(([, value]) => value !== false)
        .map(([key]) => key),
    }),
  );
  if (result instanceof Error) {
    return;
  }
  toast.add({
    title: "Successfully exported database!",
    description: result,
  });
}

const activeSchema = ref("");
const activeTable = ref("");

function viewSchema(schema: string) {
  if (activeSchema.value === schema) {
    activeSchema.value = "";
  } else {
    activeSchema.value = schema;
  }
}

function viewTable(table: string) {
  if (activeTable.value === table) {
    activeTable.value = "";
  } else {
    activeTable.value = table;
  }
}

function selectSchema(schema: string) {
  if (state.selected[schema]) {
    const md = metadata.value[connection.value].columns;
    const schemas = Object.keys(md);
    const otherSchemas = schemas.filter((s) => s !== schema);
    otherSchemas.forEach((otherSchema) => {
      state.selected[otherSchema] = false;
      Object.keys(state.selected).forEach((key) => {
        if (key.startsWith(`${otherSchema}.`)) {
          state.selected[key] = false;
        }
      });
    });
    const tables = Object.keys(md[schema]);
    tables.forEach((table) => {
      state.selected[`${schema}.${table}`] = true;
      const columns = md[schema][table];
      columns.forEach((column) => {
        state.selected[`${schema}.${table}.${column}`] = true;
      });
    });
    viewSchema(schema);
  } else {
    Object.keys(state.selected).forEach((key) => {
      if (key.startsWith(`${schema}.`)) {
        state.selected[key] = false;
      }
    });
  }
}

function selectTable(table: string) {
  const md = metadata.value[connection.value].columns;
  if (state.selected[`${activeSchema.value}.${table}`] === true) {
    if (state.selected[activeSchema.value] !== true) {
      const schemaTables = Object.keys(md[activeSchema.value]);
      if (
        schemaTables.every((t) => state.selected[`${activeSchema.value}.${t}`])
      ) {
        state.selected[activeSchema.value] = true;
      } else {
        state.selected[activeSchema.value] = "indeterminate";
      }
    }
    const columns = md[activeSchema.value][table];
    columns.forEach((column) => {
      state.selected[`${activeSchema.value}.${table}.${column}`] = true;
    });
    viewTable(table);
  } else {
    if (state.selected[activeSchema.value] !== false) {
      const schemaTables = Object.keys(md[activeSchema.value]);
      if (
        schemaTables.some((t) => state.selected[`${activeSchema.value}.${t}`])
      ) {
        state.selected[activeSchema.value] = "indeterminate";
      } else {
        state.selected[activeSchema.value] = false;
      }
    }
    Object.keys(state.selected).forEach((key) => {
      if (key.startsWith(`${activeSchema.value}.${table}.`)) {
        state.selected[key] = false;
      }
    });
  }
}

function selectColumn(column: string) {
  const md = metadata.value[connection.value].columns;
  if (
    state.selected[`${activeSchema.value}.${activeTable.value}.${column}`] ===
    true
  ) {
    if (state.selected[`${activeSchema.value}.${activeTable.value}`] !== true) {
      const tableColumns = md[activeSchema.value][activeTable.value];
      if (
        tableColumns.every(
          (c) =>
            state.selected[`${activeSchema.value}.${activeTable.value}.${c}`],
        )
      ) {
        state.selected[`${activeSchema.value}.${activeTable.value}`] = true;
      } else {
        state.selected[`${activeSchema.value}.${activeTable.value}`] =
          "indeterminate";
      }
      const schemaTables = Object.keys(md[activeSchema.value]);
      if (
        schemaTables.some((t) => state.selected[`${activeSchema.value}.${t}`])
      ) {
        state.selected[activeSchema.value] = "indeterminate";
      } else {
        state.selected[activeSchema.value] = false;
      }
    }
  } else {
    if (
      state.selected[`${activeSchema.value}.${activeTable.value}`] !== false
    ) {
      const tableColumns = md[activeSchema.value][activeTable.value];
      if (
        tableColumns.some(
          (c) =>
            state.selected[`${activeSchema.value}.${activeTable.value}.${c}`],
        )
      ) {
        state.selected[`${activeSchema.value}.${activeTable.value}`] =
          "indeterminate";
      } else {
        state.selected[`${activeSchema.value}.${activeTable.value}`] = false;
      }
      const schemaTables = Object.keys(md[activeSchema.value]);
      if (
        schemaTables.some((t) => state.selected[`${activeSchema.value}.${t}`])
      ) {
        state.selected[activeSchema.value] = "indeterminate";
      } else {
        state.selected[activeSchema.value] = false;
      }
    }
  }
}

const disabled = computed(() => {
  return (
    Object.entries(state.selected).filter(([, value]) => value === true)
      .length === 0
  );
});
</script>

<template>
  <UForm :schema="parser" :state="state" @submit="submit">
    <div class="w-full p-2">
      <div class="flex flex-col gap-2">
        <div class="flex h-48 w-full gap-2">
          <div class="flex max-h-96 flex-1 flex-col overflow-auto">
            <div
              v-for="schema of schemas"
              :key="schema"
              class="flex cursor-pointer items-center gap-2"
            >
              <UCheckbox
                v-model="state.selected[schema]"
                @change="selectSchema(schema)"
              />
              <div
                :class="[
                  'flex flex-auto items-center justify-between rounded px-2 py-1 transition-colors',
                  activeSchema === schema
                    ? 'bg-primary-500/50'
                    : 'bg-transparent',
                ]"
                @click="viewSchema(schema)"
              >
                <span>{{ schema }}</span>
                <UIcon
                  name="lucide:chevron-right"
                  :class="[
                    'size-5 transition-opacity',
                    activeSchema === schema ? 'opacity-100' : 'opacity-0',
                  ]"
                />
              </div>
            </div>
          </div>
          <USeparator orientation="vertical" />
          <div class="flex max-h-96 flex-1 flex-col overflow-auto">
            <div
              v-for="table of tables"
              :key="table"
              class="flex cursor-pointer items-center gap-2"
            >
              <UCheckbox
                v-model="state.selected[`${activeSchema}.${table}`]"
                @change="selectTable(table)"
              />
              <div
                :class="[
                  'flex flex-auto items-center justify-between rounded px-2 py-1 transition-colors',
                  activeTable === table
                    ? 'bg-primary-500/50'
                    : 'bg-transparent',
                ]"
                @click="viewTable(table)"
              >
                <span>{{ table }}</span>
                <UIcon
                  name="lucide:chevron-right"
                  :class="[
                    'size-5 transition-opacity',
                    activeTable === table ? 'opacity-100' : 'opacity-0',
                  ]"
                />
              </div>
            </div>
          </div>
          <USeparator orientation="vertical" />
          <div class="flex max-h-96 flex-1 flex-col overflow-auto">
            <div
              v-for="column of columns"
              :key="column"
              class="flex items-center gap-2"
            >
              <UCheckbox
                v-model="
                  state.selected[`${activeSchema}.${activeTable}.${column}`]
                "
                @change="selectColumn(column)"
              />
              <div
                class="flex flex-auto items-center justify-between rounded px-2 py-1 transition-colors"
              >
                <span>{{ column }}</span>
              </div>
            </div>
          </div>
        </div>
        <span class="text-2xl">Options</span>
        <USeparator />
        <div class="flex h-32 flex-row gap-2">
          <div class="flex flex-col gap-2">
            <UFormField label="Type">
              <USelect
                v-model="state.type"
                :items="types"
                :ui="{ base: 'w-36' }"
              />
            </UFormField>
          </div>
          <USeparator orientation="vertical" class="h-full" />
          <div class="flex flex-col gap-2">
            <UCheckbox
              v-model="state.schema_only"
              label="Export schema only"
              :disabled="state.data_only"
            />
            <UCheckbox
              v-model="state.data_only"
              label="Export data only"
              :disabled="state.schema_only"
            />
            <UCheckbox
              v-model="state.drop_schema"
              label="Drop schema before creating"
              :disabled="state.data_only"
            />
            <UCheckbox
              v-model="state.ignore_constraints"
              label="Ignore constraints"
              :disabled="state.schema_only"
            />
            <UCheckbox
              v-model="state.wrap_in_transaction"
              label="Wrap in transaction"
            />
          </div>
          <USeparator orientation="vertical" class="h-full" />
          <UFormField label="Drop/Create tables?" :disabled="state.data_only">
            <URadioGroup v-model="state.drop_table" :items="drop" />
          </UFormField>
        </div>
        <UButton
          icon="lucide:upload"
          type="submit"
          label="Export"
          :ui="{ base: 'self-center' }"
          :disabled="disabled"
        />
      </div>
    </div>
  </UForm>
</template>
