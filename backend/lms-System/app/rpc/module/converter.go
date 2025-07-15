package module

import (
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	modulepb "github.com/multi-tenants-cms-golang/lms-sys/protogen/module"
)

// ProtoToModel converts protobuf message to database model
func ProtoToModel(pb *modulepb.Module) *repo.Module {
	if pb == nil {
		return nil
	}

	return &repo.Module{
		// Add your conversion fields here
		// Example:
		// ID:        pb.GetId(),
		// Name:      pb.GetName(),
		// CreatedAt: time.Unix(pb.GetCreatedAt().GetSeconds(), int64(pb.GetCreatedAt().GetNanos())),
	}
}

// ModelToProto converts database model to protobuf message
func ModelToProto(model *repo.Module) *modulepb.Module {
	if model == nil {
		return nil
	}

	return &modulepb.Module{
		// Add your conversion fields here
		// Example:
		// Id:        model.ID,
		// Name:      model.Name,
		// CreatedAt: timestamppb.New(model.CreatedAt),
	}
}

// ModelsToProtos converts slice of models to slice of protos
func ModelsToProtos(models []*repo.Module) []*modulepb.Module {
	protos := make([]*modulepb.Module, len(models))
	for i, model := range models {
		protos[i] = ModelToProto(model)
	}
	return protos
}

// ProtosToModels converts slice of protos to slice of models
func ProtosToModels(protos []*modulepb.Module) []*repo.Module {
	models := make([]*repo.Module, len(protos))
	for i, pb := range protos {
		models[i] = ProtoToModel(pb)
	}
	return models
}
