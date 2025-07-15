package tenants

import (
	"github.com/multi-tenants-cms-golang/lms-system/internal/repo"
	tenantspb "github.com/multi-tenants-cms-golang/lms-system/protogen/tenants"
)

// ProtoToModel converts protobuf message to database model
func ProtoToModel(pb *tenantspb.Tenants) *repo.Tenant {
	if pb == nil {
		return nil
	}

	return &repo.Tenant{
		// Add your conversion fields here
		// Example:
		// ID:        pb.GetId(),
		// Name:      pb.GetName(),
		// CreatedAt: time.Unix(pb.GetCreatedAt().GetSeconds(), int64(pb.GetCreatedAt().GetNanos())),
	}
}

// ModelToProto converts database model to protobuf message
func ModelToProto(model *repo.Tenant) *tenantspb.Tenants {
	if model == nil {
		return nil
	}

	return &tenantspb.Tenants{
		// Add your conversion fields here
		// Example:
		// Id:        model.ID,
		// Name:      model.Name,
		// CreatedAt: timestamppb.New(model.CreatedAt),
	}
}

// ModelsToProtos converts slice of models to slice of protos
func ModelsToProtos(models []*repo.Tenant) []*tenantspb.Tenants {
	protos := make([]*tenantspb.Tenants, len(models))
	for i, model := range models {
		protos[i] = ModelToProto(model)
	}
	return protos
}

// ProtosToModels converts slice of protos to slice of models
func ProtosToModels(protos []*tenantspb.Tenants) []*repo.Tenant {
	models := make([]*repo.Tenant, len(protos))
	for i, pb := range protos {
		models[i] = ProtoToModel(pb)
	}
	return models
}
