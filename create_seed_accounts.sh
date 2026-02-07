#!/bin/bash

# Creates admin, worker, and customer accounts (with .kz emails)
# Usage: DB_HOST=localhost DB_PORT=5432 DB_USER=admin DB_PASSWORD=admin123 DB_NAME=construction_db ./create_seed_accounts.sh

set -e

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-admin}
DB_PASSWORD=${DB_PASSWORD:-admin123}
DB_NAME=${DB_NAME:-construction_db}
DB_CONTAINER=${DB_CONTAINER:-construction_db}

ADMIN_EMAIL="admin@stroymaster.kz"
WORKER_EMAIL="anuar.ahmedov@mail.ru"
CUSTOMER_EMAIL="aibek.zhusubekov@mail.ru"
DEFAULT_PASSWORD="admin123"

SQL="
DO \$\$
DECLARE
  admin_id uuid;
  worker_user_id uuid;
  customer_user_id uuid;
  worker_profile_id uuid;
BEGIN
  -- Ensure pgcrypto for bcrypt hashing
  CREATE EXTENSION IF NOT EXISTS pgcrypto;
  -- Allow admin user type
  IF EXISTS (
      SELECT 1 FROM pg_constraint
      WHERE conname = 'users_user_type_check'
  ) THEN
      ALTER TABLE users DROP CONSTRAINT users_user_type_check;
  END IF;
  ALTER TABLE users ADD CONSTRAINT users_user_type_check
    CHECK (user_type IN ('customer', 'worker', 'admin'));

  -- Admin user
  INSERT INTO users (id, email, password_hash, first_name, last_name, phone, user_type, is_active, is_verified, created_at, updated_at)
  VALUES (gen_random_uuid(), '$ADMIN_EMAIL', crypt('$DEFAULT_PASSWORD', gen_salt('bf')), 'Admin', 'User', '+77000000001', 'admin', true, true, NOW(), NOW())
  ON CONFLICT (email) DO UPDATE SET
    password_hash = crypt('$DEFAULT_PASSWORD', gen_salt('bf')),
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    phone = EXCLUDED.phone,
    user_type = 'admin',
    is_active = true,
    is_verified = true,
    updated_at = NOW()
  RETURNING id INTO admin_id;

  -- Worker user
  INSERT INTO users (id, email, password_hash, first_name, last_name, phone, user_type, is_active, is_verified, created_at, updated_at)
  VALUES (gen_random_uuid(), '$WORKER_EMAIL', crypt('$DEFAULT_PASSWORD', gen_salt('bf')), 'Ануар', 'Ахмедов', '+77000000002', 'worker', true, true, NOW(), NOW())
  ON CONFLICT (email) DO UPDATE SET
    password_hash = crypt('$DEFAULT_PASSWORD', gen_salt('bf')),
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    phone = EXCLUDED.phone,
    user_type = 'worker',
    is_active = true,
    is_verified = true,
    updated_at = NOW()
  RETURNING id INTO worker_user_id;

  -- Customer user
  INSERT INTO users (id, email, password_hash, first_name, last_name, phone, user_type, is_active, is_verified, created_at, updated_at)
  VALUES (gen_random_uuid(), '$CUSTOMER_EMAIL', crypt('$DEFAULT_PASSWORD', gen_salt('bf')), 'Айбек', 'Жусубеков', '+77000000003', 'customer', true, true, NOW(), NOW())
  ON CONFLICT (email) DO UPDATE SET
    password_hash = crypt('$DEFAULT_PASSWORD', gen_salt('bf')),
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    phone = EXCLUDED.phone,
    user_type = 'customer',
    is_active = true,
    is_verified = true,
    updated_at = NOW()
  RETURNING id INTO customer_user_id;

  -- Worker profile
  INSERT INTO workers (id, user_id, specialty, hourly_rate, experience_years, bio, location, availability_status, created_at, updated_at)
  VALUES (gen_random_uuid(), worker_user_id, 'Roofing Works', 1200, 5, 'Тәжірибелі шатыр жөндеу маманы', 'Astana', 'available', NOW(), NOW())
  ON CONFLICT (user_id) DO UPDATE SET
    specialty = EXCLUDED.specialty,
    hourly_rate = EXCLUDED.hourly_rate,
    experience_years = EXCLUDED.experience_years,
    bio = EXCLUDED.bio,
    location = EXCLUDED.location,
    availability_status = EXCLUDED.availability_status,
    updated_at = NOW()
  RETURNING id INTO worker_profile_id;

  -- Worker skills
  INSERT INTO worker_skills (id, worker_id, skill, created_at)
  VALUES (gen_random_uuid(), worker_profile_id, 'Roofing Works', NOW())
  ON CONFLICT (worker_id, skill) DO NOTHING;

  -- Customer profile
  INSERT INTO customers (id, user_id, address, city, state, postal_code, created_at, updated_at)
  VALUES (gen_random_uuid(), customer_user_id, 'Абылай хан 1', 'Astana', '', '010000', NOW(), NOW())
  ON CONFLICT (user_id) DO UPDATE SET
    address = EXCLUDED.address,
    city = EXCLUDED.city,
    state = EXCLUDED.state,
    postal_code = EXCLUDED.postal_code,
    updated_at = NOW();
END \$\$;
"

if command -v psql >/dev/null 2>&1; then
  PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL"
else
  if ! command -v podman >/dev/null 2>&1; then
    echo "❌ psql not found and podman not available"
    exit 1
  fi
  podman exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -c "$SQL"
fi

if [ $? -eq 0 ]; then
  echo "✅ Seed accounts created/updated"
  echo "Admin: $ADMIN_EMAIL / $DEFAULT_PASSWORD"
  echo "Worker: $WORKER_EMAIL / $DEFAULT_PASSWORD"
  echo "Customer: $CUSTOMER_EMAIL / $DEFAULT_PASSWORD"
else
  echo "❌ Failed to create seed accounts"
  exit 1
fi
