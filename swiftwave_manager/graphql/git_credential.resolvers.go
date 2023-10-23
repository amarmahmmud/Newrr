package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.39

import (
	"context"
	GIT "github.com/swiftwave-org/swiftwave/git_manager"

	dbmodel "github.com/swiftwave-org/swiftwave/swiftwave_manager/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_manager/graphql/model"
)

// CreateGitCredential is the resolver for the createGitCredential field.
func (r *mutationResolver) CreateGitCredential(ctx context.Context, input model.GitCredentialInput) (*model.GitCredential, error) {
	record := gitCredentialInputToDatabaseObject(&input)
	tx := r.ServiceManager.DbClient.Create(&record)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return gitCredentialToGraphqlObject(record), nil
}

// UpdateGitCredential is the resolver for the updateGitCredential field.
func (r *mutationResolver) UpdateGitCredential(ctx context.Context, id int, input model.GitCredentialInput) (*model.GitCredential, error) {
	// fetch record
	var record dbmodel.GitCredential
	tx := r.ServiceManager.DbClient.First(&record, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	// update record
	record.Name = input.Name
	record.Username = input.Username
	record.Password = input.Password
	tx = r.ServiceManager.DbClient.Save(&record)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return gitCredentialToGraphqlObject(&record), nil
}

// DeleteGitCredential is the resolver for the deleteGitCredential field.
func (r *mutationResolver) DeleteGitCredential(ctx context.Context, id int) (*model.GitCredential, error) {
	// fetch record
	var record dbmodel.GitCredential
	tx := r.ServiceManager.DbClient.First(&record, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	// delete record
	tx = r.ServiceManager.DbClient.Delete(&record)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return gitCredentialToGraphqlObject(&record), nil
}

// GitCredentials is the resolver for the GitCredentials field.
func (r *queryResolver) GitCredentials(ctx context.Context) ([]*model.GitCredential, error) {
	var records []dbmodel.GitCredential
	tx := r.ServiceManager.DbClient.Find(&records)
	if tx.Error != nil {
		return nil, tx.Error
	}
	var gitCredentials []*model.GitCredential
	for _, record := range records {
		gitCredentials = append(gitCredentials, gitCredentialToGraphqlObject(&record))
	}
	return gitCredentials, nil
}

// GitCredential is the resolver for the GitCredential field.
func (r *queryResolver) GitCredential(ctx context.Context, id int) (*model.GitCredential, error) {
	var record dbmodel.GitCredential
	tx := r.ServiceManager.DbClient.First(&record, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return gitCredentialToGraphqlObject(&record), nil
}

// CheckGitCredentialRepositoryAccess is the resolver for the checkGitCredentialRepositoryAccess field.
func (r *queryResolver) CheckGitCredentialRepositoryAccess(ctx context.Context, input model.GitCredentialRepositoryAccessInput) (*model.GitCredentialRepositoryAccessResult, error) {
	// Fetch git credential
	var gitCredential dbmodel.GitCredential
	tx := r.ServiceManager.DbClient.First(&gitCredential, input.GitCredentialID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	// Prepare result object
	gitCredentialTestResult := &model.GitCredentialRepositoryAccessResult{
		GitCredentialID:  input.GitCredentialID,
		RepositoryURL:    input.RepositoryURL,
		RepositoryBranch: input.RepositoryBranch,
		GitCredential:    gitCredentialToGraphqlObject(&gitCredential),
	}
	// Test git credential
	_, err := GIT.FetchLatestCommitHash(input.RepositoryURL, input.RepositoryBranch, gitCredential.Username, gitCredential.Password)
	if err != nil {
		gitCredentialTestResult.Success = false
		gitCredentialTestResult.Error = "Git credential does not have access to the repository"
	} else {
		gitCredentialTestResult.Success = true
		gitCredentialTestResult.Error = ""
	}
	return gitCredentialTestResult, nil
}
