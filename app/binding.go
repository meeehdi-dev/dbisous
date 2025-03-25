package app

import "dbisous/app/client"

func (a *App) GetConnections() ([]Connection, error) {
	return getConnections(metadataDB)
}

func (a *App) CreateConnection(connection Connection) error {
	return createConnection(metadataDB, connection)
}

func (a *App) UpdateConnection(connection Connection) error {
	return updateConnection(metadataDB, connection)
}

func (a *App) DeleteConnection(id string) error {
	return deleteConnection(metadataDB, id)
}

func (a *App) Connect(id string) (client.DatabaseMetadata, error) {
	return connect(activeConnections, metadataDB, id)
}

func (a *App) Disconnect(id string) error {
	return disconnect(activeConnections, id)
}

func (a *App) TestConnection(dbType ConnectionType, connectionString string) error {
	return testConnection(dbType, connectionString)
}
