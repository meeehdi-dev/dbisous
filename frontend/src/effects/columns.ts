import { Effect } from "effect";
import { client } from "_/go/models";
import { cell } from "@/components/database/table/table";
import { TableColumn, TableData } from "@nuxt/ui/dist/module";

function formatColumns(
  table: string,
  primaryKey: string,
  columns: client.ColumnMetadata[],
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
          disabled: false,
        }),
      }) as TableColumn<TableData>,
  );

  formatted.push({
    accessorKey: "action",
    header: "Actions",
  });

  return formatted;
}

export function formatQueryResult(result: client.QueryResult) {
  return Effect.succeed({
    ...result,
    columns: formatColumns(result.table, result.primary_key, result.columns),
  });
}
