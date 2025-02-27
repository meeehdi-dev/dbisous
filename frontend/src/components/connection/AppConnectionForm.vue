<script setup lang="ts">
import * as v from "valibot";
import { reactive, ref } from "vue";
import { useWails } from "@/composables/useWails";
import { CreateConnection, SelectFile, UpdateConnection } from "_/go/app/App";
import { FormSubmitEvent } from "@nuxt/ui/dist/module";
import { useConnections } from "@/composables/useConnections";
import { app } from "_/go/models";

const emit = defineEmits<{ connectionAdded: [] }>();
const connection = defineModel<app.Connection>();

const wails = useWails();
const { fetchConnections } = useConnections();

const schema = v.object({
  id: v.string(),
  created_at: v.string(),
  updated_at: v.string(),
  name: v.string(),
  type: v.enum(app.ConnectionType),
  connection_string: v.string(),
});
const parser = v.safeParser(schema);
type Schema = v.InferOutput<typeof schema>;

const state = reactive<Partial<Schema>>(
  connection.value ?? {
    name: "",
    connection_string: "",
  },
);

const items = [
  {
    title: "Database type",
    icon: "lucide:database",
    slot: "type",
  },
  {
    title: "Connection details",
    icon: "lucide:link",
    slot: "details",
  },
];
const active = ref(state.id ? 1 : 0);

async function submitConnection(event: FormSubmitEvent<Schema>) {
  const result = await wails(() =>
    event.data.id ? UpdateConnection(event.data) : CreateConnection(event.data),
  );
  if (result instanceof Error) {
    // TODO: specific error handling
  } else {
    fetchConnections();
    emit("connectionAdded");
  }
}

async function selectFile() {
  const result = await wails(SelectFile);
  if (result instanceof Error) {
    // TODO: specific error handling
  } else {
    state.connection_string = result;
  }
}

function selectType(type: app.ConnectionType) {
  state.type = type;
  active.value = 1;
}
</script>

<template>
  <UForm :schema="parser" :state="state" @submit="submitConnection">
    <UStepper :items="items" v-model="active" disabled>
      <template #type>
        <AppConnectionTypeSelector :value="state.type" @select="selectType" />
      </template>

      <template #details>
        <div class="flex pb-4">
          <UButton
            label="Back"
            color="neutral"
            variant="outline"
            icon="lucide:arrow-left"
            @click="active = 0"
            v-if="!state.id"
          />
        </div>

        <UFormField label="Name" name="name">
          <UInput
            placeholder="Optional name"
            v-model="state.name"
            class="w-full"
          />
        </UFormField>

        <UFormField label="File" name="url">
          <UInput
            placeholder="Select a file"
            v-model="state.connection_string"
            class="w-full"
          >
            <template #trailing v-if="state.type === app.ConnectionType.SQLite">
              <UButton
                variant="link"
                icon="lucide:upload"
                aria-label="Select SQLite file"
                @click="selectFile"
              />
            </template>
          </UInput>
        </UFormField>

        <div class="flex justify-end pt-4">
          <UButton type="submit" icon="lucide:save" label="Save" />
        </div>
      </template>
    </UStepper>
  </UForm>
</template>
