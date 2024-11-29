-- name: UpsertUser :one
INSERT INTO users (username)
VALUES ($1)
ON CONFLICT (username) DO NOTHING
RETURNING id;

-- name: GetUsersWithTotalCounts :many
SELECT
  u.id AS user_id,
  u.username,
  json_build_object(
    'keywords', json_agg(
      json_build_object(
        'keyword_id', k.id,
        'keyword', k.keyword,
        'count', um.count
      )
    ),
    'total_count', SUM(um.count),
    'last_message', MAX(um.last_message),
    'fav_word', 
    (
      SELECT k2.keyword
      FROM user_messages um2
      JOIN keywords k2 ON um2.keyword_id = k2.id
      WHERE um2.user_id = u.id
      ORDER BY um2.count DESC
      LIMIT 1
    )
  ) AS stats
FROM users u
JOIN user_messages um ON um.user_id = u.id
JOIN keywords k ON um.keyword_id = k.id
GROUP BY u.id, u.username
ORDER BY SUM(um.count) DESC;

-- name: GetCountsPerUserPerKeywordByStreamer :many
SELECT
  u.id AS user_id,
  u.username,
  json_build_object(
    'keywords', json_agg(
      json_build_object(
        'keyword_id', k.id,
        'keyword', k.keyword,
        'count', um.count
      )
    ),
    'total_count', SUM(um.count),
    'last_message', MAX(um.last_message),
    'fav_word', 
    (
      SELECT k2.keyword
      FROM user_messages um2
      JOIN keywords k2 ON um2.keyword_id = k2.id
      WHERE um2.user_id = u.id AND um2.streamer_id = $1
      ORDER BY um2.count DESC
      LIMIT 1
    )
  ) AS stats
FROM users u
JOIN user_messages um ON um.user_id = u.id
JOIN keywords k ON um.keyword_id = k.id
WHERE um.streamer_id = $1
GROUP BY u.id, u.username
ORDER BY SUM(um.count) DESC;


-- name: GetCountsPerUserPerKeywordByUsername :one
SELECT
  u.id AS user_id,
  u.username,
  json_build_object(
    'keywords', json_agg(
      json_build_object(
        'keyword_id', k.id,
        'keyword', k.keyword,
        'count', um.count
      )
    ),
    'total_count', SUM(um.count),
    'last_message', MAX(um.last_message),
    'fav_word', 
    (
      SELECT k2.keyword
      FROM user_messages um2
      JOIN keywords k2 ON um2.keyword_id = k2.id
      WHERE um2.user_id = u.id
      ORDER BY um2.count DESC
      LIMIT 1
    )
  ) AS stats
FROM users u
JOIN user_messages um ON um.user_id = u.id
JOIN keywords k ON um.keyword_id = k.id
WHERE u.username = $1
GROUP BY u.id, u.username;
