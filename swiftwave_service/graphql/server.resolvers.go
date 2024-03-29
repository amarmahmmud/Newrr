package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"bytes"
	"context"
	"errors"
	"net"
	"strings"
	"time"

	containermanger "github.com/swiftwave-org/swiftwave/container_manager"
	"github.com/swiftwave-org/swiftwave/ssh_toolkit"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/graphql/model"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/logger"
	"gorm.io/gorm"
)

// CreateServer is the resolver for the createServer field.
func (r *mutationResolver) CreateServer(ctx context.Context, input model.NewServerInput) (*model.Server, error) {
	server := newServerInputToDatabaseObject(&input)
	err := core.CreateServer(&r.ServiceManager.DbClient, server)
	if err != nil {
		return nil, err
	}
	// if localhost, insert public key
	if server.IsLocalhost() {
		publicKey, err := r.Config.SystemConfig.PublicSSHKey()
		if err != nil {
			logger.GraphQLLoggerError.Println("Failed to generate public ssh key", err.Error())
		}
		// append the public key to current server ~/.ssh/authorized_keys
		err = AppendPublicSSHKeyLocally(publicKey)
		if err != nil {
			logger.GraphQLLoggerError.Println("Failed to append public ssh key", err.Error())
		}
	}
	return serverToGraphqlObject(server), nil
}

// TestSSHAccessToServer is the resolver for the testSSHAccessToServer field.
func (r *mutationResolver) TestSSHAccessToServer(ctx context.Context, id uint) (bool, error) {
	command := "echo 'Hi'"
	server, err := core.FetchServerByID(&r.ServiceManager.DbClient, id)
	if err != nil {
		return false, err
	}
	err = ssh_toolkit.ExecCommandOverSSH(command, nil, nil, 10, server.IP, 22, server.User, r.Config.SystemConfig.SshPrivateKey, 20)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckDependenciesOnServer is the resolver for the checkDependenciesOnServer field.
func (r *mutationResolver) CheckDependenciesOnServer(ctx context.Context, id uint) ([]*model.Dependency, error) {
	server, err := core.FetchServerByID(&r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	result := make([]*model.Dependency, 0)
	for _, dependency := range core.RequiredServerDependencies {
		if dependency == "init" {
			continue
		}
		stdoutBuffer := new(bytes.Buffer)
		err = ssh_toolkit.ExecCommandOverSSH(core.DependencyCheckCommands[dependency], stdoutBuffer, nil, 5, server.IP, 22, server.User, r.Config.SystemConfig.SshPrivateKey, 30)
		if err != nil {
			if strings.Contains(err.Error(), "exited with status 1") {
				result = append(result, &model.Dependency{Name: dependency, Available: false})
				continue
			} else {
				return nil, err
			}
		}
		if stdoutBuffer.String() == "" {
			result = append(result, &model.Dependency{Name: dependency, Available: false})
		} else {
			result = append(result, &model.Dependency{Name: dependency, Available: true})
		}
	}
	return result, nil
}

// InstallDependenciesOnServer is the resolver for the installDependenciesOnServer field.
func (r *mutationResolver) InstallDependenciesOnServer(ctx context.Context, id uint) (bool, error) {
	_, err := core.FetchServerByID(&r.ServiceManager.DbClient, id)
	if err != nil {
		return false, err
	}
	// Queue the request
	// - create a server log
	// - push the request to the queue
	serverLog := &core.ServerLog{
		ServerID: id,
		Title:    "Installing dependencies",
	}
	err = core.CreateServerLog(&r.ServiceManager.DbClient, serverLog)
	if err != nil {
		return false, err
	}
	// Push the request to the queue
	err = r.WorkerManager.EnqueueInstallDependenciesOnServerRequest(id, serverLog.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// SetupServer is the resolver for the setupServer field.
func (r *mutationResolver) SetupServer(ctx context.Context, input model.ServerSetupInput) (bool, error) {
	server, err := core.FetchServerByID(&r.ServiceManager.DbClient, input.ID)
	if err != nil {
		return false, err
	}

	// Check if all dependencies are installed
	installedDependencies, err := r.CheckDependenciesOnServer(ctx, input.ID)
	if err != nil {
		return false, err
	}
	for _, dependency := range installedDependencies {
		if !dependency.Available {
			return false, errors.New("dependency " + dependency.Name + " is not installed")
		}
	}

	// Proceed request logic (reject in any other case)
	// - if, want to be manager
	//    - if, there are some managers already, need to be online any of them
	//    - if, no servers, then it will be the first manager
	// - if, want to be worker
	//   - there need to be at least one manager
	var swarmManagerServer *core.Server
	if input.SwarmMode == model.SwarmModeManager {
		// Check if there are some servers already
		exists, err := core.IsPreparedServerExists(&r.ServiceManager.DbClient)
		if err != nil {
			return false, err
		}
		if exists {
			// Try to find out if there is any manager online
			r, err := core.FetchSwarmManager(&r.ServiceManager.DbClient)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return false, errors.New("swarm manager not found")
				} else {
					return false, err
				}
			}
			swarmManagerServer = &r
		}
	} else {
		// Check if there is any manager
		r, err := core.FetchSwarmManager(&r.ServiceManager.DbClient)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, errors.New("can't setup as worker, no swarm manager found")
			}
			return false, err
		}
		swarmManagerServer = &r
	}

	if swarmManagerServer == nil && input.SwarmMode == model.SwarmModeWorker {
		return false, errors.New("no manager found")
	}

	// NOTE: From here, if `swarmManagerServer` is nil, then this new server can be initialized as first swarm manager

	// Fetch hostname
	hostnameStdoutBuffer := new(bytes.Buffer)
	err = ssh_toolkit.ExecCommandOverSSH("cat /etc/hostname", hostnameStdoutBuffer, nil, 10, server.IP, 22, server.User, r.Config.SystemConfig.SshPrivateKey, 20)
	if err != nil {
		return false, err
	}
	hostname := strings.TrimSpace(hostnameStdoutBuffer.String())
	server.HostName = hostname
	server.Status = core.ServerPreparing
	server.SwarmMode = core.SwarmMode(input.SwarmMode)
	server.DockerUnixSocketPath = input.DockerUnixSocketPath
	err = core.UpdateServer(&r.ServiceManager.DbClient, server)
	if err != nil {
		return false, err
	}

	// Enqueue the request
	// - create a server log
	// - push the request to the queue
	serverLog := &core.ServerLog{
		ServerID: input.ID,
		Title:    "Setup server",
	}
	err = core.CreateServerLog(&r.ServiceManager.DbClient, serverLog)
	if err != nil {
		return false, err
	}
	// Push the request to the queue
	err = r.WorkerManager.EnqueueSetupServerRequest(input.ID, serverLog.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// PromoteServerToManager is the resolver for the promoteServerToManager field.
func (r *mutationResolver) PromoteServerToManager(ctx context.Context, id uint) (bool, error) {
	server, err := core.FetchServerByID(&r.ServiceManager.DbClient, id)
	if err != nil {
		return false, err
	}
	if server.Status != core.ServerOnline {
		return false, errors.New("server is not online")
	}
	// Fetch any swarm manager
	swarmManagerServer, err := core.FetchSwarmManagerExceptServer(&r.ServiceManager.DbClient, server.ID)
	if err != nil {
		return false, errors.New("no manager found")
	}
	// If there is any swarm manager, then promote this server to manager
	// Fetch net.Conn to the swarm manager
	conn, err := ssh_toolkit.NetConnOverSSH("unix", swarmManagerServer.DockerUnixSocketPath, 5, swarmManagerServer.IP, 22, swarmManagerServer.User, r.Config.SystemConfig.SshPrivateKey, 30)
	if err != nil {
		return false, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.GraphQLLoggerError.Println(err.Error())
		}
	}(conn)
	// Promote this server to manager
	manager, err := containermanger.New(ctx, conn)
	if err != nil {
		return false, err
	}
	err = manager.PromoteToManager(server.HostName)
	if err != nil {
		return false, err
	}
	server.SwarmMode = core.SwarmManager
	err = core.UpdateServer(&r.ServiceManager.DbClient, server)
	return err == nil, err
}

// DemoteServerToWorker is the resolver for the demoteServerToWorker field.
func (r *mutationResolver) DemoteServerToWorker(ctx context.Context, id uint) (bool, error) {
	server, err := core.FetchServerByID(&r.ServiceManager.DbClient, id)
	if err != nil {
		return false, err
	}
	if server.Status != core.ServerOnline {
		return false, errors.New("server is not online")
	}
	// Fetch any swarm manager
	swarmManagerServer, err := core.FetchSwarmManagerExceptServer(&r.ServiceManager.DbClient, server.ID)
	if err != nil {
		return false, errors.New("no manager found")
	}
	// If there is any swarm manager, then promote this server to manager
	// Fetch net.Conn to the swarm manager
	conn, err := ssh_toolkit.NetConnOverSSH("unix", swarmManagerServer.DockerUnixSocketPath, 5, swarmManagerServer.IP, 22, swarmManagerServer.User, r.Config.SystemConfig.SshPrivateKey, 30)
	if err != nil {
		return false, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.GraphQLLoggerError.Println(err.Error())
		}
	}(conn)
	// Promote this server to manager
	manager, err := containermanger.New(ctx, conn)
	if err != nil {
		return false, err
	}
	err = manager.DemoteToWorker(server.HostName)
	if err != nil {
		return false, err
	}
	server.SwarmMode = core.SwarmWorker
	err = core.UpdateServer(&r.ServiceManager.DbClient, server)
	return err == nil, err
}

// RestrictDeploymentOnServer is the resolver for the restrictDeploymentOnServer field.
func (r *mutationResolver) RestrictDeploymentOnServer(ctx context.Context, id uint) (bool, error) {
	server, err := core.FetchServerByID(&r.ServiceManager.DbClient, id)
	if err != nil {
		return false, err
	}
	if server.Status != core.ServerOnline {
		return false, errors.New("server is not online")
	}
	// Fetch any swarm manager
	swarmManagerServer, err := core.FetchSwarmManagerExceptServer(&r.ServiceManager.DbClient, server.ID)
	if err != nil {
		return false, errors.New("no manager found")
	}
	// If there is any swarm manager, then promote this server to manager
	// Fetch net.Conn to the swarm manager
	conn, err := ssh_toolkit.NetConnOverSSH("unix", swarmManagerServer.DockerUnixSocketPath, 5, swarmManagerServer.IP, 22, swarmManagerServer.User, r.Config.SystemConfig.SshPrivateKey, 30)
	if err != nil {
		return false, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.GraphQLLoggerError.Println(err.Error())
		}
	}(conn)
	// Promote this server to manager
	manager, err := containermanger.New(ctx, conn)
	if err != nil {
		return false, err
	}
	err = manager.MarkNodeAsDrained(server.HostName)
	if err != nil {
		return false, err
	}
	server.ScheduleDeployments = false
	err = core.UpdateServer(&r.ServiceManager.DbClient, server)
	return err == nil, err
}

// AllowDeploymentOnServer is the resolver for the allowDeploymentOnServer field.
func (r *mutationResolver) AllowDeploymentOnServer(ctx context.Context, id uint) (bool, error) {
	server, err := core.FetchServerByID(&r.ServiceManager.DbClient, id)
	if err != nil {
		return false, err
	}
	if server.Status != core.ServerOnline {
		return false, errors.New("server is not online")
	}
	// Fetch any swarm manager
	swarmManagerServer, err := core.FetchSwarmManagerExceptServer(&r.ServiceManager.DbClient, server.ID)
	if err != nil {
		return false, errors.New("no manager found")
	}
	// If there is any swarm manager, then promote this server to manager
	// Fetch net.Conn to the swarm manager
	conn, err := ssh_toolkit.NetConnOverSSH("unix", swarmManagerServer.DockerUnixSocketPath, 5, swarmManagerServer.IP, 22, swarmManagerServer.User, r.Config.SystemConfig.SshPrivateKey, 30)
	if err != nil {
		return false, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.GraphQLLoggerError.Println(err.Error())
		}
	}(conn)
	// Promote this server to manager
	manager, err := containermanger.New(ctx, conn)
	if err != nil {
		return false, err
	}
	err = manager.MarkNodeAsActive(server.HostName)
	if err != nil {
		return false, err
	}
	server.ScheduleDeployments = true
	err = core.UpdateServer(&r.ServiceManager.DbClient, server)
	return err == nil, err
}

// RemoveServerFromSwarmCluster is the resolver for the removeServerFromSwarmCluster field.
func (r *mutationResolver) RemoveServerFromSwarmCluster(ctx context.Context, id uint) (bool, error) {
	server, err := core.FetchServerByID(&r.ServiceManager.DbClient, id)
	if err != nil {
		return false, err
	}
	// Fetch any swarm manager
	swarmManagerServer, err := core.FetchSwarmManagerExceptServer(&r.ServiceManager.DbClient, server.ID)
	if err != nil {
		return false, errors.New("no manager found")
	}
	// If there is any swarm manager, then promote this server to manager
	// Fetch net.Conn to the swarm manager
	conn, err := ssh_toolkit.NetConnOverSSH("unix", swarmManagerServer.DockerUnixSocketPath, 5, swarmManagerServer.IP, 22, swarmManagerServer.User, r.Config.SystemConfig.SshPrivateKey, 30)
	if err != nil {
		return false, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.GraphQLLoggerError.Println(err.Error())
		}
	}(conn)
	manager, err := containermanger.New(ctx, conn)
	if err != nil {
		return false, err
	}
	err = manager.RemoveNode(server.HostName)
	if err != nil {
		return false, err
	}
	server.Status = core.ServerNeedsSetup
	err = core.UpdateServer(&r.ServiceManager.DbClient, server)
	if err == nil {
		// try to connect to the server and leave from the swarm
		serverConn, err2 := ssh_toolkit.NetConnOverSSH("unix", server.DockerUnixSocketPath, 5, server.IP, 22, server.User, r.Config.SystemConfig.SshPrivateKey, 30)
		if err2 == nil {
			defer func(serverConn net.Conn) {
				_ = serverConn.Close()
			}(serverConn)
			serverDockerManager, err2 := containermanger.New(ctx, serverConn)
			if err2 == nil {
				_ = serverDockerManager.LeaveSwarm()
			}
		}

	}
	return err == nil, err
}

// EnableProxyOnServer is the resolver for the enableProxyOnServer field.
func (r *mutationResolver) EnableProxyOnServer(ctx context.Context, id uint, typeArg model.ProxyType) (bool, error) {
	// Fetch the server
	server, err := core.FetchServerByID(&r.ServiceManager.DbClient, id)
	if err != nil {
		return false, err
	}
	if server.ProxyConfig.Enabled {
		return false, errors.New("proxy is already enabled")
	}
	// Set the proxy type
	server.ProxyConfig.Type = core.ProxyType(typeArg)
	// update in db
	err = core.ChangeProxyType(&r.ServiceManager.DbClient, server, server.ProxyConfig.Type)
	if err != nil {
		return false, err
	}
	// Enable the proxy
	server.ProxyConfig.SetupRunning = true
	// For backup proxy, atleast 1 active proxy is required
	if server.ProxyConfig.Type == core.BackupProxy {
		activeProxies, err := core.FetchProxyActiveServers(&r.ServiceManager.DbClient)
		if err != nil {
			return false, err
		}
		if len(activeProxies) == 0 {
			return false, errors.New("for adding backup proxy, atleast 1 active proxy is required")
		}
	}
	err = core.UpdateServer(&r.ServiceManager.DbClient, server)
	if err != nil {
		return false, err
	}
	// Create a server log
	serverLog := &core.ServerLog{
		ServerID: id,
		Title:    "Enable proxy on server " + server.HostName,
	}
	err = core.CreateServerLog(&r.ServiceManager.DbClient, serverLog)
	if err != nil {
		return false, err
	}
	// Queue the request
	err = r.WorkerManager.EnqueueSetupAndEnableProxyRequest(id, serverLog.ID)
	return err == nil, err
}

// DisableProxyOnServer is the resolver for the disableProxyOnServer field.
func (r *mutationResolver) DisableProxyOnServer(ctx context.Context, id uint) (bool, error) {
	// Fetch the server
	server, err := core.FetchServerByID(&r.ServiceManager.DbClient, id)
	if err != nil {
		return false, err
	}
	// Disable the proxy
	server.ProxyConfig.Enabled = false
	server.ProxyConfig.SetupRunning = false
	err = core.UpdateServer(&r.ServiceManager.DbClient, server)
	return err == nil, err
}

// FetchAnalyticsServiceToken is the resolver for the fetchAnalyticsServiceToken field.
func (r *mutationResolver) FetchAnalyticsServiceToken(ctx context.Context, id uint, rotate bool) (string, error) {
	var tokenRecord *core.AnalyticsServiceToken
	var err error
	if !rotate {
		tokenRecord, err = core.FetchAnalyticsServiceToken(ctx, r.ServiceManager.DbClient, id)
	} else {
		tokenRecord, err = core.RotateAnalyticsServiceToken(ctx, r.ServiceManager.DbClient, id)
	}
	if err != nil {
		return "", err
	} else {
		return tokenRecord.IDToken()
	}
}

// Servers is the resolver for the servers field.
func (r *queryResolver) Servers(ctx context.Context) ([]*model.Server, error) {
	servers, err := core.FetchAllServers(&r.ServiceManager.DbClient)
	if err != nil {
		return nil, err
	}
	serverList := make([]*model.Server, 0)
	for _, server := range servers {
		serverList = append(serverList, serverToGraphqlObject(&server))
	}
	return serverList, nil
}

// Server is the resolver for the server field.
func (r *queryResolver) Server(ctx context.Context, id uint) (*model.Server, error) {
	server, err := core.FetchServerByID(&r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	return serverToGraphqlObject(server), nil
}

// PublicSSHKey is the resolver for the publicSSHKey field.
func (r *queryResolver) PublicSSHKey(ctx context.Context) (string, error) {
	return r.Config.SystemConfig.PublicSSHKey()
}

// ServerResourceAnalytics is the resolver for the serverResourceAnalytics field.
func (r *queryResolver) ServerResourceAnalytics(ctx context.Context, id uint, timeframe model.ServerResourceAnalyticsTimeframe) ([]*model.ServerResourceAnalytics, error) {
	var previousTime time.Time = time.Now()
	switch timeframe {
	case model.ServerResourceAnalyticsTimeframeLast1Hour:
		previousTime = time.Now().Add(-1 * time.Hour)
	case model.ServerResourceAnalyticsTimeframeLast24Hours:
		previousTime = time.Now().Add(-24 * time.Hour)
	case model.ServerResourceAnalyticsTimeframeLast7Days:
		previousTime = time.Now().Add(-7 * 24 * time.Hour)
	case model.ServerResourceAnalyticsTimeframeLast30Days:
		previousTime = time.Now().Add(-30 * 24 * time.Hour)
	}
	previousTimeUnix := previousTime.Unix()

	// fetch the server resource analytics
	serverResourceStat, err := core.FetchServerResourceAnalytics(ctx, r.ServiceManager.DbClient, id, uint(previousTimeUnix))
	if err != nil {
		return nil, err
	}
	// convert the server resource analytics to graphql object
	serverResourceStatList := make([]*model.ServerResourceAnalytics, 0)
	for _, record := range serverResourceStat {
		serverResourceStatList = append(serverResourceStatList, serverResourceStatToGraphqlObject(record))
	}
	return serverResourceStatList, nil
}

// ServerDiskUsage is the resolver for the serverDiskUsage field.
func (r *queryResolver) ServerDiskUsage(ctx context.Context, id uint) ([]*model.ServerDisksUsage, error) {
	// fetch the server disk usage
	serverResourceUsageRecords, err := core.FetchServerDiskUsage(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	serverDiskStatsList := make([]*model.ServerDisksUsage, 0)
	for _, record := range serverResourceUsageRecords {
		val := severDisksStatToGraphqlObject(record.DiskStats, record.RecordedAt)
		serverDiskStatsList = append(serverDiskStatsList, &val)
	}
	return serverDiskStatsList, nil
}

// ServerLatestResourceAnalytics is the resolver for the serverLatestResourceAnalytics field.
func (r *queryResolver) ServerLatestResourceAnalytics(ctx context.Context, id uint) (*model.ServerResourceAnalytics, error) {
	// fetch the latest server resource analytics
	serverResourceStat, err := core.FetchLatestServerResourceAnalytics(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	// convert the server resource analytics to graphql object
	return serverResourceStatToGraphqlObject(serverResourceStat), nil
}

// ServerLatestDiskUsage is the resolver for the serverLatestDiskUsage field.
func (r *queryResolver) ServerLatestDiskUsage(ctx context.Context, id uint) (*model.ServerDisksUsage, error) {
	// fetch the latest server disk usage
	serverDiskStats, timestamp, err := core.FetchLatestServerDiskUsage(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	// convert the server disk usage to graphql object
	res := severDisksStatToGraphqlObject(*serverDiskStats, *timestamp)
	return &res, nil
}

// Logs is the resolver for the logs field.
func (r *serverResolver) Logs(ctx context.Context, obj *model.Server) ([]*model.ServerLog, error) {
	serverLogs, err := core.FetchServerLogByServerID(&r.ServiceManager.DbClient, obj.ID)
	if err != nil {
		return nil, err
	}
	serverLogList := make([]*model.ServerLog, 0)
	for _, serverLog := range serverLogs {
		serverLogList = append(serverLogList, serverLogToGraphqlObject(&serverLog))
	}
	return serverLogList, nil
}

// Server returns ServerResolver implementation.
func (r *Resolver) Server() ServerResolver { return &serverResolver{r} }

type serverResolver struct{ *Resolver }
