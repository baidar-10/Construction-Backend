#!/bin/bash

echo "=== MIGRATION VERIFICATION REPORT ==="
echo ""

echo "✓ Checking Migration 001 - Booking Messages"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name='messages');"

echo "✓ Checking Migration 002 - Admin Support"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT COUNT(*) FROM users WHERE user_type='admin';"

echo "✓ Checking Migration 003 - Open Bookings"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='bookings' AND column_name='is_open');"

echo "✓ Checking Migration 005 - Booking Applications"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name='booking_applications');"

echo "✓ Checking Migration 006 - User Verification"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='verification_code');"

echo "✓ Checking Migration 007 - Team Members"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name='team_members');"

echo "✓ Checking Migration 008 - Last Login At"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='last_login_at');"

echo "✓ Checking Migration 008 - Worker Promotion"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='workers' AND column_name='is_promoted');"

echo "✓ Checking Migration 009 - Worker Payment Fields"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='workers' AND column_name='payment_type');"

echo "✓ Checking Migration 010 - Promotion Requests"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name='promotion_requests');"

echo "✓ Checking Migration 011 - Review Media"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='reviews' AND column_name='media_urls');"

echo "✓ Checking Migration 011 - Verification Documents"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name='verification_documents');"

echo "✓ Checking Migration 012 - Portfolio"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name='portfolio_items');"

echo "✓ Checking Migration 012 - Rework Status"
docker exec construction_db psql -U admin -d construction_db -t -c "SELECT COUNT(*) FROM pg_enum WHERE enumlabel='rework_required' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'verification_status');"

echo ""
echo "=== ALL DATABASE TABLES ==="
docker exec construction_db psql -U admin -d construction_db -c "\dt"
