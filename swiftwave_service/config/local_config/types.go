package local_config

import _ "embed"

//go:embed .version
var softwareVersion string

type Config struct {
	IsDevelopmentMode bool             `yaml:"dev_mode"`
	Version           string           `yaml:"-"`
	ServiceConfig     ServiceConfig    `yaml:"service"`
	PostgresqlConfig  PostgresqlConfig `yaml:"postgresql"`
}

type ServiceConfig struct {
	UseTLS                      bool   `yaml:"use_tls"`
	ManagementNodeAddress       string `yaml:"management_node_address"`
	BindAddress                 string `yaml:"bind_address"`
	BindPort                    int    `yaml:"bind_port"`
	SocketPathDirectory         string `yaml:"-"`
	DataDirectory               string `yaml:"-"`
	LocalPostgresDataDirectory  string `yaml:"-"`
	TarballDirectoryPath        string `yaml:"-"`
	PVBackupDirectoryPath       string `yaml:"-"`
	PVRestoreDirectoryPath      string `yaml:"-"`
	NetworkName                 string `yaml:"-"`
	HAProxyServiceName          string `yaml:"-"`
	HAProxyUnixSocketDirectory  string `yaml:"-"`
	HAProxyUnixSocketPath       string `yaml:"-"`
	HAProxyDataDirectoryPath    string `yaml:"-"`
	UDPProxyServiceName         string `yaml:"-"`
	UDPProxyUnixSocketDirectory string `yaml:"-"`
	UDPProxyUnixSocketPath      string `yaml:"-"`
	UDPProxyDataDirectoryPath   string `yaml:"-"`
	SSLCertDirectoryPath        string `yaml:"-"`
	LogDirectoryPath            string `yaml:"-"`
	InfoLogFilePath             string `yaml:"-"`
	ErrorLogFilePath            string `yaml:"-"`
}

type PostgresqlConfig struct {
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	User             string `yaml:"user"`
	Password         string `yaml:"password"`
	Database         string `yaml:"database"`
	TimeZone         string `yaml:"time_zone"`
	SSLMode          string `yaml:"ssl_mode"`
	RunLocalPostgres bool   `yaml:"run_local_postgres"`
}
