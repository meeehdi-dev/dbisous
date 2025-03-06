<script setup lang="ts">
import { useConnections } from "@/composables/useConnections";
import { useUrlParams } from "@/composables/useUrlParams";
import { FormSubmitEvent } from "@nuxt/ui/dist/module";
import * as v from "valibot";
import { computed, reactive, ref } from "vue";

enum ExportType {
  SQL = "sql",
}

const { databaseId } = useUrlParams();
const { metadata } = useConnections();

const formSchema = v.object({
  type: v.enum(ExportType),
});
const parser = v.safeParser(formSchema);
type FormSchema = v.InferOutput<typeof formSchema>;

const state = reactive<Partial<FormSchema>>({ type: ExportType.SQL });

const types = ref(
  Object.entries(ExportType).map(([label, value]) => ({ label, value })),
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

function submitConnection(event: FormSubmitEvent<FormSchema>) {
  console.log(event);
}

const activeSchema = ref("");
const activeTable = ref("");
const selected = ref<Record<string, true | false | "indeterminate">>({});

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
  if (selected.value[schema]) {
    const md = metadata.value[databaseId.value].columns;
    const tables = Object.keys(md[schema]);
    console.log(tables);
    tables.forEach((table) => {
      selected.value[`${schema}.${table}`] = true;
      console.log(`${schema}.${table}`);
      const columns = md[schema][table];
      columns.forEach((column) => {
        selected.value[`${schema}.${table}.${column}`] = true;
      });
    });
  } else {
    Object.keys(selected.value).forEach((key) => {
      if (key.startsWith(`${schema}.`)) {
        selected.value[key] = false;
      }
    });
  }
}

function selectTable(table: string) {
  const md = metadata.value[databaseId.value].columns;
  if (selected.value[`${activeSchema.value}.${table}`] === true) {
    if (selected.value[activeSchema.value] !== true) {
      const schemaTables = Object.keys(md[activeSchema.value]);
      if (
        schemaTables.every((t) => selected.value[`${activeSchema.value}.${t}`])
      ) {
        selected.value[activeSchema.value] = true;
      } else {
        selected.value[activeSchema.value] = "indeterminate";
      }
    }
    const columns = md[activeSchema.value][table];
    columns.forEach((column) => {
      selected.value[`${activeSchema.value}.${table}.${column}`] = true;
    });
  } else {
    if (selected.value[activeSchema.value] !== false) {
      const schemaTables = Object.keys(md[activeSchema.value]);
      if (
        schemaTables.some((t) => selected.value[`${activeSchema.value}.${t}`])
      ) {
        selected.value[activeSchema.value] = "indeterminate";
      } else {
        selected.value[activeSchema.value] = false;
      }
    }
    Object.keys(selected.value).forEach((key) => {
      if (key.startsWith(`${activeSchema.value}.${table}.`)) {
        selected.value[key] = false;
      }
    });
  }
}

function selectColumn(column: string) {
  const md = metadata.value[databaseId.value].columns;
  if (
    selected.value[`${activeSchema.value}.${activeTable.value}.${column}`] ===
    true
  ) {
    if (selected.value[`${activeSchema.value}.${activeTable.value}`] !== true) {
      const tableColumns = md[activeSchema.value][activeTable.value];
      if (
        tableColumns.every(
          (c) =>
            selected.value[`${activeSchema.value}.${activeTable.value}.${c}`],
        )
      ) {
        selected.value[`${activeSchema.value}.${activeTable.value}`] = true;
      } else {
        selected.value[`${activeSchema.value}.${activeTable.value}`] =
          "indeterminate";
      }
    }
  } else {
    if (
      selected.value[`${activeSchema.value}.${activeTable.value}`] !== false
    ) {
      const tableColumns = md[activeSchema.value][activeTable.value];
      if (
        tableColumns.some(
          (c) =>
            selected.value[`${activeSchema.value}.${activeTable.value}.${c}`],
        )
      ) {
        selected.value[`${activeSchema.value}.${activeTable.value}`] =
          "indeterminate";
      } else {
        selected.value[`${activeSchema.value}.${activeTable.value}`] = false;
      }
    }
  }
}
</script>

<template>
  <div class="p-2 w-full">
    <div class="flex gap-4 h-48 w-full">
      <div class="flex flex-1 flex-col max-h-96">
        <div
          v-for="schema of schemas"
          :key="schema"
          class="flex gap-2 items-center cursor-pointer"
        >
          <UCheckbox
            v-model="selected[schema]"
            @change="selectSchema(schema)"
          />
          <div
            :class="[
              'flex flex-auto px-2 py-1 rounded justify-between items-center transition-colors',
              activeSchema === schema ? 'bg-primary-500/50' : 'bg-transparent',
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
      <div class="flex flex-1 flex-col max-h-96">
        <div
          v-for="table of tables"
          :key="table"
          class="flex gap-2 items-center cursor-pointer"
        >
          <UCheckbox
            v-model="selected[`${activeSchema}.${table}`]"
            @change="selectTable(table)"
          />
          <div
            :class="[
              'flex flex-auto px-2 py-1 rounded justify-between items-center transition-colors',
              activeTable === table ? 'bg-primary-500/50' : 'bg-transparent',
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
      <div class="flex flex-1 flex-col max-h-96">
        <div
          v-for="column of columns"
          :key="column"
          class="flex gap-2 items-center"
        >
          <UCheckbox
            v-model="selected[`${activeSchema}.${activeTable}.${column}`]"
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
    <UForm :schema="parser" :state="state" @submit="submitConnection">
      <UFormField label="Type" name="type">
        <USelect v-model="state.type" :items="types" :ui="{ base: 'w-36' }" />
      </UFormField>
    </UForm>
  </div>
</template>
