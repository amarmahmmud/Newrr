package server

import (
	"context"
	"fmt"
	"log"
	"strings"
)

func (src ApplicationSource) RepositoryURL() string {
	if src.GitProvider == GitProviderGithub {
		return "https://github.com/" + src.RepositoryUsername + "/" + src.RepositoryName + ".git"
	}
	if src.GitProvider == GitProviderGitlab {
		return "https://gitlab.com/" + src.RepositoryUsername + "/" + src.RepositoryName + ".git"
	}
	return ""
}

func (s *Server) AddLogToApplicationBuildLog(log_id string, message string, loglevel string, add_newline bool) {
	// push all the logs to redis topic for realtime log streaming
	// topic -> log_update/<log_id>

	var logRecord ApplicationBuildLog
	tx := s.DB_CLIENT.Where("id = ?", log_id).First(&logRecord)
	if tx.Error != nil {
		return
	}
	if add_newline {
		message += "\n"
	}

	// info, warning, error, success
	loglevel = strings.ToLower(loglevel)
	if loglevel == "error" {
		message = fmt.Sprintf("\x1b[1;31m\x1b[0m%s", message)
	} else if loglevel == "warning" {
		message = fmt.Sprintf("\x1b[1;33m\x1b[0m%s", message)
	} else if loglevel == "success" {
		message = fmt.Sprintf("\x1b[1;32m\x1b[0m%s", message)
	}

	// don't need special modification for -> loglevel == "info" || loglevel == ""
	logRecord.Logs = logRecord.Logs + message
	s.DB_CLIENT.Save(&logRecord)

	// push to redis topic
	s.REDIS_CLIENT.Publish(context.Background(), "log_update/"+log_id, message)
}

func (s *Server) MarkBuildLogAsCompleted(log_id string) {
	var logRecord ApplicationBuildLog
	tx := s.DB_CLIENT.Where("id = ?", log_id).First(&logRecord)
	if tx.Error != nil {
		return
	}
	logRecord.Completed = true
	s.DB_CLIENT.Save(&logRecord)
}

func (s *Server) CreateDefaultGitUser() {
	var git_credential GitCredential
	tx := s.DB_CLIENT.Where("name = ?", "default").First(&git_credential)
	if tx.Error != nil {
		log.Println("`default` git user not found, creating...")
		git_credential.Name = "default"
		git_credential.Username = ""
		git_credential.Password = ""
		tx2 := s.DB_CLIENT.Create(&git_credential)
		if tx2.Error != nil {
			log.Println("Failed to create `default` git user")
			panic(tx2.Error)
		}
	}
}

func (s ApplicationSource) GetSourceSummary() string {
	if s.Type == ApplicationSourceTypeGit {
		return fmt.Sprintf("%s Branch: %s", s.RepositoryURL(), s.Branch)
	}
	if s.Type == ApplicationSourceTypeTarball {
		return "Source Code"
	}
	if s.Type == ApplicationSourceTypeImage {
		return fmt.Sprintf("Image: %s", s.DockerImage)
	}
	return "Unknown"
}

func (s *Server) isProductionEnvironment() bool {
	return strings.Compare(s.ENVIRONMENT, "production") == 0
}

func (a *ApplicationBuildLog) GetRedisPubSubChannel() string {
	return "log_update/" + a.ID
}
