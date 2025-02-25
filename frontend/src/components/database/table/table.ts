import { h } from "vue";
import AppCell from "@/components/database/table/cell/AppCell.vue";
import type { TableColumn, TableData } from "@nuxt/ui/dist/module";
import { client } from "_/go/models";

export enum RowAction {
  View = "view",
  Duplicate = "duplicate",
  Delete = "delete",
}

export type RowEmits = {
  view: [TableData];
  duplicate: [TableData];
  delete: [TableData];
  paginationChange: [number, number];
};

export type FormattedQueryResult = Omit<
  client.QueryResult,
  "convertValues" | "columns"
> & {
  columns: Array<TableColumn<TableData>>;
};

export type CellProps = {
  table?: string;
  primaryKey?: string;
  column?: string;
  row?: unknown;
  type?: string;
  defaultValue?: unknown;
  nullable?: boolean;
  disabled: boolean;
};
export const cell =
  ({
    table,
    primaryKey,
    column,
    type,
    defaultValue,
    nullable,
    disabled,
  }: CellProps) =>
  (ctx: CellContext<unknown, unknown>) =>
    h(AppCell, {
      table,
      primaryKey,
      column,
      row: ctx.row.original,
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

export function formatColumns(
  columns: client.ColumnMetadata[],
  table?: string,
  primaryKey?: string,
  disabled = false,
) {
  const formatted = columns.map(
    ({ name, type, default_value: defaultValue, nullable }) =>
      ({
        accessorKey: name,
        header: name,
        cell: cell({
          table,
          primaryKey,
          column: name,
          type,
          defaultValue,
          nullable,
          disabled,
        }),
      }) as TableColumn<TableData>,
  );

  formatted.push({
    accessorKey: "action",
    header: "Actions",
  });

  return formatted;
}
