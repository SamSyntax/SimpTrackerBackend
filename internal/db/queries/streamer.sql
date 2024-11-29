-- name: UpsertStreamer :one
INSERT INTO streamers (twitch_id, username, access_token, refresh_token, expires_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (twitch_id) DO UPDATE SET
    username = EXCLUDED.username,
    access_token = EXCLUDED.access_token,
    refresh_token = EXCLUDED.refresh_token,
    expires_at = EXCLUDED.expires_at,
    updated_at = NOW()
RETURNING id;

-- name: GetStreamerByTwitchID :one
SELECT * FROM streamers WHERE twitch_id = $1;

-- name: GetStreamers :many
SELECT * FROM streamers;
