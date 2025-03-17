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
    options: [],
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
    options: [],
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
    options: [],
  });
});

test("parses connection string (postgres, port, database, options)", () => {
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
    options: [{ name: "sslmode", value: "disable" }],
  });
});

test("parses connection string (mysql)", () => {
  expect(parseConnectionString("mysql://root:mysql@tcp")).toStrictEqual({
    host: "tcp",
    port: "",
    user: "root",
    pass: "mysql",
    database: "",
    options: [],
  });
});

test("parses connection string (mysql, port)", () => {
  expect(parseConnectionString("mysql://root:mysql@tcp:5432")).toStrictEqual({
    host: "tcp",
    port: "5432",
    user: "root",
    pass: "mysql",
    database: "",
    options: [],
  });
});

test("parses connection string (mysql, port, database)", () => {
  expect(
    parseConnectionString("mysql://root:mysql@tcp:5432/mysql"),
  ).toStrictEqual({
    host: "tcp",
    port: "5432",
    user: "root",
    pass: "mysql",
    database: "mysql",
    options: [],
  });
});

test("parses connection string (mysql, port, database, options)", () => {
  expect(
    parseConnectionString("mysql://root:mysql@tcp:5432/mysql?sslmode=disable"),
  ).toStrictEqual({
    host: "tcp",
    port: "5432",
    user: "root",
    pass: "mysql",
    database: "mysql",
    options: [{ name: "sslmode", value: "disable" }],
  });
});
