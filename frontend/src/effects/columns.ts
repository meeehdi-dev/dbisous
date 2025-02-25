import { Effect } from "effect";
import { client } from "_/go/models";
import { cell } from "@/components/database/table/table";
import { TableColumn, TableData } from "@nuxt/ui/dist/module";

function formatColumns(
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

export function formatQueryResult(
  result: client.QueryResult,
  disabled?: boolean,
  table?: string,
  primaryKey?: string,
) {
  return Effect.succeed({
    ...result,
    columns: formatColumns(result.columns, table, primaryKey, disabled),
  });
}
