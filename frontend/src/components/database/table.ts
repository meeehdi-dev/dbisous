import { h } from "vue";
import AppCell from "./AppCell.vue";
import type { TableColumn, TableData } from "@nuxt/ui/dist/module";
import { client } from "../../../wailsjs/go/models";

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

export const cell =
  (type?: string, nullable?: boolean) => (ctx: CellContext<unknown, unknown>) =>
    h(AppCell, {
      value: ctx.getValue(),
      type,
      nullable,
    });

export function formatColumns(
  columns: client.ColumnMetadata[],
  actions = true,
) {
  const formatted = columns.map(
    (column) =>
      ({
        accessorKey: column.name,
        header: column.name,
        cell: cell(column.type, column.nullable),
      }) as TableColumn<TableData>,
  );

  if (actions) {
    formatted.push({
      accessorKey: "action",
      header: "Actions",
    });
  }

  return formatted;
}
