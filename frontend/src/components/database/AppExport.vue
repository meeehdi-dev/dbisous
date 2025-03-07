<script setup lang="ts">
import { useConnections } from "@/composables/useConnections";
import { useUrlParams } from "@/composables/useUrlParams";
import { useWails } from "@/composables/useWails";
import { FormSubmitEvent } from "@nuxt/ui/dist/module";
import { Export } from "_/go/app/App";
import { client } from "_/go/models";
import * as v from "valibot";
import { computed, reactive, ref } from "vue";

const { databaseId } = useUrlParams();
const { metadata } = useConnections();
const wails = useWails();
// eslint-disable-next-line no-undef
const toast = useToast();

const exportSchema = v.object({
  type: v.enum(client.ExportType),
  drop_schema: v.enum(client.ExportDrop),
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
  drop_schema: client.ExportDrop.Do_nothing,
  drop_table: client.ExportDrop.Drop_and_create,
  selected: {},
});

const types = ref(
  Object.entries(client.ExportType).map(([label, value]) => ({ label, value })),
);
const drop = ref(
  Object.entries(client.ExportDrop).map(([label, value]) => ({
    label: label.replaceAll("_", " "),
    value,
  })),
);

const schemas = computed(() => {
  const md = metadata.value[databaseId.value].columns;
  return Object.keys(md);
});

const tables = computed(() => {
  if (activeSchema.value === "") {
    return [];
  }
  const md = metadata.value[databaseId.value].columns;
  return Object.keys(md[activeSchema.value]);
});

const columns = computed(() => {
  if (activeSchema.value === "" || activeTable.value === "") {
    return [];
  }
  const md = metadata.value[databaseId.value].columns;
  return md[activeSchema.value][activeTable.value];
});

async function submitConnection(event: FormSubmitEvent<ExportSchema>) {
  const result = await wails(() =>
    Export(databaseId.value, {
      ...event.data,
      selected: Object.entries(state.selected)
        .filter(([, value]) => value === true)
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
    const md = metadata.value[databaseId.value].columns;
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
  } else {
    Object.keys(state.selected).forEach((key) => {
      if (key.startsWith(`${schema}.`)) {
        state.selected[key] = false;
      }
    });
  }
}

function selectTable(table: string) {
  const md = metadata.value[databaseId.value].columns;
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
  const md = metadata.value[databaseId.value].columns;
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
  <UForm :schema="parser" :state="state" @submit="submitConnection">
    <div class="p-4 w-full">
      <div class="flex flex-col gap-4">
        <div class="flex gap-4 h-48 w-full">
          <div class="flex flex-1 flex-col max-h-96 overflow-auto">
            <div
              v-for="schema of schemas"
              :key="schema"
              class="flex gap-2 items-center cursor-pointer"
            >
              <UCheckbox
                v-model="state.selected[schema]"
                @change="selectSchema(schema)"
              />
              <div
                :class="[
                  'flex flex-auto px-2 py-1 rounded justify-between items-center transition-colors',
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
                    'transition-opacity size-5',
                    activeSchema === schema ? 'opacity-100' : 'opacity-0',
                  ]"
                />
              </div>
            </div>
          </div>
          <USeparator orientation="vertical" />
          <div class="flex flex-1 flex-col max-h-96 overflow-auto">
            <div
              v-for="table of tables"
              :key="table"
              class="flex gap-2 items-center cursor-pointer"
            >
              <UCheckbox
                v-model="state.selected[`${activeSchema}.${table}`]"
                @change="selectTable(table)"
              />
              <div
                :class="[
                  'flex flex-auto px-2 py-1 rounded justify-between items-center transition-colors',
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
                    'transition-opacity size-5',
                    activeTable === table ? 'opacity-100' : 'opacity-0',
                  ]"
                />
              </div>
            </div>
          </div>
          <USeparator orientation="vertical" />
          <div class="flex flex-1 flex-col max-h-96 overflow-auto">
            <div
              v-for="column of columns"
              :key="column"
              class="flex gap-2 items-center"
            >
              <UCheckbox
                v-model="
                  state.selected[`${activeSchema}.${activeTable}.${column}`]
                "
                @change="selectColumn(column)"
              />
              <div
                class="flex flex-auto px-2 py-1 rounded justify-between items-center transition-colors"
              >
                <span>{{ column }}</span>
              </div>
            </div>
          </div>
        </div>
        <UFormField label="Type" name="type">
          <USelect v-model="state.type" :items="types" :ui="{ base: 'w-36' }" />
        </UFormField>
        <div
          :class="[
            'flex flew-row gap-4 h-20 transition-opacity',
            state.type === 'sql' ? 'opacity-100' : 'opacity-50',
          ]"
        >
          <UFormField label="Drop/Create schemas?" name="drop">
            <URadioGroup
              v-model="state.drop_schema"
              :items="drop"
              :disabled="state.type !== 'sql'"
            />
          </UFormField>
          <USeparator orientation="vertical" />
          <UFormField label="Drop/Create tables?" name="drop">
            <URadioGroup
              v-model="state.drop_table"
              :items="drop"
              :disabled="state.type !== 'sql'"
            />
          </UFormField>
          <USeparator orientation="vertical" />
          <UFormField
            label="Drop constraints during import?"
            name="constraints"
          >
          </UFormField>
          <USeparator orientation="vertical" />
          <UFormField label="Wrap in transaction?" name="transaction">
          </UFormField>
          <USeparator orientation="vertical" />
          <div
            v-if="state.type !== 'sql'"
            class="flex items-center gap-2 text-warning-400"
          >
            <UIcon name="lucide:triangle-alert" />
            SQL Only
          </div>
        </div>
        <div class="flex justify-center">
          <UButton
            icon="lucide:upload"
            type="submit"
            label="Export"
            :disabled="disabled"
          />
        </div>
      </div>
    </div>
  </UForm>
</template>
