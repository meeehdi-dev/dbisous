import { createSharedComposable } from "@vueuse/core";
import type * as app from "_/dbisous/app/models.js";
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
