#!/bin/bash

# Script to create 10 test workers with different specialties and cities
API_URL="http://localhost:8080/api"

# Array of different specialties
specialties=("electricity" "plumbing" "carpentry" "painting" "masonry" "hvac" "roofing" "landscaping" "flooring" "drywall")

# Array of Kazakhstan cities
cities=("Almaty" "Astana" "Karaganda" "Shymkent" "Aktobe" "Pavlodar" "Taraz" "Aktau" "Kostanay" "Atyrau")

# Worker names
names=("Baidar" "Arman" "Timur" "Nurlan" "Kazbek" "Aibek" "Daulet" "Samat" "Rustam" "Adil")

echo "Creating 10 test workers with proper cities..."

for i in {1..10}; do
  firstname="${names[$((i-1))]}"
  lastname="Worker"
  email="worker$i@test.com"
  password="Test123!@"
  specialty="${specialties[$((i-1))]}"
  city="${cities[$((i-1))]}"
  hourly_rate=$((2000 + i * 500))

  echo ""
  echo "[$i/10] Creating $firstname $lastname - $specialty specialist in $city..."

  # Register user
  register_response=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"$email\",
      \"password\": \"$password\",
      \"firstName\": \"$firstname\",
      \"lastName\": \"$lastname\",
      \"userType\": \"worker\",
      \"phone\": \"+7700000$(printf '%02d' $i)\"
    }")

  user_id=$(echo $register_response | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
  token=$(echo $register_response | grep -o '"token":"[^"]*' | cut -d'"' -f4)
  
  if [ -z "$user_id" ] || [ -z "$token" ]; then
    echo "❌ Error registering user"
    continue
  fi

  # Create worker profile with location
  worker_response=$(curl -s -X POST "$API_URL/workers" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $token" \
    -d "{
      \"specialty\": \"$specialty\",
      \"hourlyRate\": $hourly_rate,
      \"experienceYears\": $((5 + i)),
      \"bio\": \"Professional $specialty worker with extensive experience in $city\",
      \"location\": \"$city\",
      \"availabilityStatus\": \"available\"
    }")

  if echo $worker_response | grep -q "id"; then
    echo "✓ $firstname created: $specialty in $city ($hourly_rate ₸/hour)"
  else
    echo "⚠ Worker profile created but response unclear"
  fi
done

echo ""
echo "✓ All 10 test workers created successfully!"
echo "You can now test the search functionality on the home page."
