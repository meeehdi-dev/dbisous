<script setup lang="ts">
import * as v from "valibot";
import { reactive, ref } from "vue";
import { useWails } from "@/composables/useWails";
import { SelectFile } from "_/go/app/App";
import { FormSubmitEvent } from "@nuxt/ui/dist/module";
import { useConnections } from "@/composables/useConnections";
import { app } from "_/go/models";

const emit = defineEmits<{ connectionAdded: [] }>();
const connection = defineModel<app.Connection>();

const wails = useWails();
const { fetchConnections, addConnection, updateConnection } = useConnections();

const formSchema = v.object({
  id: v.optional(v.string()),
  created_at: v.optional(v.string()),
  updated_at: v.optional(v.string()),
  type: v.optional(v.enum(app.ConnectionType)),
  name: v.string(),
  connection_string: v.string(),
});
const parser = v.safeParser(formSchema);
type FormSchema = v.InferOutput<typeof formSchema>;

const state = reactive<FormSchema>(
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

async function submit(event: FormSubmitEvent<FormSchema>) {
  if (event.data.id) {
    await updateConnection(event.data as app.Connection);
  } else {
    await addConnection(event.data as app.Connection);
  }
  await fetchConnections();
  emit("connectionAdded");
}

async function selectFile() {
  const result = await wails(SelectFile);
  if (result instanceof Error) {
    return;
  }
  state.connection_string = result;
}

function selectType(type: app.ConnectionType) {
  state.type = type;
  active.value = 1;
}
</script>

<template>
  <UForm :schema="parser" :state="state" @submit="submit">
    <UStepper v-model="active" :items="items" disabled>
      <template #type>
        <AppConnectionTypeSelector :value="state.type" @select="selectType" />
      </template>

      <template #details>
        <div class="flex pb-4">
          <UButton
            v-if="!state.id"
            label="Back"
            color="neutral"
            variant="outline"
            icon="lucide:arrow-left"
            @click="active = 0"
          />
        </div>

        <UFormField label="Name">
          <UInput
            v-model="state.name"
            placeholder="Optional name"
            class="w-full"
          />
        </UFormField>

        <UFormField label="File">
          <UInput
            v-model="state.connection_string"
            placeholder="Select a file"
            class="w-full"
          >
            <template v-if="state.type === app.ConnectionType.SQLite" #trailing>
              <UButton
                variant="link"
                icon="lucide:upload"
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
