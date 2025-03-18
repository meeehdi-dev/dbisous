<script setup lang="ts">
import * as v from "valibot";
import { computed, reactive, ref, watch } from "vue";
import { useWails } from "@/composables/useWails";
import { SelectFile, TestConnection } from "_/go/app/App";
import { useConnections } from "@/composables/shared/useConnections";
import { app } from "_/go/models";
import type { FormSubmitEvent } from "@nuxt/ui";
import { parseConnectionString } from "@/utils/connection";

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

const connectionHost = ref("");
const connectionPort = ref("");
const connectionUser = ref("");
const connectionPass = ref("");
const connectionDatabase = ref("");
const connectionOptions = ref<Array<{ name: string; value: string }>>([]);

const postgresPrefix = "postgres://";
function onConnectionStringChange() {
  const { host, port, user, pass, database, options } = parseConnectionString(
    state.connection_string,
  );

  connectionHost.value = host;
  connectionPort.value = port;
  connectionUser.value = user;
  connectionPass.value = pass;
  connectionDatabase.value = database;
  connectionOptions.value = options;
}
onConnectionStringChange();
function onConnectionInfoChange() {
  state.connection_string = `${state.type === app.ConnectionType.PostgreSQL ? postgresPrefix : ""}${connectionUser.value}:${connectionPass.value}@${connectionHost.value}${connectionPort.value ? `:${connectionPort.value}` : ""}/${connectionDatabase.value}${connectionOptions.value.length > 0 ? "?" : ""}${connectionOptions.value.map((option) => [option.name, option.value].join(option.value ? "=" : "")).join("&")}`;
}

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
const active = ref(state.type ? 1 : 0);
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

const placeholders = {
  [app.ConnectionType.MySQL]: {
    connectionString: "mysql:mysql@tcp/mysql",
    port: "3306",
    database: "mysql",
  },
  [app.ConnectionType.PostgreSQL]: {
    connectionString: "postgres://postgres:postgres@localhost:5432/postgres",
    port: "5432",
    database: "postgres",
  },
};

function onConnectionOptionAdd() {
  connectionOptions.value.push({ name: "", value: "" });
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
        <div class="flex flex-col gap-2">
          <UButton
            v-if="!state.id"
            label="Back"
            color="neutral"
            variant="outline"
            icon="lucide:arrow-left"
            :ui="{ base: 'p-2 self-start' }"
            @click="active = 0"
          />

          <UFormField label="Name">
            <UInput
              v-model="state.name"
              placeholder="Optional name"
              class="w-full"
            />
          </UFormField>
        </div>

        <UFormField
          v-if="state.type === app.ConnectionType.SQLite"
          label="File"
        >
          <UInput
            v-model="state.connection_string"
            placeholder="Select a file"
            class="w-full"
          >
            <template #trailing>
              <UButton
                variant="link"
                icon="lucide:upload"
                @click="selectFile"
              />
            </template>
          </UInput>
        </UFormField>

        <template
          v-if="
            state.type &&
            (state.type === app.ConnectionType.PostgreSQL ||
              state.type === app.ConnectionType.MySQL)
          "
        >
          <UFormField label="Connection string">
            <UInput
              v-model="state.connection_string"
              :placeholder="placeholders[state.type].connectionString"
              class="w-full"
              @update:model-value="onConnectionStringChange"
            />
          </UFormField>
          <UFormField label="Host">
            <UInput
              v-model="connectionHost"
              placeholder="localhost"
              class="w-full"
              @update:model-value="onConnectionInfoChange"
            />
          </UFormField>
          <UFormField label="Port">
            <UInput
              v-model="connectionPort"
              :placeholder="placeholders[state.type].port"
              class="w-full"
              @update:model-value="onConnectionInfoChange"
            />
          </UFormField>
          <UFormField label="User">
            <UInput
              v-model="connectionUser"
              placeholder="user"
              class="w-full"
              @update:model-value="onConnectionInfoChange"
            />
          </UFormField>
          <UFormField label="Password">
            <UInput
              v-model="connectionPass"
              placeholder="pass"
              class="w-full"
              @update:model-value="onConnectionInfoChange"
            />
          </UFormField>
          <UFormField label="Database">
            <UInput
              v-model="connectionDatabase"
              :placeholder="placeholders[state.type].database"
              class="w-full"
              @update:model-value="onConnectionInfoChange"
            />
          </UFormField>

          <UFormField label="Options">
            <div class="flex gap-2">
              <div class="flex flex-auto">
                <USeparator orientation="vertical" />
              </div>
              <div class="flex flex-col gap-2">
                <div
                  v-for="(connectionOption, i) in connectionOptions"
                  :key="i"
                  class="flex gap-2"
                >
                  <UInput
                    v-model="connectionOption.name"
                    placeholder="Name"
                    @update:model-value="onConnectionInfoChange"
                  />
                  <UInput
                    v-model="connectionOption.value"
                    placeholder="Value"
                    @update:model-value="onConnectionInfoChange"
                  />
                </div>
                <UButton
                  label="Add option"
                  icon="lucide:plus"
                  variant="soft"
                  :ui="{ base: 'self-start' }"
                  @click="onConnectionOptionAdd"
                />
              </div>
            </div>
          </UFormField>
        </template>

        <div class="flex justify-end gap-2 pt-4">
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
