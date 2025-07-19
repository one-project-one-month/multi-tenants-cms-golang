package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	dbHost     = "ep-purple-haze-a15hvjtw.ap-southeast-1.aws.neon.tech"
	dbPort     = 5432
	dbUser     = "neondb_owner"
	dbPassword = "npg_mFrGIyn0Sq9B"
	dbName     = "neondb"
)

type Seeder struct {
	db *sql.DB
}

func main() {
	psqlInfo := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=require&prefer_simple_protocol=true&statement_cache_mode=describe",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Failed to close database:", err)
		}
	}(db)

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to database!")

	seeder := &Seeder{db: db}

	//if err := seeder.seedRoles(); err != nil {
	//	log.Fatal("Failed to seed roles:", err)
	//}
	//
	//if err := seeder.seedUsers(); err != nil {
	//	log.Fatal("Failed to seed users:", err)
	//}
	//
	//if err := seeder.seedMFATokens(); err != nil {
	//	log.Fatal("Failed to seed MFA tokens:", err)
	//}

	//if err := seeder.seedPurchases(); err != nil {
	//	log.Fatal("Failed to seed purchases:", err)
	//}

	//if err := seeder.seedPageRequests(); err != nil {
	//	log.Fatal("Failed to seed page requests:", err)
	//}
	//
	//if err := seeder.seedCMSPageRequests(); err != nil {
	//	log.Fatal("Failed to seed CMS page requests:", err)
	//}

	if err := seeder.seedCMSPages(); err != nil {
		log.Fatal("Failed to seed CMS pages:", err)
	}

	fmt.Println("Database seeding completed successfully!")
}

func (s *Seeder) seedRoles() error {
	fmt.Println("Seeding roles...")

	roles := []string{
		"CMS_ADMIN",
		"CMS_CUSTOMER",
		"CMS_EDITOR",
		"CMS_DEVELOPER",
		"CMS_SUPPORT",
	}

	for _, role := range roles {
		_, err := s.db.Exec(`
            INSERT INTO cms_whole_sys_role (role_name) 
            VALUES ($1) 
            ON CONFLICT (role_name) DO NOTHING`,
			role)
		if err != nil {
			return fmt.Errorf("failed to insert role %s: %w", role, err)
		}
	}

	fmt.Printf("Seeded %d roles\n", len(roles))
	return nil
}

func (s *Seeder) seedUsers() error {
	fmt.Println("Seeding users...")

	// Generate 50+ users with realistic data
	firstNames := []string{
		"John", "Jane", "Robert", "Emily", "Michael", "Sarah", "David", "Jessica",
		"William", "Jennifer", "Richard", "Lisa", "Joseph", "Nancy", "Thomas",
		"Karen", "Charles", "Betty", "Christopher", "Margaret", "Daniel", "Sandra",
		"Matthew", "Ashley", "Anthony", "Kimberly", "Donald", "Donna", "Mark",
		"Dorothy", "Paul", "Michelle", "Steven", "Carol", "Andrew", "Amanda",
		"Kenneth", "Melissa", "Joshua", "Deborah", "Kevin", "Stephanie", "Brian",
		"Rebecca", "George", "Laura", "Edward", "Sharon", "Ronald", "Cynthia",
	}

	lastNames := []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones", "Miller", "Davis",
		"Garcia", "Rodriguez", "Wilson", "Martinez", "Anderson", "Taylor",
		"Thomas", "Hernandez", "Moore", "Martin", "Jackson", "Thompson",
		"White", "Lopez", "Lee", "Gonzalez", "Harris", "Clark", "Lewis",
		"Robinson", "Walker", "Perez", "Hall", "Young", "Allen", "Sanchez",
		"Wright", "King", "Scott", "Green", "Baker", "Adams", "Nelson",
		"Hill", "Ramirez", "Campbell", "Mitchell", "Roberts", "Carter",
		"Phillips", "Evans", "Turner", "Torres",
	}

	domains := []string{"example.com", "test.org", "demo.net", "sample.edu", "company.io"}
	namespaces := []string{"lms_system", "customers", "partners", "developers", "support", "testing"}

	users := make([]struct {
		name       string
		email      string
		namespace  string
		password   string
		role       string
		verified   bool
		mfaEnabled bool
	}, 0, 50)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 50; i++ {
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		domain := domains[rand.Intn(len(domains))]
		namespace := namespaces[rand.Intn(len(namespaces))]

		// First 5 users are admins/editors/support
		var role string
		if i < 2 {
			role = "CMS_ADMIN"
		} else if i < 4 {
			role = "CMS_EDITOR"
		} else if i < 6 {
			role = "CMS_SUPPORT"
		} else if i < 8 {
			role = "CMS_DEVELOPER"
		} else {
			role = "CMS_CUSTOMER"
		}

		users = append(users, struct {
			name       string
			email      string
			namespace  string
			password   string
			role       string
			verified   bool
			mfaEnabled bool
		}{
			name:       fmt.Sprintf("%s %s", firstName, lastName),
			email:      fmt.Sprintf("%s.%s%d@%s", firstName, lastName, i, domain),
			namespace:  namespace,
			password:   fmt.Sprintf("%s%s%d!", firstName, lastName, i),
			role:       role,
			verified:   rand.Intn(2) == 1, // 50% chance
			mfaEnabled: rand.Intn(2) == 1, // 50% chance
		})
	}

	for _, user := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", user.email, err)
		}

		userID := uuid.New()

		_, err = s.db.Exec(`
            INSERT INTO cms_user (
                cms_user_id, cms_user_name, cms_user_email, cms_name_space, 
                password, cms_user_role, verified, mfa_enabled, created_at, updated_at
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
            ON CONFLICT (cms_user_email) DO NOTHING`,
			userID, user.name, user.email, user.namespace,
			string(hashedPassword), user.role, user.verified, user.mfaEnabled,
			time.Now().Add(-time.Duration(rand.Intn(365))*24*time.Hour), // Random creation date in past year
			time.Now())
		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", user.email, err)
		}
	}

	fmt.Printf("Seeded %d users\n", len(users))
	return nil
}

func (s *Seeder) seedMFATokens() error {
	fmt.Println("Seeding MFA tokens...")

	rows, err := s.db.Query("SELECT cms_user_id FROM cms_user WHERE mfa_enabled = true")
	if err != nil {
		return fmt.Errorf("failed to get MFA users: %w", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			continue
		}

		// Generate a random MFA token
		mfaToken := generateRandomMFAToken()

		_, err = s.db.Exec(`
            INSERT INTO mfa_token (mfa_token, user_id, expires_at, created_at)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (user_id) DO NOTHING`,
			mfaToken, userID, time.Now().Add(24*time.Hour), time.Now())
		if err != nil {
			return fmt.Errorf("failed to insert MFA token: %w", err)
		}
		count++
	}

	fmt.Printf("Seeded %d MFA tokens\n", count)
	return nil
}

func generateRandomMFAToken() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	token := make([]byte, 16)
	rand.Seed(time.Now().UnixNano())
	for i := range token {
		token[i] = charset[rand.Intn(len(charset))]
	}
	return string(token)
}

func (s *Seeder) seedPurchases() error {
	fmt.Println("Seeding purchases...")

	rows, err := s.db.Query("SELECT cms_user_id FROM cms_user WHERE cms_user_role = 'CMS_CUSTOMER'")
	if err != nil {
		return fmt.Errorf("failed to get customers: %w", err)
	}
	defer rows.Close()

	systemNames := []string{
		"Learn Management System",
		"Content Management Pro",
		"EduPortal",
		"CourseBuilder",
		"KnowledgeBase Pro",
		"Academy Platform",
		"Training Hub",
		"SkillMaster",
	}

	count := 0
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			continue
		}

		// Each customer gets 1-3 purchases
		purchaseCount := 1 + rand.Intn(3)
		for i := 0; i < purchaseCount; i++ {
			purchaseDate := time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(30))
			systemName := systemNames[rand.Intn(len(systemNames))]

			_, err = s.db.Exec(`
                INSERT INTO cms_cus_purchase (cms_cus_id, system_name, purchase_date, created_at)
                VALUES ($1, $2, $3, $4)`,
				userID, systemName, purchaseDate, time.Now())
			if err != nil {
				return fmt.Errorf("failed to insert purchase: %w", err)
			}
			count++
		}
	}

	fmt.Printf("Seeded %d purchases\n", count)
	return nil
}

func (s *Seeder) seedPageRequests() error {
	fmt.Println("Seeding user page requests...")

	rows, err := s.db.Query("SELECT cms_user_id FROM cms_user WHERE cms_user_role = 'CMS_CUSTOMER'")
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	pageTypes := []string{
		"LANDING", "BLOG", "PORTFOLIO", "BUSINESS", "PERSONAL",
		"E-COMMERCE", "EDUCATIONAL", "NONPROFIT", "PORTAL", "WIKI",
	}

	statuses := []string{
		"PENDING", "APPROVED", "REJECTED", "REVISION_REQUESTED", "CANCELLED",
	}

	count := 0
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			continue
		}

		// Each user gets 1-2 page requests
		requestCount := 1 + rand.Intn(2)
		for i := 0; i < requestCount; i++ {
			pageType := pageTypes[rand.Intn(len(pageTypes))]
			status := statuses[rand.Intn(len(statuses))]

			// Use ON CONFLICT DO NOTHING to handle the unique constraint
			_, err = s.db.Exec(`
                INSERT INTO user_page_request (user_id, page_type, status, created_at)
                VALUES ($1, $2, $3, $4)
                ON CONFLICT (user_id) DO NOTHING`,
				userID, pageType, status, time.Now().Add(-time.Duration(rand.Intn(30))*24*time.Hour))
			if err != nil {
				return fmt.Errorf("failed to insert page request: %w", err)
			}
			count++
		}
	}

	fmt.Printf("Seeded %d user page requests\n", count)
	return nil
}

func (s *Seeder) seedCMSPageRequests() error {
	fmt.Println("Seeding CMS page requests...")

	// Get all admin IDs
	adminRows, err := s.db.Query("SELECT cms_user_id FROM cms_user WHERE cms_user_role IN ('CMS_ADMIN', 'CMS_EDITOR', 'CMS_SUPPORT')")
	if err != nil {
		return fmt.Errorf("failed to get admins: %w", err)
	}
	defer adminRows.Close()

	var adminIDs []uuid.UUID
	for adminRows.Next() {
		var id uuid.UUID
		if err := adminRows.Scan(&id); err != nil {
			continue
		}
		adminIDs = append(adminIDs, id)
	}

	ownerRows, err := s.db.Query("SELECT cms_user_id FROM cms_user WHERE cms_user_role = 'CMS_CUSTOMER'")
	if err != nil {
		return fmt.Errorf("failed to get owners: %w", err)
	}
	defer ownerRows.Close()

	requestTypes := []string{
		"WEBSITE_CREATION", "BLOG_SETUP", "PORTFOLIO_BUILD", "BUSINESS_PAGE",
		"LANDING_PAGE", "E-COMMERCE_SETUP", "EDUCATIONAL_PORTAL", "NONPROFIT_SITE",
		"WIKI_CREATION", "RESTRUCTURE", "CONTENT_MIGRATION", "DESIGN_REVAMP",
	}

	statuses := []string{
		"PENDING", "IN_PROGRESS", "COMPLETED", "REJECTED", "ON_HOLD",
		"NEEDS_REVIEW", "AWAITING_FEEDBACK", "ARCHIVED",
	}

	titles := []string{
		"New Business Website", "Personal Blog Setup", "Company Portfolio",
		"Product Landing Page", "Online Store Creation", "Educational Platform",
		"Nonprofit Organization Site", "Knowledge Base Wiki", "Corporate Portal",
		"Marketing Campaign Page", "Event Microsite", "Community Forum",
	}

	count := 0
	for ownerRows.Next() {
		var ownerID uuid.UUID
		if err := ownerRows.Scan(&ownerID); err != nil {
			continue
		}

		// Each owner gets 2-5 requests
		requestCount := 2 + rand.Intn(4)
		for i := 0; i < requestCount; i++ {
			requestID := uuid.New()
			requestType := requestTypes[rand.Intn(len(requestTypes))]
			status := statuses[rand.Intn(len(statuses))]
			title := titles[rand.Intn(len(titles))]
			adminID := adminIDs[rand.Intn(len(adminIDs))]

			_, err = s.db.Exec(`
                INSERT INTO cms_page_request (
                    request_id, owner_id, request_type, title, description, 
                    page_url, logo_url, status, admin_id, created_at, updated_at
                ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
				requestID, ownerID, requestType,
				title,
				fmt.Sprintf("Detailed description for %s request. Includes all necessary requirements and specifications.", title),
				fmt.Sprintf("https://%s-site-%d.example.com", requestType, count+1),
				fmt.Sprintf("https://assets.example.com/logos/logo-%d.png", count+1),
				status, adminID,
				time.Now().Add(-time.Duration(rand.Intn(90))*24*time.Hour),
				time.Now().Add(-time.Duration(rand.Intn(30))*24*time.Hour))
			if err != nil {
				return fmt.Errorf("failed to insert CMS page request: %w", err)
			}
			count++
		}
	}

	fmt.Printf("Seeded %d CMS page requests\n", count)
	return nil
}

func (s *Seeder) seedCMSPages() error {
	fmt.Println("Seeding CMS pages...")

	// Get all staff IDs (admins, editors, developers)
	staffRows, err := s.db.Query("SELECT cms_user_id FROM cms_user WHERE cms_user_role IN ('CMS_ADMIN', 'CMS_EDITOR', 'CMS_DEVELOPER', 'CMS_SUPPORT')")
	if err != nil {
		return fmt.Errorf("failed to get staff: %w", err)
	}
	defer staffRows.Close()

	var staffIDs []uuid.UUID
	for staffRows.Next() {
		var id uuid.UUID
		if err := staffRows.Scan(&id); err != nil {
			continue
		}
		staffIDs = append(staffIDs, id)
	}

	rows, err := s.db.Query(`
        SELECT request_id, owner_id 
        FROM cms_page_request 
        WHERE status IN ('COMPLETED', 'IN_PROGRESS', 'NEEDS_REVIEW', 'AWAITING_FEEDBACK')`)
	if err != nil {
		return fmt.Errorf("failed to get page requests: %w", err)
	}
	defer rows.Close()

	statuses := []string{"DRAFT", "PUBLISHED", "ARCHIVED", "PENDING_REVIEW", "SCHEDULED"}

	pageTemplates := []struct {
		title   string
		content string
	}{
		{"Welcome Page", "<h1>Welcome to Our Site</h1><p>This is the homepage content.</p>"},
		{"About Us", "<h1>About Our Company</h1><p>Learn about our history and mission.</p>"},
		{"Services", "<h1>Our Services</h1><ul><li>Service 1</li><li>Service 2</li></ul>"},
		{"Blog", "<h1>Latest Articles</h1><div class='article'>...</div>"},
		{"Contact", "<h1>Get in Touch</h1><form>...</form>"},
		{"Portfolio", "<h1>Our Work</h1><div class='gallery'>...</div>"},
		{"Products", "<h1>Our Products</h1><div class='product-grid'>...</div>"},
		{"FAQ", "<h1>Frequently Asked Questions</h1><div class='accordion'>...</div>"},
	}

	count := 0
	for rows.Next() {
		var requestID, ownerID uuid.UUID
		if err := rows.Scan(&requestID, &ownerID); err != nil {
			continue
		}

		// Each request gets 1-3 pages
		pageCount := 1 + rand.Intn(3)
		for i := 0; i < pageCount; i++ {
			status := statuses[rand.Intn(len(statuses))]
			template := pageTemplates[rand.Intn(len(pageTemplates))]
			staffID := staffIDs[rand.Intn(len(staffIDs))]

			_, err = s.db.Exec(`
    INSERT INTO cms_page (
        page_request_id, title, content, image_url, status, 
        owner_id, published_by_staff_id, created_at, updated_at
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
				requestID,
				fmt.Sprintf("%s %d", template.title, i+1),
				template.content,
				fmt.Sprintf("https://assets.example.com/images/page-%d-%d.jpg", count+1, i+1),
				status, ownerID, staffID,
				time.Now().Add(-time.Duration(rand.Intn(60))*24*time.Hour),
				time.Now().Add(-time.Duration(rand.Intn(15))*24*time.Hour))
			if err != nil {
				return fmt.Errorf("failed to insert CMS page: %w", err)
			}
			count++
		}
	}

	fmt.Printf("Seeded %d CMS pages\n", count)
	return nil
}
