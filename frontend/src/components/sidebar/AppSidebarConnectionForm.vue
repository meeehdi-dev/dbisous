<script setup lang="ts">
import * as v from "valibot";
import { computed, reactive, ref, watch } from "vue";
import { useWails } from "@/composables/useWails";
import { SelectFile, TestConnection } from "_/go/app/App";
import { useConnections } from "@/composables/shared/useConnections";
import { app } from "_/go/models";
import type { FormSubmitEvent } from "@nuxt/ui";

const emit = defineEmits<{ connectionAdded: [] }>();
const connection = defineModel<app.Connection>();

const wails = useWails();
// eslint-disable-next-line no-undef
const toast = useToast();
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
const status = ref<"idle" | "loading" | "success" | "failed">("idle");
watch([() => state.type, () => state.connection_string], () => {
  status.value = "idle";
});
const statusColor = computed(() => {
  switch (status.value) {
    case "loading":
      return "warning";
    case "success":
      return "success";
    case "failed":
      return "error";
    case "idle":
    default:
      return "neutral";
  }
});

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

async function testConnection() {
  const { type, connection_string } = state;
  if (!type || !connection_string) {
    return;
  }
  status.value = "loading";
  const result = await wails(() => TestConnection(type, connection_string));
  if (result instanceof Error) {
    status.value = "failed";
  } else {
    status.value = "success";
    toast.add({
      title: "Test connection",
      description: "Successfully pinged database!",
    });
  }
}
</script>

<template>
  <UForm :schema="parser" :state="state" @submit="submit">
    <UStepper v-model="active" :items="items" disabled>
      <template #type>
        <AppSidebarConnectionTypeSelector
          :value="state.type"
          @select="selectType"
        />
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

        <div class="flex justify-end gap-4 pt-4">
          <UButton
            icon="lucide:plug-zap"
            label="Test"
            variant="soft"
            :loading="status === 'loading'"
            :disabled="
              status !== 'idle' || !state.type || !state.connection_string
            "
            @click="testConnection()"
          >
            <template #trailing>
              <UChip class="ml-1" standalone inset :color="statusColor" />
            </template>
          </UButton>
          <UButton type="submit" icon="lucide:save" label="Save" />
        </div>
      </template>
    </UStepper>
  </UForm>
</template>
