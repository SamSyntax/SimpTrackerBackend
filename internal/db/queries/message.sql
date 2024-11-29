-- name: UpsertUserMessage :exec
INSERT INTO user_messages (user_id, streamer_id, keyword_id, count, last_message, message_date)
VALUES ($1, $2, $3, 1, $4, CURRENT_DATE)
ON CONFLICT (user_id, keyword_id, streamer_id, message_date)
DO UPDATE SET
    count = user_messages.count + 1,
    last_message = EXCLUDED.last_message,
    updated_at = NOW();

-- name: GetMessagesByStreamer :many
SELECT * FROM user_messages
WHERE streamer_id = $1
  AND message_date BETWEEN $2 AND $3;
