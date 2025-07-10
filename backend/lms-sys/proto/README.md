# LMS Protobuf Definitions

This directory contains Protocol Buffer definitions for the LMS (Learning Management System).

## Structure

### Core Tables
- `roles/` - User roles (LMS_ADMIN, STUDENT, INSTRUCTOR)
- `users/` - User management
- `tenants/` - Multi-tenant support
- `tenant_members/` - Tenant membership

### Course System
- `course_categories/` - Course categorization
- `courses/` - Course management
- `modules/` - Course modules
- `lessons/` - Individual lessons

### Assessment System
- `quizzes/` - Quiz management
- `student_quiz/` - Quiz attempts and scores
- `assignments/` - Assignment management
- `submissions/` - Assignment submissions

### Enrollment System
- `enrollments/` - Student enrollments
- `certificates/` - Certificate management

### Rating & Review System
- `ratings/` - Course ratings
- `reviews/` - Course reviews

### Other Systems
- `namespace_consumers/` - Consumer namespace management
- `reports/` - Reporting system

### Common
- `common/` - Shared types, enums, and utilities

## Usage

Each folder should contain:
- `{table_name}.proto` - Main protobuf definition
- `{table_name}_service.proto` - Service definitions (if applicable)

## Database Schema Mapping

This structure maps to the following database tables:

| Proto Folder | Database Table |
|--------------|----------------|
| roles | LMS_USER_Role |
| users | LMS_USER |
| tenants | Tenants |
| tenant_members | Tenants_Members |
| course_categories | Course_Category |
| courses | Course |
| modules | Module |
| lessons | Lesson |
| quizzes | Quiz |
| student_quiz | Student_Quiz |
| assignments | Assignment |
| submissions | Submission |
| enrollments | enrollment |
| certificates | Certificate |
| ratings | Rating |
| reviews | Review |
| namespace_consumers | namespace_consumer |
| reports | Report |
