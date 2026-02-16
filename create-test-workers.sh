#!/bin/bash

# Script to create 10 test workers and 10 test customers
API_URL="http://localhost:8080/api"

echo "========================================"
echo "Cleaning up existing test accounts..."
echo "========================================"

# Delete users created by this script based on email patterns
docker exec construction_db psql -U admin -d construction_db << EOF
DELETE FROM users WHERE email ~ '^(Baidar|Arman|Timur|Nurlan|Kazbek|Aibek|Daulet|Samat|Rustam|Adil|Aigerim|Madina|Saltanat|Zhanna|Karlygash|Gaukhar|Kamila|Zarina|Symbat|Asel)[0-9]+@(inbox\.com|gmail\.com|gmail\.kz)$';
EOF

# Run cleanup SQL script
docker exec construction_db psql -U admin -d construction_db -f /app/cleanup-test-data.sql > /dev/null 2>&1

echo "✅ Cleanup complete"
echo ""

# Array of different specialties
specialties=("Electrical Works" "Plumbing" "Apartment Renovation" "Painting Works" "Roofing Works" "Solar Panel Installation" "Flooring Works" "Furniture Maker" "Designer" "High-Altitude Works")

# Array of Kazakhstan cities (expanded for more variety)
cities=("Almaty" "Astana" "Karaganda" "Shymkent" "Aktobe" "Pavlodar" "Taraz" "Aktau" "Kostanay" "Atyrau" "Semey" "Oral" "Petropavl" "Temirtau" "Turkestan" "Kyzylorda" "Rudny" "Zhezkazgan" "Taldykorgan" "Ust-Kamenogorsk")

# Hourly rates with more variation
hourly_rates=(1500 2500 3000 2200 4000 1800 3500 2800 5000 3200)

# Payment types for variety
payment_types=("hourly" "hourly" "m2" "hourly" "project" "hourly" "m2" "hourly" "project" "hourly")

# Worker names
worker_names=("Baidar" "Arman" "Timur" "Nurlan" "Kazbek" "Aibek" "Daulet" "Samat" "Rustam" "Adil")

# Worker surnames (different for each)
worker_surnames=("Bekbayev" "Tokayev" "Shaimardanov" "Kozhakhmetov" "Assayev" "Mukhanov" "Tarasov" "Orynbayev" "Galiakhmetov" "Dossaliyev")

# Customer names
customer_names=("Aigerim" "Madina" "Saltanat" "Zhanna" "Karlygash" "Gaukhar" "Kamila" "Zarina" "Symbat" "Asel")

# Customer surnames (different for each)
customer_surnames=("Akimova" "Suleimenova" "Ismailova" "Khasanova" "Ospanova" "Erimbetova" "Mukasheva" "Nurgaliyeva" "Bekbayeva" "Saduakassova")

# Arrays to store created IDs
declare -a worker_ids
declare -a worker_tokens
declare -a customer_ids
declare -a customer_tokens

echo "========================================"
echo "Creating 10 test workers..."
echo "========================================"

for i in {1..10}; do
  firstname="${worker_names[$((i-1))]}"
  lastname="${worker_surnames[$((i-1))]}"
  email_domain=$([ $((i % 2)) -eq 0 ] && echo "inbox.com" || echo "gmail.com")
  email="${worker_names[$((i-1))]}$(($RANDOM % 1000))@$email_domain"
  password="qwerty123"
  specialty="${specialties[$((i-1))]}"
  city="${cities[$((i-1))]}"
  hourly_rate="${hourly_rates[$((i-1))]}"
  payment_type="${payment_types[$((i-1))]}"

  echo ""
  echo "[$i/10] Creating worker: $firstname $lastname - $specialty in $city..."

  # Register user with worker profile
  register_response=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"$email\",
      \"password\": \"$password\",
      \"firstName\": \"$firstname\",
      \"lastName\": \"$lastname\",
      \"userType\": \"worker\",
      \"phone\": \"+7700100$(printf '%02d' $i)00\",
      \"location\": \"$city\",
      \"specialty\": \"$specialty\",
      \"hourlyRate\": $hourly_rate,
      \"paymentType\": \"$payment_type\",
      \"currency\": \"KZT\",
      \"experienceYears\": $((3 + i)),
      \"bio\": \"Professional $specialty specialist with over $((3 + i)) years of experience in $city. Quality work guaranteed.\",
      \"availabilityStatus\": \"available\"
    }")

  user_id=$(echo $register_response | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
  
  if [ -z "$user_id" ]; then
    echo "❌ Error registering worker $firstname"
    echo "Response: $register_response"
    continue
  fi
  
  # Login to get token
  login_response=$(curl -s -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"$email\",
      \"password\": \"$password\"
    }")
  
  token=$(echo $login_response | grep -o '"token":"[^"]*' | cut -d'"' -f4)
  
  if [ -z "$token" ]; then
    echo "❌ Error getting token for worker $firstname"
    continue
  fi
  
  # Get worker profile ID
  worker_profile=$(curl -s -X GET "$API_URL/workers/user/$user_id" \
    -H "Authorization: Bearer $token")
  
  worker_id=$(echo $worker_profile | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
  
  if [ -z "$worker_id" ]; then
    echo "❌ Error getting worker profile ID for $firstname"
    echo "   User ID: $user_id"
    echo "   Response: $worker_profile"
    continue
  fi
  
  # Store worker ID and token
  worker_ids[$((i-1))]=$worker_id
  worker_tokens[$((i-1))]=$token
  
  echo "✅ Worker created: $firstname $lastname - $specialty in $city ($hourly_rate ₸/$payment_type) [ID: $worker_id]"
done

echo ""
echo "========================================"
echo "Creating 10 test customers..."
echo "========================================"

for i in {1..10}; do
  firstname="${customer_names[$((i-1))]}"
  lastname="${customer_surnames[$((i-1))]}"
  email="${customer_names[$((i-1))]}$(($RANDOM % 1000))@gmail.kz"
  password="qwerty123"
  # Use different cities for customers (offset by 10 to get different cities than workers)
  city_index=$(((i + 9) % 20))
  city="${cities[$city_index]}"

  echo ""
  echo "[$i/10] Creating customer: $firstname $lastname from $city..."

  # Register customer
  register_response=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"$email\",
      \"password\": \"$password\",
      \"firstName\": \"$firstname\",
      \"lastName\": \"$lastname\",
      \"userType\": \"customer\",
      \"phone\": \"+7700200$(printf '%02d' $i)00\",
      \"location\": \"$city\"
    }")

  user_id=$(echo $register_response | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
  
  if [ -z "$user_id" ]; then
    echo "❌ Error registering customer $firstname"
    echo "Response: $register_response"
    continue
  fi
  
  # Login to get token
  login_response=$(curl -s -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"$email\",
      \"password\": \"$password\"
    }")
  
  token=$(echo $login_response | grep -o '"token":"[^"]*' | cut -d'"' -f4)
  
  # Get customer profile ID
  customer_profile=$(curl -s -X GET "$API_URL/customers/user/$user_id" \
    -H "Authorization: Bearer $token")
  
  customer_id=$(echo $customer_profile | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
  
  if [ -z "$customer_id" ]; then
    echo "❌ Error getting customer profile ID for $firstname"
    echo "   User ID: $user_id"
    echo "   Response: $customer_profile"
    continue
  fi
  
  # Store customer ID and token
  customer_ids[$((i-1))]=$customer_id
  customer_tokens[$((i-1))]=$token

  echo "✅ Customer created: $firstname $lastname from $city"
done

echo ""
echo "========================================"
echo "Creating bookings and reviews..."
echo "========================================"

# Review comments with variety
review_comments=(
  "Excellent work! Very professional and finished on time. Highly recommend!"
  "Good quality work, but took a bit longer than expected. Overall satisfied."
  "Outstanding service! Attention to detail was impressive. Will hire again!"
  "Decent work, met basic requirements. Communication could be better."
  "Amazing craftsmanship! Exceeded all expectations. Worth every penny!"
  "Satisfactory work. Got the job done but nothing exceptional."
  "Fantastic experience! Professional, clean, and efficient. Top quality!"
  "Work was okay. Some minor issues but resolved them promptly."
  "Superb quality! Very skilled and knowledgeable. Best in the business!"
  "Average work. Did what was asked but lacked creativity."
)

# Ratings (1-5 stars) corresponding to each comment
ratings=(5 3 5 3 5 3 5 3 5 3)

# Create bookings and reviews for each worker
for worker_idx in {0..9}; do
  worker_id="${worker_ids[$worker_idx]}"
  worker_name="${worker_names[$worker_idx]}"
  
  # Each worker gets 1-3 reviews from different customers
  num_reviews=$((1 + RANDOM % 3))
  
  echo ""
  echo "Creating reviews for worker: $worker_name..."
  
  for review_num in $(seq 1 $num_reviews); do
    # Select a random customer (but ensure variety)
    customer_idx=$(((worker_idx + review_num - 1) % 10))
    customer_id="${customer_ids[$customer_idx]}"
    customer_token="${customer_tokens[$customer_idx]}"
    customer_name="${customer_names[$customer_idx]}"
    
    # Create a booking first
    city_idx=$(( (customer_idx + 10) % 20 ))
    booking_location="${cities[$city_idx]}"
    
    booking_response=$(curl -s -X POST "$API_URL/bookings" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $customer_token" \
      -d "{
        \"workerID\": \"$worker_id\",
        \"title\": \"${specialties[$worker_idx]} project\",
        \"description\": \"Need quality ${specialties[$worker_idx]} work done\",
        \"location\": \"$booking_location\",
        \"startDate\": \"2026-02-01T09:00:00Z\",
        \"duration\": 8,
        \"status\": \"completed\"
      }")
    
    booking_id=$(echo $booking_response | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
    
    if [ -z "$booking_id" ]; then
      echo "⚠️  Could not create booking for review"
      echo "     Worker ID: $worker_id"
      echo "     Customer ID: $customer_id"
      echo "     Response: $booking_response"
      continue
    fi
    
    # Use different rating and comment for variety
    comment_idx=$(((worker_idx + review_num - 1) % 10))
    rating="${ratings[$comment_idx]}"
    comment="${review_comments[$comment_idx]}"
    
    # Create review
    review_response=$(curl -s -X POST "$API_URL/workers/$worker_id/reviews" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $customer_token" \
      -d "{
        \"bookingId\": \"$booking_id\",
        \"customerId\": \"$customer_id\",
        \"rating\": $rating,
        \"comment\": \"$comment\"
      }")
    
    if echo $review_response | grep -q "review"; then
      echo "  ✅ Review #$review_num: $rating⭐ from $customer_name"
    else
      echo "  ⚠️  Review creation unclear: $review_response"
    fi
  done
done

echo ""
echo "========================================"
echo "✅ All data created successfully!"
echo "========================================"
echo "Created:"
echo "  - 10 worker accounts (worker1@test.com - worker10@test.com)"
echo "  - 10 customer accounts (customer1@test.com - customer10@test.com)"
echo ""
echo "All accounts use password: qwerty123"
echo ""
echo "You can now test the platform with these accounts!"
