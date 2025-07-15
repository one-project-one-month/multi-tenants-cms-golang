package api

import (
	asv "github.com/multi-tenants-cms-golang/lms-sys/app/api/assignment"
	csv "github.com/multi-tenants-cms-golang/lms-sys/app/api/course"
	msv "github.com/multi-tenants-cms-golang/lms-sys/app/api/module"
	tsv "github.com/multi-tenants-cms-golang/lms-sys/app/api/tenants"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	pb "github.com/multi-tenants-cms-golang/lms-sys/protogen/assignments"
	cpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/courses"
	mpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
	tpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/tenants"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

// Server implements the AssignmentService
type Server struct {
	pb.UnimplementedAssignmentServiceServer // Embed the unimplemented server
	tpb.UnimplementedTenantServiceServer
	mpb.UnimplementedModuleServiceServer

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

// Run Protocol Buff ,proto buf
// REsponse:
func (s *Server) Run() error {
	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err.Error())
	}
	grpcServer := grpc.NewServer()

	assignmentService := asv.NewAssignmentsService(s.store, s.logger)
	tenantService := tsv.NewTenantService(s.store, s.logger)
	courseService := csv.NewCoursesService(s.store, s.logger)
	moduleService := msv.NewModuleService(s.store, s.logger)

	// Register services with gRPC server
	pb.RegisterAssignmentServiceServer(grpcServer, assignmentService)
	tpb.RegisterTenantServiceServer(grpcServer, tenantService)
	cpb.RegisterCourseServiceServer(grpcServer, courseService)
	mpb.RegisterModuleServiceServer(grpcServer, moduleService)

	s.logger.Println("Starting server on port 9001")
	err = grpcServer.Serve(listener)
	if err != nil {
		s.logger.Fatal(err.Error())
		return err
	}
	return nil
}
