
CREATE  SCHEMA  IF NOT EXISTS  content_management_system;
SET  search_path  to content_management_system,public;
-- Drop existing tables if they exist
DROP TABLE IF EXISTS cms_cus_purchase;
DROP TABLE IF EXISTS cms_whole_sys_role;
DROP TABLE IF EXISTS cms_user;

-- Create extension for cryptographic functions
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create custom types
CREATE TYPE role_type AS ENUM ('ROOT_ADMIN', 'CMS_CUSTOMER');
CREATE TYPE system_type AS ENUM ('LMS', 'EMS');

-- Create tables
CREATE TABLE cms_whole_sys_role (
                                    role_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
                                    role_name role_type NOT NULL UNIQUE
);

CREATE TABLE cms_user (
                          cms_user_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
                          cms_user_name varchar(100) NOT NULL,
                          cms_user_email varchar(150) NOT NULL UNIQUE,
                          cms_name_space varchar(100),
                          password varchar(90) NOT NULL,
                          cms_user_role_id UUID NOT NULL,
                          verified bool DEFAULT FALSE,
                          created_at timestamp DEFAULT CURRENT_TIMESTAMP,
                          updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
                          CONSTRAINT fk_cms_user_role
                              FOREIGN KEY (cms_user_role_id)
                                  REFERENCES cms_whole_sys_role(role_id) ON DELETE CASCADE
);

CREATE TABLE cms_cus_purchase (
                                  relation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                  cms_cus_id UUID NOT NULL,
                                  system_name system_type NOT NULL,
                                  purchase_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                  CONSTRAINT fk_cms_purchase_user
                                      FOREIGN KEY (cms_cus_id)
                                          REFERENCES cms_user(cms_user_id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_cms_user_email ON cms_user USING HASH (cms_user_email);
CREATE INDEX idx_cms_user_id ON cms_user USING HASH (cms_user_id);
CREATE INDEX idx_role_id ON cms_whole_sys_role USING HASH (role_id);
CREATE INDEX idx_cms_user_role ON cms_user(cms_user_role_id);
CREATE INDEX idx_purchase_user ON cms_cus_purchase(cms_cus_id);
CREATE INDEX idx_purchase_date ON cms_cus_purchase(purchase_date);
CREATE INDEX idx_purchase_system ON cms_cus_purchase USING HASH (system_name);
CREATE INDEX idx_cms_user_role_lookup ON cms_user(cms_user_role_id);
CREATE INDEX idx_role_name_filter ON cms_whole_sys_role(role_name);


CREATE INDEX idx_cms_user_role_id ON cms_user(cms_user_role_id);
CREATE INDEX idx_cms_cus_purchase_customer ON cms_cus_purchase(cms_cus_id);
CREATE INDEX idx_role_name ON cms_whole_sys_role(role_name);
-- Insert initial role data
INSERT INTO cms_whole_sys_role (role_name)
VALUES
    ('ROOT_ADMIN'),
    ('CMS_CUSTOMER');

-- Insert admin users
INSERT INTO cms_user (cms_user_name, cms_user_email, password, verified, cms_user_role_id)
VALUES
    ('Super Admin', 'superadmin@company.com',
     crypt('AdminPassword123!', gen_salt('bf')), TRUE,
     (SELECT role_id FROM cms_whole_sys_role WHERE role_name = 'ROOT_ADMIN')),
    ('System Admin', 'sysadmin@company.com',
     crypt('AdminPassword456!', gen_salt('bf')), TRUE,
     (SELECT role_id FROM cms_whole_sys_role WHERE role_name = 'ROOT_ADMIN'));

-- Insert customer data and purchases
WITH new_customers AS (
    INSERT INTO cms_user (
                          cms_user_name,
                          cms_user_email,
                          cms_name_space,
                          password,
                          cms_user_role_id,
                          verified
        )
        SELECT
            customer_data.name,
            customer_data.email,
            customer_data.namespace,
            crypt(customer_data.password, gen_salt('bf')),
            (SELECT role_id FROM cms_whole_sys_role WHERE role_name = 'CMS_CUSTOMER'),
            TRUE
        FROM (VALUES
                  ('John Smith', 'john.smith@company.com', 'john_workspace', 'SecurePassword123!'),
                  ('Sarah Johnson', 'sarah.johnson@company.com', 'sarah_workspace', 'AnotherSecurePass456!')
             ) AS customer_data(name, email, namespace, password)
        RETURNING cms_user_id, cms_user_email
),
     purchase_insert AS (
         INSERT INTO cms_cus_purchase (
                                       cms_cus_id,
                                       system_name,
                                       purchase_date
             )
             SELECT
                 nc.cms_user_id,
                 'LMS',
                 CURRENT_TIMESTAMP
             FROM new_customers nc
             WHERE nc.cms_user_email = 'john.smith@company.com'
             RETURNING *
     )
SELECT 'Customers created and purchase made' as result;

-- Create MFA token table
CREATE TABLE mfa_token (
                           token_id SERIAL PRIMARY KEY,
                           mfa_token VARCHAR NOT NULL,
                           user_id UUID NOT NULL,
                           expires_at TIMESTAMP,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           FOREIGN KEY (user_id) REFERENCES cms_user(cms_user_id) ON DELETE CASCADE
);

-- Create indexes for MFA token table
CREATE INDEX idx_mfa_token_user ON mfa_token(user_id);
CREATE INDEX idx_mfa_token_expires ON mfa_token(expires_at);

-- Query to verify customer data
SELECT
    u.cms_user_name,
    u.cms_user_email,
    u.cms_name_space,
    r.role_name,
    u.verified,
    u.created_at
FROM cms_user u
         JOIN cms_whole_sys_role r ON u.cms_user_role_id = r.role_id
WHERE r.role_name = 'CMS_CUSTOMER'
ORDER BY u.created_at DESC;

-- Query to view all users
SELECT * FROM cms_user;

-- Query to view purchase data
SELECT
    p.relation_id,
    u.cms_user_name,
    u.cms_user_email,
    p.system_name,
    p.purchase_date
FROM cms_cus_purchase p
         JOIN cms_user u ON p.cms_cus_id = u.cms_user_id
ORDER BY p.purchase_date DESC;

-- Example UUID comparison
SELECT 'd69e2e9e-65ea-47e7-9fe7-7f5b661f069b'::uuid = 'd69e2e9e-65ea-47e7-9fe7-7f5b661f069b'::uuid;