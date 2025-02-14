import type { TableData } from "@nuxt/ui/dist/module";

export enum RowAction {
  View = "view",
  Copy = "copy",
  Remove = "remove",
}

export type Emits = {
  view: [TableData];
  copy: [TableData];
  remove: [TableData];
};
