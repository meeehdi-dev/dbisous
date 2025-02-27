import { createSharedComposable } from "@vueuse/core";
import { ref } from "vue";

export enum ChangeType {
  Insert = "INSERT",
  Update = "UPDATE",
  Delete = "DELETE",
}

interface IChange {
  id: number;
  type: ChangeType;
}

export interface InsertChange extends IChange {
  type: ChangeType.Insert;
  table: string;
  values: Record<string, unknown>;
  __key: number; // used as temporary row key
}

export interface UpdateChange extends IChange {
  type: ChangeType.Update;
  table: string;
  values: Record<string, unknown>;
  primaryKey: string;
  rowKey: unknown;
}

export interface DeleteChange extends IChange {
  type: ChangeType.Delete;
  table: string;
  primaryKey: string;
  rowKey: unknown;
}

function toSqlValue(value: unknown | undefined): string {
  if (value === null) {
    return "NULL";
  }

  switch (typeof value) {
    case "boolean":
      return value ? "TRUE" : "FALSE";
    case "number":
      return value.toString();
    case "string":
      return `'${value.replace(/'/g, "''")}'`; // Escape single quotes
    case "object":
      if (Array.isArray(value)) {
        return `(${value.map(toSqlValue).join(", ")})`;
      } else if (value instanceof Date) {
        return `'${value.toISOString()}'`;
      } else {
        throw new Error(`Unsupported object type: ${typeof value}`);
      }
    default:
      throw new Error(`Unsupported data type: ${typeof value}`);
  }
}

function formatInsertChangeToSql(change: InsertChange) {
  // NOTE: filter out __key and NULL values
  const values = Object.fromEntries(
    Object.entries(change.values).filter(
      ([key, value]) => key !== "__key" && value !== "NULL",
    ),
  );
  return `INSERT OR ROLLBACK INTO ${change.table} (${Object.keys(values).join(
    ", ",
  )}) VALUES (${Object.entries(values)
    .map(([, value]) => toSqlValue(value))
    .join(", ")});`;
}
function formatUpdateChangeToSql(change: UpdateChange) {
  return `UPDATE OR ROLLBACK ${change.table} SET ${Object.entries(change.values)
    .map(([key, value]) => `${key} = ${toSqlValue(value)}`)
    .join(", ")} WHERE ${change.primaryKey} = ${toSqlValue(change.rowKey)};`;
}
function formatDeleteChangeToSql(change: DeleteChange) {
  return `DELETE FROM ${change.table} WHERE ${change.primaryKey} = ${toSqlValue(change.rowKey)};`;
}

export const useTransaction = createSharedComposable(() => {
  const changeId = ref(0);
  const insertChanges = ref<Array<InsertChange>>([]);
  const updateChanges = ref<Array<UpdateChange>>([]);
  const deleteChanges = ref<Array<DeleteChange>>([]);

  function commit() {
    const insertsStr = insertChanges.value
      .filter(
        (c) =>
          !deleteChanges.value.find(
            (d) => d.table === c.table && d.rowKey === c.id,
          ),
      )
      .map(formatInsertChangeToSql)
      .join("\n");
    const updatesStr = updateChanges.value
      .filter(
        (c) =>
          !deleteChanges.value.find(
            (d) => d.table === c.table && d.rowKey === c.rowKey,
          ),
      )
      .map(formatUpdateChangeToSql)
      .join("\n");
    const deletesStr = deleteChanges.value
      .map(formatDeleteChangeToSql)
      .join("\n");

    const fullInsertStr =
      insertChanges.value.length > 0
        ? `-- ${insertChanges.value.length} insert${insertChanges.value.length > 1 ? "s" : ""}\n${insertsStr}\n`
        : "";
    const fullUpdateStr =
      updateChanges.value.length > 0
        ? `-- ${updateChanges.value.length} update${updateChanges.value.length > 1 ? "s" : ""}\n${updatesStr}\n`
        : "";
    const fullDeleteStr =
      deleteChanges.value.length > 0
        ? `-- ${deleteChanges.value.length} delete${deleteChanges.value.length > 1 ? "s" : ""}\n${deletesStr}\n`
        : "";

    const sql =
      "BEGIN;\n" + fullInsertStr + fullUpdateStr + fullDeleteStr + "COMMIT;\n";

    return sql;
  }

  function abort() {
    insertChanges.value = [];
    updateChanges.value = [];
    deleteChanges.value = [];
  }

  function addUpdate(
    table: string,
    primaryKey: string,
    rowKey: unknown,
    key: string,
    value: unknown,
  ) {
    let update = updateChanges.value.find(
      (c) =>
        c.table === table && c.primaryKey === primaryKey && c.rowKey === rowKey,
    ) as UpdateChange | undefined;
    if (!update) {
      update = {
        id: changeId.value++,
        type: ChangeType.Update,
        table,
        primaryKey,
        rowKey,
        values: {},
      };
      updateChanges.value.push(update);
    }
    update.values[key] = value;
  }

  function removeUpdate(
    table: string,
    primaryKey: string,
    rowKey: unknown,
    key: string,
  ) {
    const update = updateChanges.value.find(
      (c) =>
        c.table === table && c.primaryKey === primaryKey && c.rowKey === rowKey,
    ) as UpdateChange;
    delete update.values[key];
    if (Object.keys(update.values).length === 0) {
      updateChanges.value.splice(
        updateChanges.value.findIndex((c) => c.id === update.id),
        1,
      );
    }
  }

  function addInsert(table: string, values: Record<string, unknown>) {
    const id = changeId.value++;
    insertChanges.value.push({
      id,
      type: ChangeType.Insert,
      table,
      values,
      __key: id,
    });
    return id;
  }

  function updateInsert(
    table: string,
    key: number,
    column: string,
    value: unknown,
  ) {
    const insert = insertChanges.value.find(
      (change) => change.table === table && change.__key === key,
    );
    if (!insert) {
      return;
    }
    insert.values[column] = value;
  }

  function removeInsert(table: string, key: number) {
    const insert = insertChanges.value.find(
      (c) => c.table === table && c.__key === key,
    );
    if (!insert) {
      return;
    }
    insertChanges.value.splice(
      insertChanges.value.findIndex((c) => c.id === insert.id),
      1,
    );
  }

  function toggleDelete(table: string, primaryKey: string, rowKey: unknown) {
    let delete_ = deleteChanges.value.find(
      (c) =>
        c.table === table && c.primaryKey === primaryKey && c.rowKey === rowKey,
    );
    if (!delete_) {
      delete_ = {
        id: changeId.value++,
        type: ChangeType.Delete,
        table,
        primaryKey,
        rowKey,
      };
      deleteChanges.value.push(delete_);
      return;
    }
    deleteChanges.value.splice(
      deleteChanges.value.findIndex((c) => c.id === delete_.id),
      1,
    );
  }

  return {
    insertChanges,
    updateChanges,
    deleteChanges,
    commit,
    abort,
    addUpdate,
    removeUpdate,
    addInsert,
    updateInsert,
    removeInsert,
    toggleDelete,
  };
});
