#!/bin/bash

# Script to create admin user in the database
# Usage: ./create_admin.sh

echo "🔐 Creating admin user..."

# Database connection details
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="construction_db"
DB_USER="postgres"
DB_PASSWORD="postgres123"

# Admin credentials
ADMIN_EMAIL="admin@stroyhub.com"
ADMIN_PASSWORD="admin123"
ADMIN_NAME="System Administrator"
ADMIN_PHONE="+77001234567"

# Generate bcrypt hash for the password
# Note: You'll need to run this with Go or use an online bcrypt generator
# For now, using a pre-hashed version of 'admin123'
# Generated with: bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
PASSWORD_HASH='$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1lRps.9cGLcZEiGDMVr5yUP1KUOYTa'

# SQL command to create admin user
SQL="
DO \$\$
BEGIN
    -- First, update the constraint to allow 'admin' user type
    IF EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'users_user_type_check'
    ) THEN
        ALTER TABLE users DROP CONSTRAINT users_user_type_check;
    END IF;
    
    ALTER TABLE users ADD CONSTRAINT users_user_type_check 
    CHECK (user_type IN ('customer', 'worker', 'admin'));
    
    -- Insert admin user if not exists
    INSERT INTO users (id, email, password_hash, full_name, phone_number, user_type, is_active, created_at, updated_at)
    VALUES (
        gen_random_uuid(),
        '$ADMIN_EMAIL',
        '$PASSWORD_HASH',
        '$ADMIN_NAME',
        '$ADMIN_PHONE',
        'admin',
        true,
        NOW(),
        NOW()
    )
    ON CONFLICT (email) DO UPDATE SET
        password_hash = EXCLUDED.password_hash,
        user_type = 'admin',
        is_active = true,
        updated_at = NOW();
    
    RAISE NOTICE 'Admin user created/updated successfully!';
END \$\$;
"

# Execute the SQL command
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "$SQL"

if [ $? -eq 0 ]; then
    echo "✅ Admin user created successfully!"
    echo ""
    echo "📝 Admin credentials:"
    echo "   Email: $ADMIN_EMAIL"
    echo "   Password: $ADMIN_PASSWORD"
    echo ""
    echo "🔗 Login at: http://localhost:5173/login"
    echo "🔗 Admin panel: http://localhost:5173/admin"
else
    echo "❌ Failed to create admin user"
    exit 1
fi
