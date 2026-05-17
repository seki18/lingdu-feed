package repository

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
)

// GetInteractionStatus retrieves a single Status by primary key.
func GetInteractionStatus(interactionStatus model.InteractionStatus) (model.InteractionStatus, error) {
	var status model.InteractionStatus
	err := common.DB.Get(&status, `
		SELECT post_id, user_id, status, updated_time
		FROM interaction_status
		WHERE post_id = $1
		AND user_id = $2
	`, interactionStatus.PostID, interactionStatus.UserID)

	return status, err
}

// GetInteractionStatusByUserID retrieves Status by user ID.
func GetInteractionStatusByUserID(userId int) ([]model.InteractionStatus, error) {
	var statuses []model.InteractionStatus
	err := common.DB.Select(&statuses, `
		SELECT post_id, user_id, status, updated_time
		FROM interaction_status
		WHERE user_id = $1
	`, userId)

	return statuses, err
}

// UpsertInteractionStatus inserts a new Status if it doesn't exist, or updates it if it does.
func UpsertInteractionStatus(interactionStatus model.InteractionStatus) error {

	_, err := common.DB.Exec(`
		INSERT INTO interaction_status (
			user_id,
			post_id,
			status,
			updated_time
		)
		VALUES (
			$1,
			$2,
			$3,
			NOW()
		)
		ON CONFLICT (user_id, post_id)
		DO UPDATE SET
			status = GREATEST(
				interaction_status.status,
				EXCLUDED.status
			),
			updated_time = NOW()
	`,
		interactionStatus.UserID,
		interactionStatus.PostID,
		interactionStatus.Status,
	)

	return err
}

