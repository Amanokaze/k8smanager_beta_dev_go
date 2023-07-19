package config

import "fmt"

type DBConfig struct { //[db, primary]
	connectionHost     string //= 127.0.0.1
	connectionPort     uint16 //= 5432
	connectionName     string //= xxxxx
	connectionUser     string //= xxxxx
	connectionPassword string //= xxxxx
	connectionSslmode  string //= disable
	databaseOwner      string //= xxxxx
	databaseDocker     bool   //= false
	databaseDockerpath string //= xxxxx
	maxConnection      uint16 //= 100
	//[Namespace', 'Node', 'Pod', 'Service', 'Ingress', 'Deployment', 'StatefulSet', 'DaemonSet', 'ReplicaSet', 'Event', 'PersistentVolumeClaim', 'StorageClass', 'PersistentVolume' ]
}

// GetDBHost is a function that returns the database host.
func (d *DBConfig) GetDBConnString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Seoul",
		d.connectionHost, d.connectionUser, d.connectionPassword, d.connectionName, d.connectionPort)
}

// GetDsn is a function that returns the database dsn.
func (d *DBConfig) GetDsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", d.connectionUser, d.connectionPassword, d.connectionHost, d.connectionPort, d.connectionName, d.connectionSslmode)
}

// GetDBHost is a function that returns the database host.
func (d *DBConfig) GetHost() string {
	return d.connectionHost
}

// GetDBPort is a function that returns the database port.
func (d *DBConfig) GetMaxConnection() uint16 {
	return d.maxConnection
}

// GetOSUser is a function that returns the database os user.
func (d *DBConfig) GetDBOSUser() string {
	return d.databaseOwner
}

// GetDBDocker is a function that returns the database docker.
func (d *DBConfig) GetDBDocker() bool {
	return d.databaseDocker
}

// GetDBDockerpath is a function that returns the database docker path.
func (d *DBConfig) GetDBDockerpath() string {
	return d.databaseDockerpath
}
