<script setup lang="ts">
import { Effect } from "effect";
import * as v from "valibot";
import { reactive, ref } from "vue";
import { useWails } from "../../wails";
import { CreateConnection, SelectFile } from "../../../wailsjs/go/app/App";
import { FormSubmitEvent } from "@nuxt/ui/dist/module";

const emit = defineEmits<{ (e: "connectionAdded"): void }>();

const wails = useWails();

const schema = v.object({
  id: v.string(),
  name: v.string(),
  type: v.string(),
  connection_string: v.string(),
});
const parser = v.safeParser(schema);
type Schema = v.InferOutput<typeof schema>;

const state = reactive<Schema>({
  id: "",
  name: "",
  type: "",
  connection_string: "",
});

const types = ref([{ label: "Sqlite", value: "sqlite" }]);

function addConnection(event: FormSubmitEvent<Schema>) {
  Effect.runPromise(
    wails(() => CreateConnection(event.data)).pipe(
      Effect.tap(() => {
        emit("connectionAdded");
      }),
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
</script>

<template>
  <UForm
    :schema="parser"
    :state="state"
    class="space-y-4"
    @submit="addConnection"
  >
    <UFormField label="Type" name="type">
      <USelect v-model="state.type" :items="types" class="w-full" />
    </UFormField>

    <UFormField name="id" hidden>
      <UInput v-model="state.id" />
    </UFormField>

    <UFormField label="Name" name="name">
      <UInput placeholder="Optional name" v-model="state.name" class="w-full" />
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

    <div class="flex justify-end">
      <UButton type="submit" label="Submit" />
    </div>
  </UForm>
</template>
