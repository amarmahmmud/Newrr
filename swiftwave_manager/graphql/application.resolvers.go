package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.39

import (
	"context"
	"fmt"

	dbmodel "github.com/swiftwave-org/swiftwave/swiftwave_manager/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_manager/graphql/model"
)

// EnvironmentVariables is the resolver for the environmentVariables field.
func (r *applicationResolver) EnvironmentVariables(ctx context.Context, obj *model.Application) ([]*model.EnvironmentVariable, error) {
	// fetch record
	records, err := dbmodel.FindEnvironmentVariablesByApplicationId(ctx, r.ServiceManager.DbClient, obj.ID)
	if err != nil {
		return nil, err
	}
	// convert to graphql object
	var result = make([]*model.EnvironmentVariable, 0)
	for _, record := range records {
		result = append(result, environmentVariableToGraphqlObject(record))
	}
	return result, nil
}

// PersistentVolumeBindings is the resolver for the persistentVolumeBindings field.
func (r *applicationResolver) PersistentVolumeBindings(ctx context.Context, obj *model.Application) ([]*model.PersistentVolumeBinding, error) {
	// fetch record
	records, err := dbmodel.FindPersistentVolumeBindingsByApplicationId(ctx, r.ServiceManager.DbClient, obj.ID)
	if err != nil {
		return nil, err
	}
	// convert to graphql object
	var result = make([]*model.PersistentVolumeBinding, 0)
	for _, record := range records {
		result = append(result, persistentVolumeBindingToGraphqlObject(record))
	}
	return result, nil
}

// LatestDeployment is the resolver for the latestDeployment field.
func (r *applicationResolver) LatestDeployment(ctx context.Context, obj *model.Application) (*model.Deployment, error) {
	// fetch record
	record, err := dbmodel.FindLatestDeploymentByApplicationId(ctx, r.ServiceManager.DbClient, obj.ID)
	if err != nil {
		return nil, err
	}
	return deploymentToGraphqlObject(record), nil
}

// Deployments is the resolver for the deployments field.
func (r *applicationResolver) Deployments(ctx context.Context, obj *model.Application) ([]*model.Deployment, error) {
	panic(fmt.Errorf("not implemented: Deployments - deployments"))
}

// CreateApplication is the resolver for the createApplication field.
func (r *mutationResolver) CreateApplication(ctx context.Context, input model.ApplicationInput) (*model.Application, error) {
	record := applicationInputToDatabaseObject(&input)
	// create transaction
	transaction := r.ServiceManager.DbClient.Begin()
	err := record.Create(ctx, *transaction, r.ServiceManager.DockerManager, r.ServiceConfig.CodeTarballDir)
	if err != nil {
		transaction.Rollback()
		return nil, err
	}
	err = transaction.Commit().Error
	if err != nil {
		return nil, err
	}
	return applicationToGraphqlObject(record), nil
}

// UpdateApplication is the resolver for the updateApplication field.
func (r *mutationResolver) UpdateApplication(ctx context.Context, id string, input model.ApplicationInput) (*model.Application, error) {
	// fetch record
	var record = &dbmodel.Application{}
	err := record.FindById(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	// convert input to database object
	var databaseObject = applicationInputToDatabaseObject(&input)
	databaseObject.ID = record.ID
	// update record
	err = databaseObject.Update(ctx, r.ServiceManager.DbClient, r.ServiceManager.DockerManager)
	if err != nil {
		return nil, err
	}
	return applicationToGraphqlObject(databaseObject), nil
}

// DeleteApplication is the resolver for the deleteApplication field.
func (r *mutationResolver) DeleteApplication(ctx context.Context, id string) (*model.Application, error) {
	// fetch record
	var record = &dbmodel.Application{}
	err := record.FindById(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	// delete record
	err = record.Delete(ctx, r.ServiceManager.DbClient, r.ServiceManager.DockerManager)
	if err != nil {
		return nil, err
	}
	return applicationToGraphqlObject(record), nil
}

// Application is the resolver for the application field.
func (r *queryResolver) Application(ctx context.Context, id string) (*model.Application, error) {
	var record = &dbmodel.Application{}
	err := record.FindById(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	return applicationToGraphqlObject(record), nil
}

// Applications is the resolver for the applications field.
func (r *queryResolver) Applications(ctx context.Context) ([]*model.Application, error) {
	var records []*dbmodel.Application
	records, err := dbmodel.FindAllApplications(ctx, r.ServiceManager.DbClient)
	if err != nil {
		return nil, err
	}
	var result = make([]*model.Application, 0)
	for _, record := range records {
		result = append(result, applicationToGraphqlObject(record))
	}
	return result, nil
}

// IsExistApplicationName is the resolver for the isExistApplicationName field.
func (r *queryResolver) IsExistApplicationName(ctx context.Context, name string) (bool, error) {
	return dbmodel.IsExistApplicationName(ctx, r.ServiceManager.DbClient, r.ServiceManager.DockerManager, name)
}

// Application returns ApplicationResolver implementation.
func (r *Resolver) Application() ApplicationResolver { return &applicationResolver{r} }

type applicationResolver struct{ *Resolver }
