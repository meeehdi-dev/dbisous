import { createSharedComposable } from "@vueuse/shared";
import { ref } from "vue";

export const useApp = createSharedComposable(() => {
  // TODO: replace database with connection
  // const connection = ref<string>("");
  // TODO: add back database
  const database = ref<string>("");
  const schema = ref<string>("");
  const table = ref<string>("");

  return {
    // connection,
    database,
    schema,
    table,
  };
});
