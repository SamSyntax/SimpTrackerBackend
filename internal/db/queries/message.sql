-- name: UpsertUserMessage :exec
INSERT INTO user_messages (user_id, keyword_id, count, last_message)
VALUES ($1, $2, $3, $4)
  ON CONFLICT (user_id, keyword_id) DO UPDATE
  SET count = user_messages.count + EXCLUDED.count,
  last_message = EXCLUDED.last_message,
  updated_at = NOW();

