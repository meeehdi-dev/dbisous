import { createSharedComposable } from "@vueuse/core";
import { ref } from "vue";

export enum ChangeType {
  Insert = "INSERT",
  Update = "UPDATE",
  Delete = "DELETE",
}

interface ChangeId {
  id: number;
}

interface InsertChange extends ChangeId {
  type: ChangeType.Insert;
  table: string;
  values: Record<string, unknown>;
}

interface UpdateChange extends ChangeId {
  type: ChangeType.Update;
  table: string;
  values: Record<string, unknown>;
  primaryKey: string;
  rowKey: unknown;
}

interface DeleteChange extends ChangeId {
  type: ChangeType.Delete;
  table: string;
  primaryKey: string;
  rowKey: unknown;
}

type Change = InsertChange | UpdateChange | DeleteChange;

function toSqlValue(value: unknown): string {
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
  return `INSERT INTO ${change.table} (${Object.keys(change.values).join(", ")}) VALUES (${Object.values(
    change.values,
  )
    .map(toSqlValue)
    .join(", ")});`;
}
function formatUpdateChangeToSql(change: UpdateChange) {
  return `UPDATE ${change.table} SET ${Object.entries(change.values)
    .map(([key, value]) => `${key} = ${toSqlValue(value)}`)
    .join(", ")} WHERE ${change.primaryKey} = ${toSqlValue(change.rowKey)};`;
}
function formatDeleteChangeToSql(change: DeleteChange) {
  return `DELETE FROM ${change.table} WHERE ${change.primaryKey} = ${toSqlValue(change.rowKey)};`;
}

export const useTransaction = createSharedComposable(() => {
  const changeId = ref(0);
  const changes = ref<Array<Change>>([]);

  const abortListenerKey = ref(0);
  const abortListeners = ref<Record<number, () => void>>({});

  function commit() {
    const inserts: Array<InsertChange> = [];
    const updates: Array<UpdateChange> = [];
    const deletes: Array<DeleteChange> = [];

    changes.value.forEach((change) => {
      if (isInsertChange(change)) {
        inserts.push(change);
      } else if (isUpdateChange(change)) {
        updates.push(change);
      } else if (isDeleteChange(change)) {
        deletes.push(change);
      }
    });

    const insertsStr = inserts.map(formatInsertChangeToSql).join("\n");
    const updatesStr = updates.map(formatUpdateChangeToSql).join("\n");
    const deletesStr = deletes.map(formatDeleteChangeToSql).join("\n");

    const fullInsertStr =
      inserts.length > 0
        ? `-- ${inserts.length} insert${inserts.length > 1 ? "s" : ""}\n${insertsStr}\n`
        : "";
    const fullUpdateStr =
      updates.length > 0
        ? `-- ${updates.length} update${updates.length > 1 ? "s" : ""}\n${updatesStr}\n`
        : "";
    const fullDeleteStr =
      deletes.length > 0
        ? `-- ${deletes.length} delete${deletes.length > 1 ? "s" : ""}\n${deletesStr}\n`
        : "";

    const sql =
      "START TRANSACTION;\n" +
      fullInsertStr +
      fullUpdateStr +
      fullDeleteStr +
      "COMMIT;\n";

    return sql;
  }

  function addAbortListener(cb: () => void) {
    abortListeners.value[abortListenerKey.value] = cb;
    return abortListenerKey.value++;
  }
  function removeAbortListener(key: number) {
    delete abortListeners.value[key];
  }

  function abort() {
    changes.value = [];
    Object.values(abortListeners.value).forEach((cb) => cb());
  }

  function addUpdate(
    table: string,
    primaryKey: string,
    rowKey: unknown,
    key: string,
    value: unknown,
  ) {
    let update = changes.value.find(
      (c) =>
        isUpdateChange(c) &&
        c.table === table &&
        c.primaryKey === primaryKey &&
        c.rowKey === rowKey,
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
      changes.value.push(update);
    }
    update.values[key] = value;
  }

  function removeUpdate(
    table: string,
    primaryKey: string,
    rowKey: unknown,
    key: string,
  ) {
    const update = changes.value.find(
      (c) =>
        isUpdateChange(c) &&
        c.table === table &&
        c.primaryKey === primaryKey &&
        c.rowKey === rowKey,
    ) as UpdateChange | undefined;
    if (!update) {
      return;
    }
    delete update.values[key];
    if (Object.keys(update.values).length === 0) {
      changes.value = changes.value.filter((v) => v.id !== update.id);
    }
  }

  // TODO:
  // function addInsert() {}
  // function removeInsert() {}

  function toggleDelete(table: string, primaryKey: string, rowKey: unknown) {
    let delete_ = changes.value.find(
      (c) =>
        isDeleteChange(c) &&
        c.table === table &&
        c.primaryKey === primaryKey &&
        c.rowKey === rowKey,
    ) as DeleteChange | undefined;
    if (delete_) {
      // @ts-expect-error wtf ts?
      changes.value = changes.value.filter((v) => v.id !== delete_.id);
      return;
    }
    if (!delete_) {
      delete_ = {
        id: changeId.value++,
        type: ChangeType.Delete,
        table,
        primaryKey,
        rowKey,
      };
      changes.value.push(delete_);
    }
  }

  return {
    changes,
    commit,
    abort,
    addUpdate,
    removeUpdate,
    toggleDelete,
    addAbortListener,
    removeAbortListener,
  };
});

export function isInsertChange(change: Change): change is InsertChange {
  return change.type === ChangeType.Insert;
}

export function isUpdateChange(change: Change): change is UpdateChange {
  return change.type === ChangeType.Update;
}

export function isDeleteChange(change: Change): change is DeleteChange {
  return change.type === ChangeType.Delete;
}
