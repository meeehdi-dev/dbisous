import { Effect } from "effect";
import { client } from "../../wailsjs/go/models";
import { cell } from "../components/database/table";
import { TableColumn, TableData } from "@nuxt/ui/dist/module";

function formatColumns(columns: client.ColumnMetadata[], actions = true) {
  const formatted = columns.map(
    ({ name, type, default_value: defaultValue, nullable }) =>
      ({
        accessorKey: name,
        header: name,
        cell: cell({ type, defaultValue, nullable, disabled: !actions }),
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

export function formatQueryResult(result: client.QueryResult) {
  return Effect.succeed({
    ...result,
    columns: formatColumns(result.columns),
  });
}
