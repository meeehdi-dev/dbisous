# DBisous

**Cross-platform database manager built with Wails (Golang) and Vue.js for viewing, editing, and querying databases.**

<p align="center">
<img width="1136" alt="Screenshot 2025-02-17 at 00 12 56" src="https://github.com/user-attachments/assets/0afdf0a8-b70d-4fbe-9f38-68901f166695" />
</p>

## Features

- **Database Connection Management**: Connect to SQLite, MySQL, PostgreSQL, and more to come.
- **Table Viewer**: List tables, view structures, and paginate data.
- **Query Execution**: Execute SQL queries with syntax highlighting and save frequently used queries.
- **Data Editing**: Add, update, and delete rows directly from the viewer.
- **Data Export**: Export table data to SQL, CSV, JSON, and other formats.
- **Data Import**: Import table data from SQL, CSV, JSON, and other formats.
- **Cross-Platform Support**: Runs on Windows, macOS, and Linux.

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/meeehdi-dev/dbisous.git
   cd dbisous
   ```

2. **Install dependencies**:
   ```bash
   make tidy
   cd frontend && npm install
   ```

3. **Build and run**:
   ```bash
   make build
   ./build/bin/dbisous # (.exe on win, .app on macos)
   ```

## Usage

1. **Connect to a Database**: Set up a new connection profile and connect to your database.
2. **View Tables**: Browse and view the structure and data of your tables.
3. **Execute Queries**: Use the SQL query editor to run and save queries.
4. **Edit Data**: Directly edit table data from the viewer.
5. **Export Data**: Export table data to various formats.
5. **Import Data**: Import table data from various formats.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
