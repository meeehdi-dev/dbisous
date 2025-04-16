import { createSharedComposable } from "@vueuse/core";
import { app } from "_/go/models";
import { ref } from "vue";

export const useSidebar = createSharedComposable(() => {
  const slideoverOpen = ref(false);
  const editedConnection = ref<
    Omit<app.Connection, "id" | "created_at" | "updated_at"> & {
      id?: string;
      created_at?: string;
      updated_at?: string;
    }
  >();

  return { slideoverOpen, editedConnection };
});
