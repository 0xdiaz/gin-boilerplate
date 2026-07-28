package config

import (
	"fmt"
)

type DatabaseConfiguration struct {
	Driver   string
	Dbname   string
	Username string
	Password string
	Host     string
	Port     string
	LogMode  bool
	SslMode  string // master SSL mode, e.g. disable, require
	// Replica (read-only); if empty, master values are used
	ReplicaDbname   string
	ReplicaUsername string
	ReplicaPassword string
	ReplicaHost     string
	ReplicaPort     string
	ReplicaSslMode  string
}

// DbConfiguration returns master and replica DSN from loaded config.
// Call after SetupConfig(). Returns empty strings if config not loaded.
func DbConfiguration() (masterDSN, replicaDSN string) {
	c := Get()
	if c == nil {
		return "", ""
	}
	d := &c.Database
	masterDSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		d.Host, d.Username, d.Password, d.Dbname, d.Port, d.SslMode,
	)
	replicaDSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		d.ReplicaHost, d.ReplicaUsername, d.ReplicaPassword, d.ReplicaDbname, d.ReplicaPort, d.ReplicaSslMode,
	)
	return masterDSN, replicaDSN
}
