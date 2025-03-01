import { h } from "vue";
import AppCell from "@/components/database/table/cell/AppCell.vue";
import AppColumnHeader, {
  SortDirection,
} from "@/components/database/table/column/AppColumnHeader.vue";
import type { TableColumn, TableData } from "@nuxt/ui/dist/module";
import { client } from "_/go/models";

export enum RowAction {
  View = "view",
  Insert = "insert",
  Duplicate = "duplicate",
  Delete = "delete",
}

export interface RowEmits<T> {
  view: [T];
  insert: [];
  duplicate: [T];
  delete: [T];
  paginationChange: [number, number];
}

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
  // @ts-expect-error missing type
  (ctx: CellContext<unknown, unknown>) =>
    h(AppCell, {
      table,
      primaryKey,
      column,
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      row: ctx.row.original,
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-unsafe-call
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
  const formatted: TableColumn<TableData>[] = columns.map(
    ({ name, type, default_value: defaultValue, nullable }) => ({
      accessorKey: name,
      cell: cell({
        table,
        primaryKey,
        column: name,
        type,
        defaultValue,
        nullable,
        disabled,
      }),
      header: getHeader(name),
    }),
  );

  formatted.push({
    accessorKey: "action",
    header: "Actions",
  });

  // eslint-disable-next-line @typescript-eslint/no-unsafe-return
  return formatted;
}

export function getHeader(name: string) {
  return ({ column }: { column: unknown }) => {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error tkt
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-call
    const sort = column.getIsSorted();

    return h(AppColumnHeader, {
      label: name,
      sort,
      onClick: (s: SortDirection) => {
        if (s) {
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-expect-error tkt
          // eslint-disable-next-line @typescript-eslint/no-unsafe-call
          column.toggleSorting(s === "desc");
        } else {
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-expect-error tkt
          // eslint-disable-next-line @typescript-eslint/no-unsafe-call
          column.clearSorting();
        }
      },
    });
  };
}
