import { h } from "vue";
import AppCell from "./AppCell.vue";
import type { TableColumn, TableData } from "@nuxt/ui/dist/module";
import { client } from "../../../../wailsjs/go/models";

export enum RowAction {
  View = "view",
  Copy = "copy",
  Remove = "remove",
}

export type RowEmits = {
  view: [TableData];
  copy: [TableData];
  remove: [TableData];
  paginationChange: [number, number];
};

export type FormattedQueryResult = Omit<
  client.QueryResult,
  "convertValues" | "columns"
> & {
  columns: Array<TableColumn<TableData>>;
};

export type CellProps = {
  type?: string;
  defaultValue?: unknown;
  nullable?: boolean;
  disabled: boolean;
};
export const cell =
  ({ type, defaultValue, nullable, disabled }: CellProps) =>
  (ctx: CellContext<unknown, unknown>) =>
    h(AppCell, {
      initialValue: ctx.getValue(),
      type,
      defaultValue,
      nullable,
      disabled,
    });

export const booleanTypes = ["BOOL", "BOOLEAN", "TINYINT"];
export const textTypes = ["NAME", "TEXT", "VARCHAR"];
export const dateTypes = ["TIMESTAMP", "DATETIME", "DATE"];
export const numberTypes = ["UNSIGNED BIGINT", "INT", "INT4", "FLOAT8"];
