# dbisous

**Cross-platform database manager built with Wails (Golang) and Vue.js for viewing, editing, and querying databases.**

## Features

- **Database Connection Management**: Connect to SQLite, MySQL, PostgreSQL, and more.
- **Table Viewer**: List tables, view structures, and paginate data.
- **Query Execution**: Execute SQL queries with syntax highlighting and save frequently used queries.
- **Data Editing**: Add, update, and delete rows directly from the viewer.
- **Data Export**: Export table data to CSV, JSON, and other formats.
- **User Management**: Manage database users and roles.
- **Performance Monitoring**: View query performance metrics.
- **Security**: Encrypt sensitive data and implement role-based access control (RBAC).
- **Cross-Platform Support**: Runs on Windows, macOS, and Linux.

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/meeehdi-dev/dbisous.git
   cd dbisous
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   cd frontend && npm install
   ```

3. **Build and run**:
   ```bash
   wails build
   ./build/bin/dbisous # (.exe on win, .app on macos)
   ```

## Usage

1. **Connect to a Database**: Set up a new connection profile and connect to your database.
2. **View Tables**: Browse and view the structure and data of your tables.
3. **Execute Queries**: Use the SQL query editor to run and save queries.
4. **Edit Data**: Directly edit table data from the viewer.
5. **Export Data**: Export table data to various formats.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See [LICENSE.md](LICENSE.md) for details.
