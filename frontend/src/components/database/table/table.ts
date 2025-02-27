import { h } from "vue";
import AppCell from "@/components/database/table/cell/AppCell.vue";
import type { TableColumn, TableData } from "@nuxt/ui/dist/module";
import { client } from "_/go/models";

export enum RowAction {
  View = "view",
  Insert = "insert",
  Duplicate = "duplicate",
  Delete = "delete",
}

type Row = Record<string, unknown>;

export type RowEmits = {
  view: [Row];
  insert: [];
  duplicate: [Row];
  delete: [Row];
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

export const booleanTypes = ["bool", "boolean"];
export const textTypes = ["name", "text", "varchar", "character varying"];
export const dateTypes = [
  "timestamp",
  "datetime",
  "date",
  "timestamp without time zone",
];
export const numberTypes = [
  "integer",
  "tinyint",
  "unsigned bigint",
  "int",
  "int4",
  "float8",
  "bigint",
  "double precision",
  "double",
];

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
