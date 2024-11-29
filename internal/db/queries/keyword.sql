-- name: UpsertKeyword :one
INSERT INTO keywords (streamer_id, keyword, active)
VALUES ($1, $2, TRUE)
ON CONFLICT (streamer_id, keyword) DO UPDATE SET
    active = EXCLUDED.active
RETURNING id;

-- name: DeleteKeyword :one
DELETE FROM keywords
WHERE id = $1
RETURNING id, keyword;

-- name: GetGlobalKeywordsCount :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um ON um.keyword_id = k.id
  AND um.message_date BETWEEN $1 AND $2
WHERE k.streamer_id = $3
GROUP BY
  k.id, k.keyword, k.active;

-- name: GetGlobalKeywordsCountDesc :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um ON um.keyword_id = k.id
  AND um.message_date BETWEEN $1 AND $2
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  total_count DESC;

-- name: GetKeywordsByStreamer :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um ON um.keyword_id = k.id
  AND um.streamer_id = $1
  AND um.message_date BETWEEN $2 AND $3
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  total_count DESC;

-- name: GetGlobalKeywordsCountAscPaginated :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um
  ON um.keyword_id = k.id AND um.streamer_id = $1
WHERE
  um.message_date BETWEEN $2 AND $3
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  total_count ASC
LIMIT $4 OFFSET $5;

-- name: GetGlobalKeywordsCountDescPaginated :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um
  ON um.keyword_id = k.id AND um.streamer_id = $1
WHERE
  um.message_date BETWEEN $2 AND $3
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  total_count DESC
LIMIT $4 OFFSET $5;

-- name: GetCountsPerUserPerKeywordById :one
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
WHERE u.id = $1
GROUP BY u.id, u.username;


-- name: GetGlobalKeywordsCountPaginated :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um
  ON um.keyword_id = k.id
WHERE
  um.message_date BETWEEN $1 AND $2
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  k.id ASC
LIMIT $3 OFFSET $4;


-- name: GetGlobalKeywordsCountAsc :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um
  ON um.keyword_id = k.id
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  total_count ASC;


-- name: GetKeywordById :one
SELECT 
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um
  ON um.keyword_id = k.id
WHERE
  k.id = $1
GROUP BY
  k.id, k.keyword, k.active;
