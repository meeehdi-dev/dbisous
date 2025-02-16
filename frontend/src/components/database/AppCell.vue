<script setup lang="ts">
import { ref } from "vue";
import { booleanTypes, dateTypes, numberTypes, textTypes } from "./table";

const {
  value,
  type = "",
  nullable = false,
} = defineProps<{
  value: unknown;
  type?: string;
  nullable?: boolean;
}>();

const bool = ref(!!value);
const t = ref(value as string);
// const wails = useWails();
// const { currentDatabase } = useConnection();

const types = ref<Array<string>>([]);
// await Effect.runPromise(
//   wails(() => GetTypes(currentDatabase.value)).pipe(
//     Effect.tap((data) => {
//       types.value = data.rows.map((d) => d.type);
//     }),
//   ),
// );
// TODO: SEPARATE EACH COMPONENT + EMIT CHANGES
</script>

<template>
  <USelect
    v-if="type === 'TYPE'"
    variant="ghost"
    :items="types"
    v-model="t"
    class="w-full"
  />
  <UCheckbox v-else-if="booleanTypes.includes(type)" v-model="bool" />
  <UInput
    v-else-if="textTypes.includes(type)"
    variant="ghost"
    :value="value"
    class="w-full"
  />
  <UInput
    v-else-if="dateTypes.includes(type)"
    variant="ghost"
    :value="value"
    class="w-full"
  />
  <UInputNumber
    v-else-if="numberTypes.includes(type)"
    variant="ghost"
    :value="value"
    class="w-full"
  />
  <span v-else-if="type === ''" class="italic">{{ value }}</span>
  <span v-else class="font-bold text-red-400">{{ value }} ({{ type }})</span>
</template>
