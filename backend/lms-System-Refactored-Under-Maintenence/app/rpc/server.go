package rpc

import (
	assignmentsv "github.com/multi-tenants-cms-golang/lms-sys/app/rpc/assignment"
	coursesv "github.com/multi-tenants-cms-golang/lms-sys/app/rpc/course"
	modulesv "github.com/multi-tenants-cms-golang/lms-sys/app/rpc/module"
	tenantsv "github.com/multi-tenants-cms-golang/lms-sys/app/rpc/tenant"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	assignmentpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/assignment"
	coursepb "github.com/multi-tenants-cms-golang/lms-sys/protogen/course"
	modulepb "github.com/multi-tenants-cms-golang/lms-sys/protogen/module"
	tenantpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/tenant"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

// Server implements the gRPC services
type Server struct {
	assignmentpb.UnimplementedAssignmentServiceServer
	coursepb.UnimplementedCourseServiceServer
	modulepb.UnimplementedModuleServiceServer
	tenantpb.UnimplementedTenantServiceServer
	store  *db.Store
	logger *logrus.Logger
}

func NewServer(
	db *db.Store,
	logger *logrus.Logger,
) *Server {
	return &Server{
		store:  db,
		logger: logger,
	}
}

// Run starts the gRPC server
func (s *Server) Run() error {
	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err.Error())
	}
	grpcServer := grpc.NewServer()

	assignmentService := assignmentsv.NewAssignmentService()
	courseService := coursesv.NewCourseService()
	moduleService := modulesv.NewModuleService()
	tenantService := tenantsv.NewTenantService()

	// Register services with gRPC server
	assignmentpb.RegisterAssignmentServiceServer(grpcServer, assignmentService)
	coursepb.RegisterCourseServiceServer(grpcServer, courseService)
	modulepb.RegisterModuleServiceServer(grpcServer, moduleService)
	tenantpb.RegisterTenantServiceServer(grpcServer, tenantService)

	s.logger.Println("Starting server on port 9001")
	err = grpcServer.Serve(listener)
	if err != nil {
		s.logger.Fatal(err.Error())
		return err
	}
	return nil
}
