import { h } from "vue";
import AppCell from "./AppCell.vue";

export const cell =
  (type: string, nullable = false) =>
  (ctx: CellContext<unknown, unknown>) =>
    h(AppCell, {
      value: ctx.getValue(),
      type,
      nullable,
    });
