const postgresPrefix = "postgres://";
const mysqlPrefix = "mysql://";

interface Connection {
  host: string;
  port: string;
  user: string;
  pass: string;
  database: string;
  options: Array<{ name: string; value: string }>;
}

export function parseConnectionString(connectionString: string): Connection {
  if (connectionString.startsWith(postgresPrefix)) {
    connectionString = connectionString.slice(postgresPrefix.length);
  }
  if (connectionString.startsWith(mysqlPrefix)) {
    connectionString = connectionString.slice(mysqlPrefix.length);
  }
  const [userInfo, connectionInfo] = connectionString.split("@");
  const [user, pass] = userInfo.split(":");
  const connectionParts = (connectionInfo || "").split("/");
  let hostInfo = connectionParts[0];
  const params = connectionParts[1];
  if (hostInfo.startsWith("tcp")) {
    const matches = hostInfo.match(/tcp\((.*)\)/);
    if (matches && matches.length > 0) {
      hostInfo = matches[1];
    }
  }
  const [host, port] = hostInfo.split(":");
  const [database, options] = (params || "").split("?");

  const connectionHost = host || "";
  const connectionPort = port || "";
  const connectionUser = user || "";
  const connectionPass = pass || "";
  const connectionDatabase = database || "";
  const connectionOptions = options
    ? options.split("&").map((option) => {
        const [name, value] = option.split("=");
        return { name, value };
      })
    : [];

  return {
    host: connectionHost,
    port: connectionPort,
    user: connectionUser,
    pass: connectionPass,
    database: connectionDatabase,
    options: connectionOptions,
  };
}
