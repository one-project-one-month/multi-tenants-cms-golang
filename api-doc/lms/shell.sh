#!/bin/bash



# Helper function to write a file with model + service interface
write_tsp() {
  local filename=$1
  local namespace=$2
  local model_name=$3
  local route_name=$4
  local fields=$5

  cat > ./$filename << EOF
import "@typespec/http";

using Http;

@service(#{ title: "${model_name} Service" })
namespace ${namespace} {

  model ${model_name} {
${fields}
  }

  @route("/${route_name}")
  interface ${model_name}s {
    @get list(): ${model_name}[];
    @get read(@path id: string): ${model_name};
    @post create(@body body: ${model_name}): ${model_name};
    @patch update(@path id: string, @body body: Partial<${model_name}>): ${model_name};
    @delete delete(@path id: string): void;
  }
}
EOF
}

# Generate lms_user_role.tsp
write_tsp "lms_user_role.tsp" "LMSUserRoleService" "UserRole" "roles" "    id: string; // UUID\n    name: string; // lms_role_type"

# Generate tenants.tsp
write_tsp "tenants.tsp" "TenantsService" "Tenant" "tenants" "    id: string; // UUID\n    namespace: string;\n    cmsOwnerId: string; // UUID\n    createdAt: string; // timestamp\n    isActive: boolean;"

# Generate lms_user.tsp
write_tsp "lms_user.tsp" "LMSUserService" "User" "users" "    id: string; // UUID\n    email: string;\n    password: string;\n    roleId: string; // UUID\n    tenantId: string | null;\n    address: string | null;\n    phoneNumber: string | null;\n    registrationDate: string | null; // date\n    createdAt: string; // timestamp\n    updatedAt: string; // timestamp"

# Generate tenants_members.tsp
write_tsp "tenants_members.tsp" "TenantsMembersService" "TenantMember" "tenant-members" "    id: string; // UUID\n    userId: string; // UUID\n    tenantId: string; // UUID\n    joinedDate: string; // timestamp\n    isActive: boolean;"

# Generate course_category.tsp
write_tsp "course_category.tsp" "CourseCategoryService" "CourseCategory" "course-categories" "    id: string; // UUID\n    name: string;\n    description: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp"

# Generate course.tsp
write_tsp "course.tsp" "CourseService" "Course" "courses" "    id: string; // UUID\n    title: string;\n    description: string | null;\n    instructorId: string; // UUID\n    overallRating: int32 | null;\n    categoryId: string; // UUID\n    status: string; // course_status\n    durationDayCount: int32 | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp\n    ownedBy: string; // UUID"

# Generate namespace_consumer.tsp
write_tsp "namespace_consumer.tsp" "NamespaceConsumerService" "NamespaceConsumer" "namespace-consumers" "    id: string; // UUID\n    userId: string; // UUID\n    namespace: string;\n    joinedDate: string; // timestamp\n    isActive: boolean;"

# Generate enrollment.tsp
write_tsp "enrollment.tsp" "EnrollmentService" "Enrollment" "enrollments" "    id: string; // UUID\n    studentId: string; // UUID\n    courseId: string; // UUID\n    enrollmentDate: string; // timestamp\n    progress: float32 | null;\n    status: string; // enrollment_type\n    dueDate: string | null; // timestamp\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp"

# Generate certificate.tsp
write_tsp "certificate.tsp" "CertificateService" "Certificate" "certificates" "    id: string; // UUID\n    enrollmentId: string; // UUID\n    issueDate: string; // timestamp\n    certificateUrl: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp"

# Generate rating.tsp
write_tsp "rating.tsp" "RatingService" "Rating" "ratings" "    id: string; // UUID\n    userId: string; // UUID\n    courseId: string; // UUID\n    ratingCount: int32;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp"

# Generate module.tsp
write_tsp "module.tsp" "ModuleService" "Module" "modules" "    id: string; // UUID\n    name: string;\n    courseId: string; // UUID\n    description: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp"

# Generate quiz.tsp
write_tsp "quiz.tsp" "QuizService" "Quiz" "quizzes" "    id: string; // UUID\n    question: string;\n    answer: string;\n    moduleId: string; // UUID\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp"

# Generate student_quiz.tsp
write_tsp "student_quiz.tsp" "StudentQuizService" "StudentQuiz" "student-quizzes" "    id: string; // UUID\n    studentId: string; // UUID\n    quizId: string; // UUID\n    score: int32 | null;\n    attempt: int32;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp"

# Generate lesson.tsp
write_tsp "lesson.tsp" "LessonService" "Lesson" "lessons" "    id: string; // UUID\n    title: string;\n    content: string | null;\n    materialType: string | null;\n    moduleId: string; // UUID\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp"

# Generate assignment.tsp
write_tsp "assignment.tsp" "AssignmentService" "Assignment" "assignments" "    id: string; // UUID\n    courseId: string; // UUID\n    title: string;\n    instructions: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp"

# Generate submission.tsp
write_tsp "submission.tsp" "SubmissionService" "Submission" "submissions" "    id: string; // UUID\n    assignmentId: string; // UUID\n    studentId: string; // UUID\n    submittedAt: string; // timestamp\n    fileUrl: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp"

# Generate review.tsp
write_tsp "review.tsp" "ReviewService" "Review" "reviews" "    id: string; // UUID\n    courseId: string; // UUID\n    userId: string; // UUID\n    title: string | null;\n    description: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp"

# Generate report.tsp
write_tsp "report.tsp" "ReportService" "Report" "reports" "    id: string; // UUID\n    reportName: string;\n    generatedByUserId: string; // UUID\n    generatedDate: string; // timestamp\n    dataSnapshot: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp"

# Generate common.tsp for shared enums or types (empty for now)
cat > ./common.tsp << EOF
// Add shared enums and types here
EOF

# Generate main.tsp that imports and exposes all services

cat > ./main.tsp << EOF
import "./lms_user_role.tsp";
import "./tenants.tsp";
import "./lms_user.tsp";
import "./tenants_members.tsp";
import "./course_category.tsp";
import "./course.tsp";
import "./namespace_consumer.tsp";
import "./enrollment.tsp";
import "./certificate.tsp";
import "./rating.tsp";
import "./module.tsp";
import "./quiz.tsp";
import "./student_quiz.tsp";
import "./lesson.tsp";
import "./assignment.tsp";
import "./submission.tsp";
import "./review.tsp";
import "./report.tsp";

@service({ title: "LMS Combined Service" })
namespace LMS {
  // Expose all imported namespaces here if needed
}
EOF

echo "TypeSpec files generated in ./lms folder."
