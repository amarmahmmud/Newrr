// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Application struct {
	ID                       string                     `json:"id"`
	Name                     string                     `json:"name"`
	EnvironmentVariables     []*EnvironmentVariable     `json:"environmentVariables"`
	PersistentVolumeBindings []*PersistentVolumeBinding `json:"persistentVolumeBindings"`
	LatestDeployment         *Deployment                `json:"latestDeployment"`
	Deployments              []*Deployment              `json:"deployments"`
	DeploymentMode           DeploymentMode             `json:"deploymentMode"`
	Replicas                 uint                       `json:"replicas"`
	IngressRules             []*IngressRule             `json:"ingressRules"`
}

type ApplicationInput struct {
	Name                         string                          `json:"name"`
	EnvironmentVariables         []*EnvironmentVariableInput     `json:"environmentVariables"`
	PersistentVolumeBindings     []*PersistentVolumeBindingInput `json:"persistentVolumeBindings"`
	Dockerfile                   *string                         `json:"dockerfile,omitempty"`
	BuildArgs                    []*BuildArgInput                `json:"buildArgs"`
	DeploymentMode               DeploymentMode                  `json:"deploymentMode"`
	Replicas                     *uint                           `json:"replicas,omitempty"`
	UpstreamType                 UpstreamType                    `json:"upstreamType"`
	GitCredentialID              *uint                           `json:"gitCredentialID,omitempty"`
	GitProvider                  *GitProvider                    `json:"gitProvider,omitempty"`
	RepositoryOwner              *string                         `json:"repositoryOwner,omitempty"`
	RepositoryName               *string                         `json:"repositoryName,omitempty"`
	RepositoryBranch             *string                         `json:"repositoryBranch,omitempty"`
	CommitHash                   *string                         `json:"commitHash,omitempty"`
	SourceCodeCompressedFileName *string                         `json:"sourceCodeCompressedFileName,omitempty"`
	DockerImage                  *string                         `json:"dockerImage,omitempty"`
	ImageRegistryCredentialID    *uint                           `json:"imageRegistryCredentialID,omitempty"`
}

type BuildArg struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type BuildArgInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CustomSSLInput struct {
	FullChain  string `json:"fullChain"`
	PrivateKey string `json:"privateKey"`
	SslIssuer  string `json:"sslIssuer"`
}

type Deployment struct {
	ID                           string                   `json:"id"`
	ApplicationID                string                   `json:"applicationID"`
	Application                  *Application             `json:"application"`
	UpstreamType                 UpstreamType             `json:"upstreamType"`
	GitCredentialID              uint                     `json:"gitCredentialID"`
	GitCredential                *GitCredential           `json:"gitCredential"`
	GitProvider                  GitProvider              `json:"gitProvider"`
	RepositoryOwner              string                   `json:"repositoryOwner"`
	RepositoryName               string                   `json:"repositoryName"`
	RepositoryBranch             string                   `json:"repositoryBranch"`
	CommitHash                   string                   `json:"commitHash"`
	SourceCodeCompressedFileName string                   `json:"sourceCodeCompressedFileName"`
	DockerImage                  string                   `json:"dockerImage"`
	ImageRegistryCredentialID    uint                     `json:"imageRegistryCredentialID"`
	ImageRegistryCredential      *ImageRegistryCredential `json:"imageRegistryCredential"`
	BuildArgs                    []*BuildArg              `json:"buildArgs"`
	Dockerfile                   string                   `json:"dockerfile"`
	DeploymentLogs               []*DeploymentLog         `json:"deploymentLogs"`
	Status                       DeploymentStatus         `json:"status"`
	CreatedAt                    time.Time                `json:"createdAt"`
}

type DeploymentLog struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type Domain struct {
	ID            uint            `json:"id"`
	Name          string          `json:"name"`
	SslStatus     DomainSSLStatus `json:"sslStatus"`
	SslFullChain  string          `json:"sslFullChain"`
	SslPrivateKey string          `json:"sslPrivateKey"`
	SslIssuedAt   string          `json:"sslIssuedAt"`
	SslIssuer     string          `json:"sslIssuer"`
	SslAutoRenew  bool            `json:"sslAutoRenew"`
	IngressRules  []*IngressRule  `json:"ingressRules"`
	RedirectRules []*RedirectRule `json:"redirectRules"`
}

type DomainInput struct {
	Name string `json:"name"`
}

type EnvironmentVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EnvironmentVariableInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GitCredential struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	Deployments []*Deployment `json:"deployments"`
}

type GitCredentialInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type GitCredentialRepositoryAccessInput struct {
	GitCredentialID  uint   `json:"gitCredentialId"`
	RepositoryURL    string `json:"repositoryUrl"`
	RepositoryBranch string `json:"repositoryBranch"`
}

type GitCredentialRepositoryAccessResult struct {
	GitCredentialID  uint           `json:"gitCredentialId"`
	GitCredential    *GitCredential `json:"gitCredential"`
	RepositoryURL    string         `json:"repositoryUrl"`
	RepositoryBranch string         `json:"repositoryBranch"`
	Success          bool           `json:"success"`
	Error            string         `json:"error"`
}

type ImageRegistryCredential struct {
	ID          uint          `json:"id"`
	URL         string        `json:"url"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	Deployments []*Deployment `json:"deployments"`
}

type ImageRegistryCredentialInput struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type IngressRule struct {
	ID          uint              `json:"id"`
	Domain      *Domain           `json:"domain"`
	Protocol    ProtocolType      `json:"protocol"`
	Port        uint              `json:"port"`
	Application *Application      `json:"application"`
	TargetPort  uint              `json:"targetPort"`
	Status      IngressRuleStatus `json:"status"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type IngressRuleInput struct {
	DomainID      uint         `json:"domainId"`
	ApplicationID string       `json:"applicationId"`
	Protocol      ProtocolType `json:"protocol"`
	Port          uint         `json:"port"`
	TargetPort    uint         `json:"targetPort"`
}

type PersistentVolume struct {
	ID                       uint                       `json:"id"`
	Name                     string                     `json:"name"`
	PersistentVolumeBindings []*PersistentVolumeBinding `json:"persistentVolumeBindings"`
}

type PersistentVolumeBinding struct {
	ID                 uint              `json:"id"`
	PersistentVolumeID uint              `json:"persistentVolumeID"`
	PersistentVolume   *PersistentVolume `json:"persistentVolume"`
	ApplicationID      string            `json:"applicationID"`
	Application        *Application      `json:"application"`
	MountingPath       string            `json:"mountingPath"`
}

type PersistentVolumeBindingInput struct {
	PersistentVolumeID uint   `json:"persistentVolumeID"`
	MountingPath       string `json:"mountingPath"`
}

type PersistentVolumeInput struct {
	Name string `json:"name"`
}

type RedirectRule struct {
	ID          uint               `json:"id"`
	Domain      *Domain            `json:"domain"`
	Protocol    ProtocolType       `json:"protocol"`
	Port        uint               `json:"port"`
	RedirectURL string             `json:"redirectURL"`
	Status      RedirectRuleStatus `json:"status"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

type RedirectRuleInput struct {
	DomainID    uint         `json:"domainId"`
	Protocol    ProtocolType `json:"protocol"`
	Port        uint         `json:"port"`
	RedirectURL string       `json:"redirectURL"`
}

type DeploymentMode string

const (
	DeploymentModeReplicated DeploymentMode = "replicated"
	DeploymentModeGlobal     DeploymentMode = "global"
)

var AllDeploymentMode = []DeploymentMode{
	DeploymentModeReplicated,
	DeploymentModeGlobal,
}

func (e DeploymentMode) IsValid() bool {
	switch e {
	case DeploymentModeReplicated, DeploymentModeGlobal:
		return true
	}
	return false
}

func (e DeploymentMode) String() string {
	return string(e)
}

func (e *DeploymentMode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DeploymentMode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DeploymentMode", str)
	}
	return nil
}

func (e DeploymentMode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type DeploymentStatus string

const (
	DeploymentStatusPending   DeploymentStatus = "pending"
	DeploymentStatusQueued    DeploymentStatus = "queued"
	DeploymentStatusDeploying DeploymentStatus = "deploying"
	DeploymentStatusRunning   DeploymentStatus = "running"
	DeploymentStatusStopped   DeploymentStatus = "stopped"
	DeploymentStatusFailed    DeploymentStatus = "failed"
)

var AllDeploymentStatus = []DeploymentStatus{
	DeploymentStatusPending,
	DeploymentStatusQueued,
	DeploymentStatusDeploying,
	DeploymentStatusRunning,
	DeploymentStatusStopped,
	DeploymentStatusFailed,
}

func (e DeploymentStatus) IsValid() bool {
	switch e {
	case DeploymentStatusPending, DeploymentStatusQueued, DeploymentStatusDeploying, DeploymentStatusRunning, DeploymentStatusStopped, DeploymentStatusFailed:
		return true
	}
	return false
}

func (e DeploymentStatus) String() string {
	return string(e)
}

func (e *DeploymentStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DeploymentStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DeploymentStatus", str)
	}
	return nil
}

func (e DeploymentStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type DomainSSLStatus string

const (
	DomainSSLStatusNone    DomainSSLStatus = "none"
	DomainSSLStatusPending DomainSSLStatus = "pending"
	DomainSSLStatusIssued  DomainSSLStatus = "issued"
)

var AllDomainSSLStatus = []DomainSSLStatus{
	DomainSSLStatusNone,
	DomainSSLStatusPending,
	DomainSSLStatusIssued,
}

func (e DomainSSLStatus) IsValid() bool {
	switch e {
	case DomainSSLStatusNone, DomainSSLStatusPending, DomainSSLStatusIssued:
		return true
	}
	return false
}

func (e DomainSSLStatus) String() string {
	return string(e)
}

func (e *DomainSSLStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DomainSSLStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DomainSSLStatus", str)
	}
	return nil
}

func (e DomainSSLStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type GitProvider string

const (
	GitProviderNone   GitProvider = "none"
	GitProviderGithub GitProvider = "github"
	GitProviderGitlab GitProvider = "gitlab"
)

var AllGitProvider = []GitProvider{
	GitProviderNone,
	GitProviderGithub,
	GitProviderGitlab,
}

func (e GitProvider) IsValid() bool {
	switch e {
	case GitProviderNone, GitProviderGithub, GitProviderGitlab:
		return true
	}
	return false
}

func (e GitProvider) String() string {
	return string(e)
}

func (e *GitProvider) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GitProvider(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GitProvider", str)
	}
	return nil
}

func (e GitProvider) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type IngressRuleStatus string

const (
	IngressRuleStatusPending  IngressRuleStatus = "pending"
	IngressRuleStatusApplied  IngressRuleStatus = "applied"
	IngressRuleStatusDeleting IngressRuleStatus = "deleting"
	IngressRuleStatusFailed   IngressRuleStatus = "failed"
)

var AllIngressRuleStatus = []IngressRuleStatus{
	IngressRuleStatusPending,
	IngressRuleStatusApplied,
	IngressRuleStatusDeleting,
	IngressRuleStatusFailed,
}

func (e IngressRuleStatus) IsValid() bool {
	switch e {
	case IngressRuleStatusPending, IngressRuleStatusApplied, IngressRuleStatusDeleting, IngressRuleStatusFailed:
		return true
	}
	return false
}

func (e IngressRuleStatus) String() string {
	return string(e)
}

func (e *IngressRuleStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IngressRuleStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IngressRuleStatus", str)
	}
	return nil
}

func (e IngressRuleStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProtocolType string

const (
	ProtocolTypeHTTP  ProtocolType = "http"
	ProtocolTypeHTTPS ProtocolType = "https"
	ProtocolTypeTCP   ProtocolType = "tcp"
)

var AllProtocolType = []ProtocolType{
	ProtocolTypeHTTP,
	ProtocolTypeHTTPS,
	ProtocolTypeTCP,
}

func (e ProtocolType) IsValid() bool {
	switch e {
	case ProtocolTypeHTTP, ProtocolTypeHTTPS, ProtocolTypeTCP:
		return true
	}
	return false
}

func (e ProtocolType) String() string {
	return string(e)
}

func (e *ProtocolType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProtocolType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProtocolType", str)
	}
	return nil
}

func (e ProtocolType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RedirectRuleStatus string

const (
	RedirectRuleStatusPending  RedirectRuleStatus = "pending"
	RedirectRuleStatusApplied  RedirectRuleStatus = "applied"
	RedirectRuleStatusFailed   RedirectRuleStatus = "failed"
	RedirectRuleStatusDeleting RedirectRuleStatus = "deleting"
)

var AllRedirectRuleStatus = []RedirectRuleStatus{
	RedirectRuleStatusPending,
	RedirectRuleStatusApplied,
	RedirectRuleStatusFailed,
	RedirectRuleStatusDeleting,
}

func (e RedirectRuleStatus) IsValid() bool {
	switch e {
	case RedirectRuleStatusPending, RedirectRuleStatusApplied, RedirectRuleStatusFailed, RedirectRuleStatusDeleting:
		return true
	}
	return false
}

func (e RedirectRuleStatus) String() string {
	return string(e)
}

func (e *RedirectRuleStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RedirectRuleStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RedirectRuleStatus", str)
	}
	return nil
}

func (e RedirectRuleStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UpstreamType string

const (
	UpstreamTypeGit        UpstreamType = "git"
	UpstreamTypeSourceCode UpstreamType = "sourceCode"
	UpstreamTypeImage      UpstreamType = "image"
)

var AllUpstreamType = []UpstreamType{
	UpstreamTypeGit,
	UpstreamTypeSourceCode,
	UpstreamTypeImage,
}

func (e UpstreamType) IsValid() bool {
	switch e {
	case UpstreamTypeGit, UpstreamTypeSourceCode, UpstreamTypeImage:
		return true
	}
	return false
}

func (e UpstreamType) String() string {
	return string(e)
}

func (e *UpstreamType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UpstreamType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UpstreamType", str)
	}
	return nil
}

func (e UpstreamType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
