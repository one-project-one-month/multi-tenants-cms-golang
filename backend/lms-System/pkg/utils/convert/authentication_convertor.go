package convertor

import (
	"database/sql"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
)

func RegisterProtoToModel(pb *authenticationpb.RegisterRequest) *repo.RegisterLMSUserParams {
	if pb == nil {
		return nil
	}
	hashedPassword, err := utils.HashPassword(pb.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Failed to hash password")
		return nil
	}
	return &repo.RegisterLMSUserParams{
		LmsUserName:  pb.Username,
		LmsUserEmail: pb.Email,
		Password:     hashedPassword,
		Address: sql.NullString{
			String: pb.Address,
			Valid:  true,
		},
		PhoneNumber: sql.NullString{
			String: pb.PhoneNumber,
			Valid:  true,
		},
	}
}
func RegisterModelToProto(model *repo.RegisterLMSUserRow) *authenticationpb.RegisterResponse {
	return &authenticationpb.RegisterResponse{
		RegisteredUser: &authenticationpb.SystemUser{
			Id:               model.LmsUserID.String(),
			Username:         model.LmsUserName,
			Email:            model.LmsUserEmail,
			PhoneNumber:      model.PhoneNumber.String,
			Address:          model.Address.String,
			RegistrationDate: NullableTimeToProto(model.RegistrationDate),
			CreatedAt:        NullableTimeToProto(model.UpdatedAt),
			UpdatedAt:        NullableTimeToProto(model.UpdatedAt),
		},
		Message: "",
	}
}

//// ModelToProto converts database model to protobuf message
//func ModelToProto(model *repo.Authentication) *authenticationpb.Authentication {
//	if model == nil {
//		return nil
//	}
//
//	return &authenticationpb.Authentication{
//		// Add your conversion fields here
//		// Example:
//		// Id:        model.ID,
//		// Name:      model.Name,
//		// CreatedAt: timestamppb.New(model.CreatedAt),
//	}
//}
//
//// ModelsToProtos converts slice of models to slice of protos
//func ModelsToProtos(models []*repo.Authentication) []*authenticationpb.Authentication {
//	protos := make([]*authenticationpb.Authentication, len(models))
//	for i, model := range models {
//		protos[i] = ModelToProto(model)
//	}
//	return protos
//}
//
//// ProtosToModels converts slice of protos to slice of models
//func ProtosToModels(protos []*authenticationpb.Authentication) []*repo.Authentication {
//	models := make([]*repo.Authentication, len(protos))
//	for i, pb := range protos {
//		models[i] = ProtoToModel(pb)
//	}
//	return models
//}
