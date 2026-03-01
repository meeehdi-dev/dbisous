import { expect, test } from "vitest";
import { parseConnectionString } from "./connection";

test("parses connection string (postgres)", () => {
  expect(
    parseConnectionString("postgres://postgres:postgres@localhost"),
  ).toStrictEqual({
    host: "localhost",
    port: "",
    user: "postgres",
    pass: "postgres",
    database: "",
    opts: [],
  });
});

test("parses connection string (postgres, port)", () => {
  expect(
    parseConnectionString("postgres://postgres:postgres@localhost:5432"),
  ).toStrictEqual({
    host: "localhost",
    port: "5432",
    user: "postgres",
    pass: "postgres",
    database: "",
    opts: [],
  });
});

test("parses connection string (postgres, port, database)", () => {
  expect(
    parseConnectionString(
      "postgres://postgres:postgres@localhost:5432/postgres",
    ),
  ).toStrictEqual({
    host: "localhost",
    port: "5432",
    user: "postgres",
    pass: "postgres",
    database: "postgres",
    opts: [],
  });
});

test("parses connection string (postgres, port, database, opts)", () => {
  expect(
    parseConnectionString(
      "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
    ),
  ).toStrictEqual({
    host: "localhost",
    port: "5432",
    user: "postgres",
    pass: "postgres",
    database: "postgres",
    opts: [{ name: "sslmode", value: "disable" }],
  });
});

test("parses connection string (mysql)", () => {
  expect(parseConnectionString("mysql://root:mysql@/")).toStrictEqual({
    host: "",
    port: "",
    user: "root",
    pass: "mysql",
    database: "",
    opts: [],
  });
});

test("parses connection string (mysql, host, port)", () => {
  expect(
    parseConnectionString("mysql://root:mysql@tcp(localhost:3306)"),
  ).toStrictEqual({
    host: "localhost",
    port: "3306",
    user: "root",
    pass: "mysql",
    database: "",
    opts: [],
  });
});

test("parses connection string (mysql, host, port, database)", () => {
  expect(
    parseConnectionString("mysql://root:mysql@tcp(localhost:3306)/mysql"),
  ).toStrictEqual({
    host: "localhost",
    port: "3306",
    user: "root",
    pass: "mysql",
    database: "mysql",
    opts: [],
  });
});

test("parses connection string (mysql, host, port, database, opts)", () => {
  expect(
    parseConnectionString(
      "mysql://root:mysql@tcp(localhost:3306)/mysql?sslmode=disable",
    ),
  ).toStrictEqual({
    host: "localhost",
    port: "3306",
    user: "root",
    pass: "mysql",
    database: "mysql",
    opts: [{ name: "sslmode", value: "disable" }],
  });
});
