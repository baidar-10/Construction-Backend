# Chat and Booking Management Features

## Overview
This update adds comprehensive booking management and chat functionality to the Construction Platform, allowing workers to accept/decline booking requests and enabling real-time communication between workers and customers.

## New Features

### 1. Accept/Decline Booking Requests (Worker)
Workers can now:
- View all booking requests in their dashboard
- Accept pending booking requests
- Decline pending booking requests
- Only pending bookings show action buttons

### 2. Booking-Associated Chat System
- Real-time chat between workers and customers
- Chat becomes available once a booking is accepted
- All messages are stored in the database with booking associations
- Message history is preserved for each booking
- Auto-polling for new messages every 3 seconds

### 3. Enhanced Message Management
- Messages can be associated with specific bookings
- Chat history is accessible for accepted bookings
- Mark messages as read functionality
- Conversation list shows latest messages

## Backend Changes

### API Endpoints

#### Booking Endpoints
```
PUT /api/bookings/:id/accept    - Accept a booking (worker only)
PUT /api/bookings/:id/decline   - Decline a booking (worker only)
```

#### Message Endpoints
```
POST /api/messages                           - Send message (with optional bookingId)
GET /api/messages/booking/:bookingId         - Get messages for a specific booking
PATCH /api/messages/booking/:bookingId/read  - Mark all booking messages as read
```

### Database Changes

#### Messages Table
- Added `booking_id` column (UUID, nullable, foreign key to bookings)
- Added indexes for better query performance

#### Bookings Table
- Updated status constraint to include 'accepted' and 'declined' states
- Status flow: pending → accepted/declined → confirmed → in_progress → completed

### Models Updated
- `Message` model now includes optional `BookingID` field
- Booking status now supports: pending, accepted, declined, confirmed, in_progress, completed, cancelled

## Frontend Changes

### New Components
1. **ChatWindow.jsx** - Main chat interface component
   - Real-time message display
   - Message sending functionality
   - Auto-scroll to latest messages
   - Polling for new messages

2. **ChatMessage.jsx** - Individual message display component
   - Different styling for sent/received messages
   - Timestamp display
   - Sender information

### Updated Components

#### WorkerDashboard.jsx
- Added accept/decline booking handlers
- Integrated chat window for accepted bookings
- Real-time booking list updates

#### CustomerDashboard.jsx
- Added chat functionality for accepted bookings
- Can open chat with workers on accepted bookings

#### BookingCard.jsx
- Added "Open Chat" button for accepted bookings
- Accept/Decline buttons for pending bookings (worker view)
- Status-based button visibility

### New Services

#### messageService.js
```javascript
- sendMessage(receiverId, content, bookingId)
- getMessages(userId)
- getBookingMessages(bookingId)
- getConversations()
- markAsRead(messageId)
- markBookingMessagesAsRead(bookingId)
```

#### Updated bookingService.js
```javascript
- acceptBooking(bookingId)
- declineBooking(bookingId)
```

## Installation & Setup

### Database Migration
If you have an existing database, run the migration:
```bash
psql -U admin -d construction_db -f migrations/001_add_booking_messages.sql
```

Or if using Docker:
```bash
docker exec -i construction_db psql -U admin -d construction_db < migrations/001_add_booking_messages.sql
```

For new installations, the schema is automatically created from init.sql.

### Rebuild and Restart
```bash
# Backend
cd Construction-Backend
docker-compose down
docker-compose up --build -d

# Frontend
cd Construction-Frontend
npm install
npm run dev
```

## Usage Flow

### For Workers:
1. Log in to worker dashboard
2. View pending booking requests
3. Click "Accept" or "Decline" on pending bookings
4. For accepted bookings, click "Open Chat" to communicate with customer
5. Send and receive messages in real-time

### For Customers:
1. Create a booking request for a worker
2. Wait for worker to accept/decline
3. Once accepted, "Open Chat" button appears
4. Click to start conversation with the worker
5. Discuss project details via chat

## Translation Keys Added

```json
{
  "booking": {
    "openChat": "Open Chat",
    "acceptSuccess": "Booking accepted successfully!",
    "declineSuccess": "Booking declined successfully.",
    "confirmDecline": "Are you sure you want to decline this booking?",
    "status": {
      "accepted": "Accepted",
      "declined": "Declined"
    }
  },
  "chat": {
    "title": "Chat",
    "noMessages": "No messages yet. Start the conversation!",
    "typePlaceholder": "Type a message...",
    "send": "Send",
    "sending": "Sending..."
  }
}
```

## Security Considerations

- All endpoints require authentication via JWT
- Workers can only accept/decline their own bookings
- Users can only send messages to users involved in their bookings
- Message access is controlled per user ID

## Future Enhancements

- WebSocket integration for true real-time messaging (replacing polling)
- Push notifications for new messages
- File/image sharing in chat
- Typing indicators
- Read receipts
- Message search functionality
- Chat history export

## Troubleshooting

### Messages not appearing
- Check network tab for API errors
- Verify booking status is "accepted"
- Check browser console for JavaScript errors

### Accept/Decline buttons not working
- Verify user is logged in as a worker
- Check booking status is "pending"
- Check backend logs for authorization errors

### Chat not opening
- Ensure booking status is "accepted"
- Verify both users (worker and customer) exist
- Check browser console for errors
