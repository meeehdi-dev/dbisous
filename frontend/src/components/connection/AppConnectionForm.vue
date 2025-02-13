<script setup lang="ts">
import { Effect } from "effect";
import * as v from "valibot";
import { reactive, ref } from "vue";
import { useWails } from "../../wails";
import { CreateConnection, SelectFile, UpdateConnection } from "../../../wailsjs/go/app/App";
import { FormSubmitEvent } from "@nuxt/ui/dist/module";

import { useConnections } from "../../composables/useConnections";
import { app } from "../../../wailsjs/go/models";

const emit = defineEmits<{ connectionAdded: [] }>();
const connection = defineModel<app.Connection>();

const wails = useWails();
const { fetchConnections } = useConnections();

const schema = v.object({
  id: v.string(),
  created_at: v.string(),
  updated_at: v.string(),
  name: v.string(),
  type: v.string(),
  connection_string: v.string(),
});
const parser = v.safeParser(schema);
type Schema = v.InferOutput<typeof schema>;

const state = reactive<Schema>(
  connection.value ?? {
    id: "",
    created_at: "",
    updated_at: "",
    name: "",
    type: "",
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

function submitConnection(event: FormSubmitEvent<Schema>) {
  Effect.runPromise(
    wails(() =>
      event.data.id
        ? UpdateConnection(event.data)
        : CreateConnection(event.data),
    ).pipe(
      Effect.tap(() => {
        emit("connectionAdded");
      }),
      Effect.tap(fetchConnections),
    ),
  );
}

function selectFile() {
  Effect.runPromise(
    wails(SelectFile).pipe(
      Effect.tap((url) => {
        state.connection_string = url;
      }),
    ),
  );
}

function selectType(type: string) {
  state.type = type;
  active.value = 1;
}
</script>

<template>
  <UForm :schema="parser" :state="state" @submit="submitConnection">
    <UStepper :items="items" v-model="active" disabled>
      <template #type>
        <div class="flex justify-center">
          <div class="flex flex-col gap-4">
            <UButton
              icon="simple-icons:sqlite"
              label="SQLite"
              size="xl"
              :ui="{ label: 'w-full text-center' }"
              @click="selectType('sqlite')"
              :trailing-icon="
                state.type === 'sqlite' ? 'lucide:check' : undefined
              "
            />
            <UButton
              icon="simple-icons:postgresql"
              label="PostgreSQL"
              size="xl"
              :ui="{ label: 'w-full text-center' }"
              @click="selectType('mysql')"
              :trailing-icon="
                state.type === 'mysql' ? 'lucide:check' : undefined
              "
              disabled
            />
            <UButton
              icon="simple-icons:mysql"
              label="MySQL / MariaDB"
              size="xl"
              :ui="{ label: 'w-full text-center' }"
              @click="selectType('postgresql')"
              :trailing-icon="
                state.type === 'postgresql' ? 'lucide:check' : undefined
              "
              disabled
            />
          </div>
        </div>
      </template>

      <template #details>
        <div class="flex pb-4">
          <UButton
            label="Back"
            color="neutral"
            variant="outline"
            icon="lucide:arrow-left"
            @click="active = 0"
          />
        </div>

        <UFormField name="id" hidden>
          <UInput v-model="state.id" />
        </UFormField>

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
            <template #trailing>
              <UButton
                variant="link"
                icon="lucide:upload"
                aria-label="Upload file"
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
