import {
  ChangeType,
  DeleteChange,
  formatDeleteChangeToSql,
  formatInsertChangeToSql,
  formatUpdateChangeToSql,
  InsertChange,
  UpdateChange,
} from "@/utils/transaction";
import { createSharedComposable } from "@vueuse/core";
import { ref } from "vue";

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
        ? `-- ${insertChanges.value.length.toString()} insert${insertChanges.value.length > 1 ? "s" : ""}\n${insertsStr}\n`
        : "";
    const fullUpdateStr =
      updateChanges.value.length > 0
        ? `-- ${updateChanges.value.length.toString()} update${updateChanges.value.length > 1 ? "s" : ""}\n${updatesStr}\n`
        : "";
    const fullDeleteStr =
      deleteChanges.value.length > 0
        ? `-- ${deleteChanges.value.length.toString()} delete${deleteChanges.value.length > 1 ? "s" : ""}\n${deletesStr}\n`
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
    // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
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
