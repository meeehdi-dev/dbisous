import { createSharedComposable } from "@vueuse/shared";
import { ref } from "vue";

export const useApp = createSharedComposable(() => {
  const connection = ref<string>("");
  const database = ref<string>("");
  const schema = ref<string>("");
  const table = ref<string>("");

  return {
    connection,
    database,
    schema,
    table,
  };
});
