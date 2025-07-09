#!/bin/bash

# scripts/generate_proto_structure.sh
# Script to generate protobuf folder structure for LMS tables

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROTO_DIR="proto"

# Function to print colored output
print_color() {
    printf "${1}%s${NC}\n" "$2"
}

# Function to create directory if it doesn't exist
create_dir() {
    if [ ! -d "$1" ]; then
        mkdir -p "$1"
        print_color $GREEN "Created directory: $1"
    else
        print_color $YELLOW "Directory already exists: $1"
    fi
}

print_color $BLUE "Starting protobuf folder structure generation for LMS system..."

# Create main proto directory
create_dir "$PROTO_DIR"

# Create folders for each table/domain
print_color $BLUE "Creating folders for LMS tables..."

# Core tables
create_dir "$PROTO_DIR/roles"
create_dir "$PROTO_DIR/users"
create_dir "$PROTO_DIR/tenants"
create_dir "$PROTO_DIR/tenant_members"

# Course related tables
create_dir "$PROTO_DIR/course_categories"
create_dir "$PROTO_DIR/courses"
create_dir "$PROTO_DIR/modules"
create_dir "$PROTO_DIR/lessons"

# Assessment related tables
create_dir "$PROTO_DIR/quizzes"
create_dir "$PROTO_DIR/student_quiz"
create_dir "$PROTO_DIR/assignments"
create_dir "$PROTO_DIR/submissions"

# Enrollment and progress
create_dir "$PROTO_DIR/enrollments"
create_dir "$PROTO_DIR/certificates"

# Rating and review system
create_dir "$PROTO_DIR/ratings"
create_dir "$PROTO_DIR/reviews"

# Consumer and reporting
create_dir "$PROTO_DIR/namespace_consumers"
create_dir "$PROTO_DIR/reports"

# Common types and utilities
create_dir "$PROTO_DIR/common"

# Create a README file in proto directory
cat > "$PROTO_DIR/README.md" << 'EOF'
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
EOF

# Create .gitkeep files in empty directories
print_color $BLUE "Creating .gitkeep files..."

folders=(
    "$PROTO_DIR/roles"
    "$PROTO_DIR/users"
    "$PROTO_DIR/tenants"
    "$PROTO_DIR/tenant_members"
    "$PROTO_DIR/course_categories"
    "$PROTO_DIR/courses"
    "$PROTO_DIR/modules"
    "$PROTO_DIR/lessons"
    "$PROTO_DIR/quizzes"
    "$PROTO_DIR/student_quiz"
    "$PROTO_DIR/assignments"
    "$PROTO_DIR/submissions"
    "$PROTO_DIR/enrollments"
    "$PROTO_DIR/certificates"
    "$PROTO_DIR/ratings"
    "$PROTO_DIR/reviews"
    "$PROTO_DIR/namespace_consumers"
    "$PROTO_DIR/reports"
    "$PROTO_DIR/common"
)


print_color $GREEN "✅ Protobuf folder structure created successfully!"
print_color $BLUE "📁 Proto directory structure:"
tree $PROTO_DIR 2>/dev/null || find $PROTO_DIR -type d | sed 's|[^/]*/|  |g'

print_color $YELLOW "Next steps:"
print_color $YELLOW "1. Add protobuf definitions to each folder"
print_color $YELLOW "2. Define services and messages for each domain"
print_color $YELLOW "3. Generate Go code using protoc"