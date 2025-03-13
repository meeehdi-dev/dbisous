export enum ChangeType {
  Insert = "INSERT",
  Update = "UPDATE",
  Delete = "DELETE",
}

interface Change {
  id: number;
  type: ChangeType;
}

export interface InsertChange extends Change {
  type: ChangeType.Insert;
  table: string;
  values: Record<string, unknown>;
  __key: number; // used as temporary row key
}

export interface UpdateChange extends Change {
  type: ChangeType.Update;
  table: string;
  values: Record<string, unknown>;
  primaryKey: string;
  rowKey: unknown;
}

export interface DeleteChange extends Change {
  type: ChangeType.Delete;
  table: string;
  primaryKey: string;
  rowKey: unknown;
}

export function toSqlValue(value: unknown): string {
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

export function formatInsertChangeToSql(change: InsertChange) {
  // NOTE: filter out __key and NULL values
  const values = Object.fromEntries(
    Object.entries(change.values).filter(
      ([key, value]) => key !== "__key" && value !== "NULL",
    ),
  );
  return `INSERT INTO ${change.table} (${Object.keys(values)
    .map((column) => `"${column}"`)
    .join(", ")}) VALUES (${Object.entries(values)
    .map(([, value]) => toSqlValue(value))
    .join(", ")});`;
}
export function formatUpdateChangeToSql(change: UpdateChange) {
  return `UPDATE ${change.table} SET ${Object.entries(change.values)
    .map(([key, value]) => `"${key}" = ${toSqlValue(value)}`)
    .join(", ")} WHERE ${change.primaryKey} = ${toSqlValue(change.rowKey)};`;
}
export function formatDeleteChangeToSql(change: DeleteChange) {
  return `DELETE FROM ${change.table} WHERE ${change.primaryKey} = ${toSqlValue(change.rowKey)};`;
}
