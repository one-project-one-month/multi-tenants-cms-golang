package assignment

import (
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	assignmentpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/assignment"
)

// ProtoToModel converts protobuf message to database model
func ProtoToModel(pb *assignmentpb.Assignment) *repo.Assignment {
	if pb == nil {
		return nil
	}

	return &repo.Assignment{
		// Add your conversion fields here
		// Example:
		// ID:        pb.GetId(),
		// Name:      pb.GetName(),
		// CreatedAt: time.Unix(pb.GetCreatedAt().GetSeconds(), int64(pb.GetCreatedAt().GetNanos())),
	}
}

// ModelToProto converts database model to protobuf message
func ModelToProto(model *repo.Assignment) *assignmentpb.Assignment {
	if model == nil {
		return nil
	}

	return &assignmentpb.Assignment{
		// Add your conversion fields here
		// Example:
		// Id:        model.ID,
		// Name:      model.Name,
		// CreatedAt: timestamppb.New(model.CreatedAt),
	}
}

// ModelsToProtos converts slice of models to slice of protos
func ModelsToProtos(models []*repo.Assignment) []*assignmentpb.Assignment {
	protos := make([]*assignmentpb.Assignment, len(models))
	for i, model := range models {
		protos[i] = ModelToProto(model)
	}
	return protos
}

// ProtosToModels converts slice of protos to slice of models
func ProtosToModels(protos []*assignmentpb.Assignment) []*repo.Assignment {
	models := make([]*repo.Assignment, len(protos))
	for i, pb := range protos {
		models[i] = ProtoToModel(pb)
	}
	return models
}
