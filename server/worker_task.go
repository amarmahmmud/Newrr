package server

import (
	"log"

	"github.com/vmihailenco/taskq/v3"
)

/*
This file contains functions to store in tasks in `TASK_MAP`
This helps to fetch the task from `TASK_MAP` with the task name

Thid functions are called in `server/wroker.go` file
*/

// SSL Genration tasks
func (s *Server) RegisterSSLGenerateTask() {
	t := taskq.RegisterTask(&taskq.TaskOptions{
		Name: "ssl-generate",
		Handler: func(name string) error {
			return s.ProcessGenerateSSLRequestFromQueue(name)
		},
	})
	s.TASK_MAP["ssl-generate"] = t
}

// SSL Update in HaProxy tasks
func (s *Server) RegisterUpdateSSLHAProxyTask() {
	t := taskq.RegisterTask(&taskq.TaskOptions{
		Name: "ssl-update-haproxy",
		Handler: func(name string) error {
			return s.ProcessUpdateSSLHAProxyRequestFromQueue(name)
		},
	})
	s.TASK_MAP["ssl-update-haproxy"] = t
}

// Application deployment tasks
func (s *Server) RegisterDockerImageGenerationTask() {
	t := taskq.RegisterTask(&taskq.TaskOptions{
		Name: "docker-image-preparationAddServiceToDockerImageGenerationQueue",
		Handler: func(app_id uint, log_id string) error {
			err := s.ProcessDockerImageGenerationRequestFromQueue(app_id, log_id)
			if err != nil {
				log.Println(err)
			}
			return nil
		},
	})
	s.TASK_MAP["docker-image-preparation"] = t
}

func (s *Server) RegisterDeployServiceTask() {
	t := taskq.RegisterTask(&taskq.TaskOptions{
		Name: "deploy-service",
		Handler: func(app_id uint) error {
			err := s.ProcessDeployServiceRequestFromQueue(app_id)
			if err != nil {
				log.Println(err)
			}
			return nil
		},
	})
	s.TASK_MAP["deploy-service"] = t
}
