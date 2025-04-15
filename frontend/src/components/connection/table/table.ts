import { h } from "vue";
import AppCell from "@/components/connection/table/cell/AppCell.vue";
import AppColumnHeader, {
  SortDirection,
} from "@/components/connection/table/column/AppColumnHeader.vue";
import type { TableColumn } from "@nuxt/ui";
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
  "rows" | "columns" | "convertValues"
> & {
  rows: Array<Record<string, unknown>>;
  columns: Array<TableColumn<Record<string, unknown>>>;
};

export type CellProps = {
  table?: string;
  primaryKey?: string;
  column?: string;
  row?: unknown;
  type?: string;
  items?: string[];
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
    items,
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
      items,
      defaultValue,
      nullable,
      disabled,
    });

export const booleanTypes = ["bool", "boolean"];
export const textTypes = [
  "name",
  "text",
  "varchar",
  "character varying",
  "longtext",
];
export const enumTypes = ["user-defined", "enum"];
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
  onSort: (name: string, sort: SortDirection) => void | Promise<void>,
  onFilter: (name: string, filter: false | string) => void | Promise<void>,
  table?: string,
  enums?: client.EnumMetadata[],
  primaryKey?: string,
  disabled = false,
) {
  const formatted: TableColumn<Record<string, unknown>>[] = columns.map(
    ({ name, type, default_value: defaultValue, nullable }) => ({
      accessorKey: name,
      cell: cell({
        table,
        primaryKey,
        column: name,
        type,
        items: enums?.find((e) => e.column === name)?.values,
        defaultValue,
        nullable,
        disabled,
      }),
      header: getHeader(name, {
        onSort: async (s) => {
          await onSort(name, s);
        },
        onFilter: async (f) => {
          await onFilter(name, f);
        },
      }),
    }),
  );

  formatted.push({
    accessorKey: "action",
    header: "Actions",
  });

  // eslint-disable-next-line @typescript-eslint/no-unsafe-return
  return formatted;
}

export function getHeader(
  name: string,
  callbacks?: {
    onSort: (s: SortDirection) => Promise<void>;
    onFilter: (s: false | string) => Promise<void>;
  },
) {
  return ({
    column,
  }: {
    column: {
      getIsSorted: () => false | "asc" | "desc";
      getFilterValue: () => unknown;
    };
  }) => {
    const sort = column.getIsSorted();
    const filter = column.getFilterValue();

    return h(AppColumnHeader, {
      label: name,
      sort: sort ? (sort.toUpperCase() as SortDirection) : false,
      filter,
      onSort: callbacks?.onSort,
      onFilter: callbacks?.onFilter,
    });
  };
}
