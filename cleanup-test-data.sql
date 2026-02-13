-- Cleanup script for test data
DELETE FROM worker_skills WHERE worker_id IN (SELECT id FROM workers WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'worker%@test.com' OR email LIKE 'customer%@test.com'));
DELETE FROM worker_team_members WHERE worker_id IN (SELECT id FROM workers WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'worker%@test.com'));
DELETE FROM reviews WHERE worker_id IN (SELECT id FROM workers WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'worker%@test.com'));
DELETE FROM reviews WHERE customer_id IN (SELECT id FROM users WHERE email LIKE 'customer%@test.com');
DELETE FROM bookings WHERE worker_id IN (SELECT id FROM workers WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'worker%@test.com'));
DELETE FROM bookings WHERE customer_id IN (SELECT id FROM users WHERE email LIKE 'customer%@test.com');
DELETE FROM workers WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'worker%@test.com');
DELETE FROM users WHERE email LIKE 'worker%@test.com' OR email LIKE 'customer%@test.com';
