package course

import (
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	coursepb "github.com/multi-tenants-cms-golang/lms-sys/protogen/course"
)

// ProtoToModel converts protobuf message to database model
func ProtoToModel(pb *coursepb.Course) *repo.Course {
	if pb == nil {
		return nil
	}

	return &repo.Course{
		// Add your conversion fields here
		// Example:
		// ID:        pb.GetId(),
		// Name:      pb.GetName(),
		// CreatedAt: time.Unix(pb.GetCreatedAt().GetSeconds(), int64(pb.GetCreatedAt().GetNanos())),
	}
}

// ModelToProto converts database model to protobuf message
func ModelToProto(model *repo.Course) *coursepb.Course {
	if model == nil {
		return nil
	}

	return &coursepb.Course{
		// Add your conversion fields here
		// Example:
		// Id:        model.ID,
		// Name:      model.Name,
		// CreatedAt: timestamppb.New(model.CreatedAt),
	}
}

// ModelsToProtos converts slice of models to slice of protos
func ModelsToProtos(models []*repo.Course) []*coursepb.Course {
	protos := make([]*coursepb.Course, len(models))
	for i, model := range models {
		protos[i] = ModelToProto(model)
	}
	return protos
}

// ProtosToModels converts slice of protos to slice of models
func ProtosToModels(protos []*coursepb.Course) []*repo.Course {
	models := make([]*repo.Course, len(protos))
	for i, pb := range protos {
		models[i] = ProtoToModel(pb)
	}
	return models
}
