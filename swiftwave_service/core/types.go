package core

// ************************************************************************************* //
//                                Swiftwave System Configuration 		   			     //
// ************************************************************************************* //

// UserRole : role of the registered user
type UserRole string

const (
	// AdministratorRole : admin user can perform any operation on the system
	AdministratorRole UserRole = "admin"
	// ManagerRole : manager user can perform any operation on the system
	// except user management, system configuration and server management
	ManagerRole UserRole = "manager"
)

// ServerStatus : status of the server
type ServerStatus string

const (
	ServerNeedsSetup ServerStatus = "needs_setup"
	ServerPreparing  ServerStatus = "preparing"
	ServerOnline     ServerStatus = "online"
	ServerOffline    ServerStatus = "offline"
)

// SwarmMode : mode of the swarm
type SwarmMode string

const (
	SwarmManager SwarmMode = "manager"
	SwarmWorker  SwarmMode = "worker"
)

// ProxyType : type of the proxy
type ProxyType string

const (
	BackupProxy ProxyType = "backup"
	ActiveProxy ProxyType = "active"
)

// ProxyConfig : hold information about proxy configuration
type ProxyConfig struct {
	Enabled bool      `json:"enabled" gorm:"default:false"`
	Type    ProxyType `json:"type" gorm:"default:'active'"`
}

// ************************************************************************************* //
//                                Application Level Table       		   			     //
// ************************************************************************************* //

// UpstreamType : type of source for the codebase of the application
type UpstreamType string

const (
	UpstreamTypeGit        UpstreamType = "git"
	UpstreamTypeSourceCode UpstreamType = "sourceCode"
	UpstreamTypeImage      UpstreamType = "image"
)

// GitProvider : type of git provider
type GitProvider string

const (
	GitProviderNone   GitProvider = "none"
	GitProviderGithub GitProvider = "github"
	GitProviderGitlab GitProvider = "gitlab"
)

// DomainSSLStatus : status of the ssl certificate for a domain
type DomainSSLStatus string

const (
	DomainSSLStatusNone    DomainSSLStatus = "none"
	DomainSSLStatusPending DomainSSLStatus = "pending"
	DomainSSLStatusFailed  DomainSSLStatus = "failed"
	DomainSSLStatusIssued  DomainSSLStatus = "issued"
)

// DeploymentStatus : status of the deployment
type DeploymentStatus string

const (
	DeploymentStatusPending       DeploymentStatus = "pending"
	DeploymentStatusDeployPending DeploymentStatus = "deployPending"
	DeploymentStatusLive          DeploymentStatus = "live"
	DeploymentStatusStopped       DeploymentStatus = "stopped"
	DeploymentStatusFailed        DeploymentStatus = "failed"
	DeploymentStalled             DeploymentStatus = "stalled"
)

// ProtocolType : type of protocol for ingress rule
type ProtocolType string

const (
	HTTPProtocol  ProtocolType = "http"
	HTTPSProtocol ProtocolType = "https"
	TCPProtocol   ProtocolType = "tcp"
	UDPProtocol   ProtocolType = "udp"
)

// IngressRuleStatus : status of the ingress rule
type IngressRuleStatus string

const (
	IngressRuleStatusPending  IngressRuleStatus = "pending"
	IngressRuleStatusApplied  IngressRuleStatus = "applied"
	IngressRuleStatusFailed   IngressRuleStatus = "failed"
	IngressRuleStatusDeleting IngressRuleStatus = "deleting"
)

// RedirectRuleStatus : status of the redirect rule
type RedirectRuleStatus string

const (
	RedirectRuleStatusPending  RedirectRuleStatus = "pending"
	RedirectRuleStatusApplied  RedirectRuleStatus = "applied"
	RedirectRuleStatusFailed   RedirectRuleStatus = "failed"
	RedirectRuleStatusDeleting RedirectRuleStatus = "deleting"
)

// DeploymentMode : mode of deployment of application (replicated or global)
type DeploymentMode string

const (
	DeploymentModeReplicated DeploymentMode = "replicated"
	DeploymentModeGlobal     DeploymentMode = "global"
)

// ApplicationUpdateResult : result of application update
type ApplicationUpdateResult struct {
	RebuildRequired bool
	ReloadRequired  bool
	DeploymentId    string
}

// DeploymentUpdateResult : result of deployment update
type DeploymentUpdateResult struct {
	RebuildRequired bool
	DeploymentId    string
}

// BackupType : type of backup
type BackupType string

const (
	LocalBackup BackupType = "local"
	S3Backup    BackupType = "s3"
)

// BackupStatus : status of backup
type BackupStatus string

const (
	BackupPending BackupStatus = "pending"
	BackupFailed  BackupStatus = "failed"
	BackupSuccess BackupStatus = "success"
)

// RestoreType : type of restore
type RestoreType string

const (
	LocalRestore RestoreType = "local"
)

// RestoreStatus : status of restore
type RestoreStatus string

const (
	RestorePending RestoreStatus = "pending"
	RestoreFailed  RestoreStatus = "failed"
	RestoreSuccess RestoreStatus = "success"
)

// PersistentVolumeType : type of persistent volume
type PersistentVolumeType string

const (
	PersistentVolumeTypeLocal PersistentVolumeType = "local"
	PersistentVolumeTypeNFS   PersistentVolumeType = "nfs"
)

// NFSConfig : configuration for NFS Storage
type NFSConfig struct {
	Host    string `json:"host,omitempty"`
	Path    string `json:"path,omitempty"`
	Version int    `json:"version,omitempty"`
}
