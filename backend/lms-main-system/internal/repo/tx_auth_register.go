package repo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"strings"
)

// WholeRegistrationFlow handles the complete user registration process
// This is the main function that orchestrates the entire registration workflow
//
// REGISTRATION FLOW:
// 1. Validate namespace exists in the system
// 2. Get tenant ID associated with the namespace
// 3. Check if user is already a tenant member (affects role assignment)
// 4. Create new user in the database
// 5. Assign appropriate roles based on membership status:
//   - If user is tenant member: assign ADMIN role
//   - If user is not tenant member: assign default roles (STUDENT, VIEWER)
//
// 6. Return the created user details

func ProtoToCreateParams(req *authenticationpb.RegisterRequest, tenantId uuid.UUID, namespace string) CreateUserParams {
	domainEmail := emailToDomainEmail(req.Email, namespace)

	return CreateUserParams{
		LmsUserName:  req.Username,
		LmsUserEmail: req.Email,
		Password:     req.Password,
		TenantID: pgtype.UUID{
			Bytes: tenantId,
			Valid: true,
		},
		Address: pgtype.Text{
			String: req.Address,
			Valid:  req.Address != "",
		},
		NamespaceDomain: pgtype.Text{
			String: *domainEmail,
			Valid:  true,
		},
		PhoneNumber: pgtype.Text{
			String: req.PhoneNumber,
			Valid:  req.PhoneNumber != "",
		},
	}
}

func emailToDomainEmail(email, namespace string) *string {
	if !strings.Contains(email, "@") {
		return nil
	}
	name := strings.Split(email, "@")
	domainEmail := name[0] + "@" + strings.ToLower(namespace) + ".edu"
	return &domainEmail
}
func createUserRowToAssignRow(req *CreateUserRow, rolesId uuid.UUID) AssignRolesToUserParams {
	return AssignRolesToUserParams{
		LmsUserID: req.LmsUserID,
		LmsRoleID: rolesId,
	}
}

func (store *SQLStore) WholeRegistrationFlow(
	ctx context.Context,
	req *authenticationpb.RegisterRequest,
	namespace string,
) (*authenticationpb.RegisterResponse, error) {
	store.logger.WithFields(logrus.Fields{
		"username":  req.Username,
		"email":     req.Email,
		"namespace": namespace,
	}).Info("Starting user registration flow")

	err := store.execTx(ctx, func(q *Queries) error {
		store.logger.WithField("namespace", namespace).Info("Checking namespace existence")
		data, err := store.CheckNameSpaceFromMetaData(ctx, namespace)
		if err != nil {
			store.logger.WithFields(logrus.Fields{
				"namespace": namespace,
				"error":     err.Error(),
			}).Error("Failed to check namespace existence")
			return err
		}
		if data != 1 {
			store.logger.WithField("namespace", namespace).Error("Namespace does not exist")
			return errors.New("namespace does not exist")
		}

		store.logger.WithField("namespace", namespace).Info("Retrieving tenant ID")
		tenantIdInDB, err := store.GetTenantIdByNameSpace(ctx, namespace)
		if err != nil {
			store.logger.WithFields(logrus.Fields{
				"namespace": namespace,
				"error":     err.Error(),
			}).Error("Failed to get tenant ID by namespace")
			return err
		}
		store.logger.WithField("tenant_id", tenantIdInDB).Info("Successfully retrieved tenant ID")

		store.logger.WithField("email", req.Email).Info("Checking if user is already a tenant member")
		email, err := store.CheckTenantsMemberExistenceByEmail(ctx, req.Email)
		if err != nil {
			store.logger.WithFields(logrus.Fields{
				"email": req.Email,
				"error": err.Error(),
			}).Info("User does not belong to any tenant as member, proceeding with default role")
		}
		santizedPhone := utils.SanitizePhoneNumberEnhanced(req.PhoneNumber)
		userToBeCreated := ProtoToCreateParams(req, tenantIdInDB, namespace)
		userToBeCreated.PhoneNumber.String = santizedPhone
		userToBeCreated.Password, err = utils.HashPassword(userToBeCreated.Password)
		if err != nil {
			store.logger.WithFields(logrus.Fields{
				"email": req.Email,
				"error": err.Error(),
			}).Info("Failed to hash password")
		}
		store.logger.WithField("email", req.Email).Info("Creating new user")
		user, err := store.CreateUser(ctx, userToBeCreated)
		if err != nil {
			store.logger.WithFields(logrus.Fields{
				"email": req.Email,
				"error": err.Error(),
			}).Error("Failed to create user")
			return err
		}
		store.logger.WithFields(logrus.Fields{
			"user_id": user.LmsUserID,
			"email":   req.Email,
		}).Info("User created successfully")

		if email == 1 {
			store.logger.WithField("email", req.Email).Info("User is tenant member, assigning ADMIN role")
			roleIdInDB, err := store.GetRoleIdByName(ctx, "ADMIN")
			if err != nil {
				store.logger.WithFields(logrus.Fields{
					"role_name": "ADMIN",
					"error":     err.Error(),
				}).Error("Failed to get ADMIN role ID")
				return err
			}

			err = store.AssignRolesToUser(
				ctx,
				AssignRolesToUserParams{
					LmsUserID: user.LmsUserID,
					LmsRoleID: roleIdInDB,
				},
			)
			if err != nil {
				store.logger.WithFields(logrus.Fields{
					"user_id": user.LmsUserID,
					"role_id": roleIdInDB,
					"error":   err.Error(),
				}).Error("Failed to assign ADMIN role to user")
				return err
			}
			store.logger.WithField("user_id", user.LmsUserID).Info("Successfully assigned ADMIN role")
		} else {
			store.logger.WithField("email", req.Email).Info("User is not tenant member, assigning default roles")
			defaultRolesId, err := store.GetDefaultRoleIDs(ctx)
			if err != nil {
				store.logger.WithField("error", err.Error()).Error("Failed to get default role IDs")
				return err
			}

			store.logger.WithFields(logrus.Fields{
				"user_id":     user.LmsUserID,
				"roles_count": len(defaultRolesId),
			}).Info("Assigning default roles to user")

			for _, roleId := range defaultRolesId {
				rolesAssign := createUserRowToAssignRow(&user, roleId)
				err := store.AssignRolesToUser(ctx, rolesAssign)
				if err != nil {
					store.logger.WithFields(logrus.Fields{
						"user_id": user.LmsUserID,
						"role_id": roleId,
						"error":   err.Error(),
					}).Error("Failed to assign default role to user")
					return err
				}
			}
			store.logger.WithField("user_id", user.LmsUserID).Info("Successfully assigned all default roles")
		}
		return nil
	})

	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err.Error(),
		}).Error("Registration transaction failed")
		return nil, err
	}

	store.logger.WithField("email", req.Email).Info("Retrieving created user details")
	userInDB, err := store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err.Error(),
		}).Error("Failed to get user by email after creation")
		return nil, err
	}

	store.logger.WithFields(logrus.Fields{
		"user_id": userInDB.LmsUserID,
		"email":   userInDB.LmsUserEmail,
	}).Info("User registration completed successfully")

	return &authenticationpb.RegisterResponse{
		RegisteredUser: &authenticationpb.SystemUser{
			Id:               userInDB.LmsUserID.String(),
			Username:         userInDB.LmsUserName,
			Email:            userInDB.LmsUserEmail,
			PhoneNumber:      userInDB.PhoneNumber.String,
			Address:          userInDB.Address.String,
			EmailVerified:    userInDB.EmailVerified.Bool,
			MfaEnable:        userInDB.MfaEnable.Bool,
			DomainEmail:      userInDB.NamespaceDomain.String,
			RegistrationDate: utils.ParsePgTimestamp(pgtype.Timestamptz(userInDB.RegistrationDate)),
			CreatedAt:        utils.ParsePgTimestamp(pgtype.Timestamptz(userInDB.CreatedAt)),
			UpdatedAt:        utils.ParsePgTimestamp(pgtype.Timestamptz(userInDB.UpdatedAt)),
		},
		Message: "user created successfully",
	}, nil
}
