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

function formatValue(v: unknown): string {
  return typeof v === "string"
    ? `'${v.replaceAll("'", "\\'")}'`
    : (v as string);
}

function formatInsertChangeToSql(change: InsertChange) {
  return `INSERT INTO ${change.table} (${Object.keys(change.values).join(", ")}) VALUES (${Object.values(
    change.values,
  )
    .map(formatValue)
    .join(", ")});`;
}
function formatUpdateChangeToSql(change: UpdateChange) {
  return `UPDATE ${change.table} SET ${Object.entries(change.values)
    .map(([key, value]) => `${key} = ${formatValue(value)}`)
    .join(", ")} WHERE ${change.primaryKey} = ${formatValue(change.rowKey)};`;
}
function formatDeleteChangeToSql(change: DeleteChange) {
  return `DELETE FROM ${change.table} WHERE ${change.primaryKey} = ${formatValue(change.rowKey)};`;
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

    const sql = fullInsertStr + fullUpdateStr + fullDeleteStr;

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
        c.type === ChangeType.Update &&
        c.table === table &&
        c.primaryKey === primaryKey &&
        c.rowKey === rowKey,
    );
    if (update && !isUpdateChange(update)) {
      return;
    }
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
        c.type === ChangeType.Update &&
        c.table === table &&
        c.primaryKey === primaryKey &&
        c.rowKey === rowKey,
    );
    if (!update || !isUpdateChange(update)) {
      return;
    }
    delete update.values[key];
    if (Object.keys(update.values).length === 0) {
      changes.value = changes.value.filter((v) => v.id !== update.id);
    }
  }

  return {
    changes,
    commit,
    abort,
    addUpdate,
    removeUpdate,
    addAbortListener,
    removeAbortListener,
  };
});

function isInsertChange(change: Change): change is InsertChange {
  return change.type === ChangeType.Insert;
}

function isUpdateChange(change: Change): change is UpdateChange {
  return change.type === ChangeType.Update;
}

function isDeleteChange(change: Change): change is DeleteChange {
  return change.type === ChangeType.Delete;
}
