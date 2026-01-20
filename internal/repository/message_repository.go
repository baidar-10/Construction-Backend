package repository

import (
	"construction-backend/internal/database"
	"construction-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type MessageRepository struct {
	db *database.Database
}

func NewMessageRepository(db *database.Database) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *MessageRepository) FindBetweenUsers(userID1, userID2 uuid.UUID) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID1, userID2, userID2, userID1).
		Order("created_at ASC").Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) GetConversations(userID uuid.UUID) ([]map[string]interface{}, error) {
	var conversations []map[string]interface{}

	rows, err := r.db.Raw(`
		SELECT DISTINCT ON (other_user_id)
			other_user_id,
			message_id,
			last_message,
			last_message_time,
			unread_count,
			sender_id
		FROM (
			SELECT 
				CASE 
					WHEN sender_id = ? THEN receiver_id 
					ELSE sender_id 
				END as other_user_id,
				id as message_id,
				content as last_message,
				created_at as last_message_time,
				sender_id,
				SUM(CASE WHEN receiver_id = ? AND is_read = false THEN 1 ELSE 0 END) as unread_count
			FROM messages
			WHERE sender_id = ? OR receiver_id = ?
			GROUP BY other_user_id, message_id, content, created_at, sender_id
			ORDER BY other_user_id, created_at DESC
		) sub
		ORDER BY other_user_id, last_message_time DESC
	`, userID, userID, userID, userID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var conv map[string]interface{}
		var otherUserID uuid.UUID
		var messageID uuid.UUID
		var lastMessage string
		var lastMessageTime time.Time
		var unreadCount int
		var senderID uuid.UUID

		if err := rows.Scan(&otherUserID, &messageID, &lastMessage, &lastMessageTime, &unreadCount, &senderID); err != nil {
			continue
		}

		// Get other user details
		var user models.User
		r.db.First(&user, "id = ?", otherUserID)

		isRead := senderID != userID || unreadCount == 0

		conv = map[string]interface{}{
			"userId": otherUserID,
			"otherUser": map[string]interface{}{
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"avatarUrl": user.AvatarURL,
			},
			"lastMessage": map[string]interface{}{
				"id":        messageID,
				"content":   lastMessage,
				"createdAt": lastMessageTime.Format(time.RFC3339),
				"isRead":    isRead,
				"senderID":  senderID,
			},
			"unreadCount": unreadCount,
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (r *MessageRepository) MarkAsRead(messageID uuid.UUID) error {
	return r.db.Model(&models.Message{}).Where("id = ?", messageID).
		Update("is_read", true).Error
}

func (r *MessageRepository) MarkAllAsRead(senderID, receiverID uuid.UUID) error {
	return r.db.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ?", senderID, receiverID).
		Update("is_read", true).Error
}

func (r *MessageRepository) FindByBookingID(bookingID uuid.UUID) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Preload("Sender").Preload("Receiver").
		Where("booking_id = ?", bookingID).
		Order("created_at ASC").Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) MarkBookingMessagesAsRead(bookingID, userID uuid.UUID) error {
	return r.db.Model(&models.Message{}).
		Where("booking_id = ? AND receiver_id = ?", bookingID, userID).
		Update("is_read", true).Error
}
