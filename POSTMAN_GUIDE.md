# Postman Guide for Construction Backend API

## What is Postman?

Postman is an API testing tool that helps you:
- **Test your backend endpoints** without needing a frontend
- **Debug API issues** by inspecting requests and responses
- **Document your API** by saving example requests
- **Share API collections** with your team

## Getting Started

### 1. Launch Postman
- Open Postman from Applications folder (macOS)
- Sign up for a free account (optional but recommended for saving collections)

### 2. Set Up Base URL
Your backend API is running at:
- **Local Development**: `http://localhost:8080`
- **Production Server**: `http://85.202.192.68:8080`

## Testing Your API Endpoints

### Authentication Endpoints

#### 1. Register New User (Customer)
```
Method: POST
URL: http://localhost:8080/api/auth/register
Body (JSON):
{
  "firstName": "John",
  "lastName": "Doe",
  "email": "john@example.com",
  "password": "password123",
  "phone": "+77771234567",
  "location": "Almaty, Kazakhstan",
  "userType": "customer"
}
```

#### 2. Register New Worker
```
Method: POST
URL: http://localhost:8080/api/auth/register
Body (JSON):
{
  "firstName": "Ivan",
  "lastName": "Petrov",
  "email": "ivan@example.com",
  "password": "password123",
  "phone": "+77779876543",
  "location": "Almaty, Kazakhstan",
  "userType": "worker",
  "role": "Electrician",
  "yearsExperience": 5,
  "hourlyRate": 50,
  "skills": "Wiring,Installation,Repair",
  "bio": "Experienced electrician with 5 years of experience"
}
```

#### 3. Login
```
Method: POST
URL: http://localhost:8080/api/auth/login
Body (JSON):
{
  "email": "john@example.com",
  "password": "password123"
}

Response will include a JWT token - save this for authenticated requests!
```

### Worker Endpoints

#### 4. Get All Workers
```
Method: GET
URL: http://localhost:8080/api/workers
Headers: None required (public endpoint)
```

#### 5. Search Workers
```
Method: GET
URL: http://localhost:8080/api/workers/search?query=electrician
```

#### 6. Get Worker by ID
```
Method: GET
URL: http://localhost:8080/api/workers/1
```

### Booking Endpoints (Requires Authentication)

#### 7. Create Booking
```
Method: POST
URL: http://localhost:8080/api/bookings
Headers:
  Authorization: Bearer YOUR_JWT_TOKEN_HERE
Body (JSON):
{
  "workerID": 1,
  "title": "Kitchen Renovation",
  "description": "Need help with kitchen electrical work",
  "date": "2026-01-15",
  "duration": 4,
  "location": "Almaty, Kazakhstan"
}
```

#### 8. Get My Bookings
```
Method: GET
URL: http://localhost:8080/api/bookings
Headers:
  Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

### Message Endpoints (Requires Authentication)

#### 9. Get Conversations
```
Method: GET
URL: http://localhost:8080/api/messages/conversations
Headers:
  Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

#### 10. Send Message
```
Method: POST
URL: http://localhost:8080/api/messages
Headers:
  Authorization: Bearer YOUR_JWT_TOKEN_HERE
Body (JSON):
{
  "receiverID": 2,
  "bookingID": 1,
  "content": "Hi, I'd like to discuss the project details"
}
```

### Review Endpoints

#### 11. Create Review
```
Method: POST
URL: http://localhost:8080/api/reviews
Headers:
  Authorization: Bearer YOUR_JWT_TOKEN_HERE
Body (JSON):
{
  "workerID": 1,
  "bookingID": 1,
  "rating": 5,
  "comment": "Excellent work! Very professional."
}
```

## Tips for Using Postman

### 1. Save Requests in a Collection
- Click "New" → "Collection"
- Name it "Construction Backend API"
- Add all your requests to this collection

### 2. Use Environment Variables
- Create an environment for "Development" and "Production"
- Set variables like:
  - `base_url`: `http://localhost:8080`
  - `token`: Your JWT token
- Use them in requests: `{{base_url}}/api/workers`

### 3. Set Authorization Token Globally
- In your collection settings, go to "Authorization"
- Select "Bearer Token"
- Set the token once for all requests

### 4. Test Response Data
- Go to the "Tests" tab in a request
- Add assertions to verify responses:
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has data", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('data');
});
```

## Common Issues

### CORS Errors
If you see CORS errors, they only happen in the browser. Postman bypasses CORS, so you can test freely.

### 401 Unauthorized
Make sure your JWT token is:
1. Valid (not expired)
2. Properly formatted: `Bearer YOUR_TOKEN_HERE`
3. Included in the Authorization header

### Connection Refused
Make sure your backend server is running:
```bash
cd Construction-Backend
go run cmd/api/main.go
```

## Next Steps

1. **Create a Collection**: Organize all your API requests
2. **Document Examples**: Save successful responses as examples
3. **Share with Team**: Export collection and share the JSON file
4. **Automate Testing**: Use Collection Runner for automated tests

Happy Testing! 🚀
