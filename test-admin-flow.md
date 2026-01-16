# 🧪 Admin Panel Access Guide

## How to Access Admin Panel on Frontend

### Step 1: Open Login Page
Navigate to: **http://localhost:80/login**

### Step 2: Login with Admin Credentials
```
Email:    admin@stroyhub.com
Password: admin123
```

### Step 3: After Login
You should automatically see:
- **👑 Admin** link in the navigation bar (top right, in red color)
- Click it to access the admin dashboard

### Direct Access
You can also go directly to: **http://localhost:80/admin**
(but you must be logged in as admin first)

---

## 🎨 What You'll See

### Desktop Navigation (after admin login):
- Find Workers
- About
- 👑 Admin (in red) ← Click this!
- Messages
- Profile Menu

### Admin Dashboard Features:
1. **Dashboard Tab** - Statistics overview
   - Total Users
   - Total Workers
   - Total Bookings
   - Active Bookings
   - Total Customers
   - Total Reviews

2. **Users Tab** - User management
   - View all users with pagination
   - Filter by type (customers/workers/all)
   - Toggle user status (activate/deactivate)
   - Delete users
   - Search functionality

3. **Bookings Tab** - Booking management
   - View all bookings
   - Filter by status
   - See detailed booking information
   - View customer and worker details

---

## 🔍 Troubleshooting

### "I don't see the Admin link"
✅ **Solution:** Make sure you're logged in with admin credentials
- Email: `admin@stroyhub.com`
- Password: `admin123`

### "I get redirected or see errors"
✅ **Solution:** Check that:
1. Backend is running: `docker-compose ps` (in Construction-Backend)
2. All services show "Up" status
3. Your browser console for any errors (F12)

### "Frontend doesn't load"
✅ **Solution:** 
```bash
cd /Users/yerlanbarabashkin/Work/Construction-Backend
docker-compose restart frontend
```

---

## 📱 Mobile View
On mobile, tap the menu icon (☰) to see:
- 👑 Admin Panel (appears only when logged in as admin)

---

## 🎯 Test Flow

1. Open browser: http://localhost:80/login
2. Enter admin credentials
3. Click "Login"
4. You should be redirected to homepage
5. Look at the navigation bar
6. Click **👑 Admin** (red text)
7. You're now in the admin dashboard!

---

## 🖼️ Visual Guide

### Before Login:
```
[Logo] [Find Workers] [About]          [Language] [Login] [Sign Up]
```

### After Admin Login:
```
[Logo] [Find Workers] [About] [👑 Admin] [Messages]   [Language] [Profile ▼]
                              ↑
                         Click here!
```

---

## 🚀 Quick Commands

### Restart All Services:
```bash
cd /Users/yerlanbarabashkin/Work/Construction-Backend
docker-compose restart
```

### Check Services Status:
```bash
docker-compose ps
```

### View Frontend Logs:
```bash
docker-compose logs -f frontend
```

---

**Everything is already implemented and working! Just log in with admin credentials to see the admin panel.** 🎉
