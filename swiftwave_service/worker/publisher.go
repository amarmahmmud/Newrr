package worker

func (m Manager) EnqueueBuildApplicationRequest(applicationId string, deploymentId string) error {
	return m.ServiceManager.TaskQueueClient.EnqueueTask(buildApplicationQueueName, BuildApplicationRequest{
		AppId:        applicationId,
		DeploymentId: deploymentId,
	})
}

func (m Manager) EnqueueDeployApplicationRequest(applicationId string) error {
	return m.ServiceManager.TaskQueueClient.EnqueueTask(deployApplicationQueueName, DeployApplicationRequest{
		AppId: applicationId,
	})
}

func (m Manager) EnqueueDeleteApplicationRequest(applicationId string) error {
	return m.ServiceManager.TaskQueueClient.EnqueueTask(deleteApplicationQueueName, DeleteApplicationRequest{
		Id: applicationId,
	})
}

func (m Manager) EnqueueSSLGenerateRequest(domainId uint) error {
	return m.ServiceManager.TaskQueueClient.EnqueueTask(sslGenerateQueueName, SSLGenerateRequest{
		DomainId: domainId,
	})
}

func (m Manager) EnqueueIngressRuleApplyRequest(ingressRuleId uint) error {
	return m.ServiceManager.TaskQueueClient.EnqueueTask(ingressRuleApplyQueueName, IngressRuleApplyRequest{
		Id: ingressRuleId,
	})
}

func (m Manager) EnqueueIngressRuleDeleteRequest(ingressRuleId uint) error {
	return m.ServiceManager.TaskQueueClient.EnqueueTask(ingressRuleDeleteQueueName, IngressRuleDeleteRequest{
		Id: ingressRuleId,
	})
}

func (m Manager) EnqueueRedirectRuleApplyRequest(redirectRuleId uint) error {
	return m.ServiceManager.TaskQueueClient.EnqueueTask(redirectRuleApplyQueueName, RedirectRuleApplyRequest{
		Id: redirectRuleId,
	})
}

func (m Manager) EnqueueRedirectRuleDeleteRequest(redirectRuleId uint) error {
	return m.ServiceManager.TaskQueueClient.EnqueueTask(redirectRuleDeleteQueueName, RedirectRuleDeleteRequest{
		Id: redirectRuleId,
	})
}
