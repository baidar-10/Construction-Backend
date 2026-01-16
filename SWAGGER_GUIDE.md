# Swagger API Documentation Guide

## 🎉 Swagger is Now Installed!

Swagger provides an interactive web interface where you can **view and test all your API endpoints** directly in your browser.

## 🚀 Access Swagger UI

### Local Development:
```
http://localhost:8080/swagger/index.html
```

### Production Server:
```
http://85.202.192.68:8080/swagger/index.html
```

## 📖 What You'll See in Swagger

Swagger UI shows:
- ✅ All API endpoints organized by category (Auth, Workers, Bookings, etc.)
- ✅ Request/Response formats with examples
- ✅ Required vs optional parameters
- ✅ Authentication requirements
- ✅ **"Try it out" button** to test APIs directly!

## 🧪 How to Test APIs in Swagger

### 1. Testing Public Endpoints (No Auth Required)

**Example: Get All Workers**
1. Open Swagger UI at `http://localhost:8080/swagger/index.html`
2. Find the **`GET /api/workers`** endpoint
3. Click **"Try it out"**
4. Click **"Execute"**
5. See the response with all workers!

### 2. Testing Protected Endpoints (Auth Required)

**Example: Create a Booking**
1. First, login to get a JWT token:
   - Find **`POST /api/auth/login`**
   - Click **"Try it out"**
   - Enter credentials:
   ```json
   {
     "email": "beka@inbox.com",
     "password": "your_password"
   }
   ```
   - Click **"Execute"**
   - Copy the `token` from the response

2. Authorize Swagger:
   - Click the **🔓 Authorize** button at the top right
   - Enter: `Bearer YOUR_TOKEN_HERE`
   - Click **"Authorize"**
   - Click **"Close"**

3. Now test protected endpoints:
   - Find **`POST /api/bookings`**
   - Click **"Try it out"**
   - Enter booking details
   - Click **"Execute"**
   - It will work because you're authenticated!

## 📝 Updating Swagger Documentation

Whenever you add new endpoints or modify existing ones:

```bash
cd /Users/yerlanbarabashkin/Work/Construction-Backend
./generate-swagger.sh
docker-compose up -d --build
```

## 🔍 Swagger vs Postman

| Feature | Swagger | Postman |
|---------|---------|---------|
| **Documentation** | ✅ Auto-generated, always up-to-date | ❌ Manual documentation |
| **Share with team** | ✅ Just share URL | ⚠️ Export/import collections |
| **Test APIs** | ✅ In-browser testing | ✅ Desktop app |
| **Save requests** | ❌ Can't save | ✅ Collections |
| **Best for** | Documentation & quick tests | Deep testing & automation |

**Pro Tip:** Use Swagger for documentation and quick tests, Postman for detailed testing!

## 🌐 Production Deployment

To deploy Swagger to your production server:

```bash
# On your local machine
cd /Users/yerlanbarabashkin/Work/Construction-Backend
./generate-swagger.sh

# SSH to production
ssh ubuntu@85.202.192.68

# On production server
cd ~/buildconnect/backend

# Copy updated files (run this from local machine)
scp -r /Users/yerlanbarabashkin/Work/Construction-Backend/docs ubuntu@85.202.192.68:~/buildconnect/backend/
scp /Users/yerlanbarabashkin/Work/Construction-Backend/cmd/api/main.go ubuntu@85.202.192.68:~/buildconnect/backend/cmd/api/
scp /Users/yerlanbarabashkin/Work/Construction-Backend/Dockerfile ubuntu@85.202.192.68:~/buildconnect/backend/

# Then rebuild on server
docker-compose down
docker-compose up -d --build
```

## 🎨 Swagger UI Features

1. **Schemas Section** - See all data models (User, Worker, Booking, etc.)
2. **Grouping** - Endpoints organized by category
3. **Request Examples** - Pre-filled with sample data
4. **Response Codes** - See all possible responses (200, 400, 401, etc.)
5. **Download Spec** - Export OpenAPI specification

## 🔧 Troubleshooting

### Swagger UI shows "Failed to load API definition"
- Regenerate docs: `./generate-swagger.sh`
- Rebuild containers: `docker-compose up -d --build`

### Can't authenticate in Swagger
- Make sure to include "Bearer " prefix
- Token format: `Bearer eyJhbGciOiJIUzI1NiIs...`

### Endpoints not showing in Swagger
- Make sure you regenerated docs after adding endpoints
- Check that docs folder exists and has swagger.json

## 📚 Next Steps

1. ✅ Open `http://localhost:8080/swagger/index.html`
2. ✅ Test the GET /api/workers endpoint
3. ✅ Login and get a JWT token
4. ✅ Authorize Swagger with the token
5. ✅ Test protected endpoints like bookings

Happy API testing! 🚀
