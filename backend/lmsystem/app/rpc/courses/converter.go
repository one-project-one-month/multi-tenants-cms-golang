package courses

import (
	"github.com/multi-tenants-cms-golang/lms-system/internal/repo"
	coursespb "github.com/multi-tenants-cms-golang/lms-system/protogen/courses"
)

// ProtoToModel converts protobuf message to database model
func ProtoToModel(pb *coursespb.Courses) *repo.Courses {
	if pb == nil {
		return nil
	}

	return &repo.Courses{
		// Add your conversion fields here
		// Example:
		// ID:        pb.GetId(),
		// Name:      pb.GetName(),
		// CreatedAt: time.Unix(pb.GetCreatedAt().GetSeconds(), int64(pb.GetCreatedAt().GetNanos())),
	}
}

// ModelToProto converts database model to protobuf message
func ModelToProto(model *repo.Courses) *coursespb.Courses {
	if model == nil {
		return nil
	}

	return &coursespb.Courses{
		// Add your conversion fields here
		// Example:
		// Id:        model.ID,
		// Name:      model.Name,
		// CreatedAt: timestamppb.New(model.CreatedAt),
	}
}

// ModelsToProtos converts slice of models to slice of protos
func ModelsToProtos(models []*repo.Courses) []*coursespb.Courses {
	protos := make([]*coursespb.Courses, len(models))
	for i, model := range models {
		protos[i] = ModelToProto(model)
	}
	return protos
}

// ProtosToModels converts slice of protos to slice of models
func ProtosToModels(protos []*coursespb.Courses) []*repo.Courses {
	models := make([]*repo.Courses, len(protos))
	for i, pb := range protos {
		models[i] = ProtoToModel(pb)
	}
	return models
}
