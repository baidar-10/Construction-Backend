#!/bin/bash

# Script to create admin user in the database
# Usage: DB_HOST=localhost DB_PORT=5432 DB_USER=admin DB_PASSWORD=admin123 DB_NAME=construction_db ./create_admin.sh

echo "🔐 Creating admin user..."

set -e

# Database connection details
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-admin}
DB_PASSWORD=${DB_PASSWORD:-admin123}
DB_NAME=${DB_NAME:-construction_db}
DB_CONTAINER=${DB_CONTAINER:-buildconnect-postgres}

# Admin credentials
ADMIN_EMAIL="admin@stroymaster.com"
ADMIN_PASSWORD="admin123"
ADMIN_FIRST_NAME="System"
ADMIN_LAST_NAME="Administrator"
ADMIN_PHONE="+77001234567"

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
    
    -- Ensure pgcrypto for bcrypt hashing
    CREATE EXTENSION IF NOT EXISTS pgcrypto;

    -- Insert admin user if not exists
    INSERT INTO users (id, email, password_hash, first_name, last_name, phone, user_type, is_active, is_verified, created_at, updated_at)
    VALUES (
        uuid_generate_v4(),
        '$ADMIN_EMAIL',
        crypt('$ADMIN_PASSWORD', gen_salt('bf')),
        '$ADMIN_FIRST_NAME',
        '$ADMIN_LAST_NAME',
        '$ADMIN_PHONE',
        'admin',
        true,
        true,
        NOW(),
        NOW()
    )
    ON CONFLICT (email) DO UPDATE SET
        password_hash = crypt('$ADMIN_PASSWORD', gen_salt('bf')),
        first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        phone = EXCLUDED.phone,
        user_type = 'admin',
        is_active = true,
        is_verified = true,
        updated_at = NOW();
    
    RAISE NOTICE 'Admin user created/updated successfully!';
END \$\$;
"

run_psql() {
    if command -v psql >/dev/null 2>&1; then
        PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL"
        return $?
    fi

    if command -v docker-compose >/dev/null 2>&1; then
        docker-compose exec -T postgres env PGPASSWORD="$DB_PASSWORD" psql -U "$DB_USER" -d "$DB_NAME" -c "$SQL"
        return $?
    fi

    if command -v docker >/dev/null 2>&1; then
        docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -c "$SQL"
        return $?
    fi

    if command -v podman >/dev/null 2>&1; then
        podman exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -c "$SQL"
        return $?
    fi

    echo "❌ psql not found and no container runtime available"
    return 1
}

run_psql

if [ $? -eq 0 ]; then
    echo "✅ Admin user created successfully!"
    echo ""
    echo "📝 Admin credentials:"
    echo "   Email: $ADMIN_EMAIL"
    echo "   Password is set"
    echo ""
    echo "🔗 Login at: http://$DB_HOST:5173/login"
    echo "🔗 Admin panel: http://$DB_HOST:5173/admin"
else
    echo "❌ Failed to create admin user"
    exit 1
fi
