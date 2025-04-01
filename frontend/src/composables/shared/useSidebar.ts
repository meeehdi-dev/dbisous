import { createSharedComposable } from "@vueuse/core";
import { app } from "_/go/models";
import { ref } from "vue";

export const useSidebar = createSharedComposable(() => {
  const slideoverOpen = ref(false);
  const editedConnection = ref<Omit<app.Connection, "id"> & { id?: string }>();

  return { slideoverOpen, editedConnection };
});
