-- name: UpsertKeyword :one
INSERT INTO keywords (keyword)
VALUES ($1)
ON CONFLICT (keyword) DO UPDATE SET keyword = EXCLUDED.keyword
RETURNING id;

-- name: GetGlobalKeywordsCountDesc :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count -- Using COALESCE to handle NULL
FROM
  keywords k
LEFT JOIN
  user_messages um ON um.keyword_id = k.id
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  total_count DESC;

-- name: GetGlobalKeywordsCountAsc :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um ON um.keyword_id = k.id
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  total_count ASC;

-- name: GetUsedKeywords :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um ON um.keyword_id = k.id
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  total_count ASC;

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
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  k.id ASC;

-- name: GetGlobalKeywordsCountPaginated :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um ON um.keyword_id = k.id
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  k.id ASC
LIMIT $1 OFFSET $2;

-- name: GetKeywordById :one
SELECT 
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um ON um.keyword_id = k.id
WHERE
  k.id = $1
GROUP BY
  k.id, k.keyword, k.active;

-- name: GetGlobalKeywordsCountDescPaginated :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um ON um.keyword_id = k.id
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  total_count DESC
LIMIT $1 OFFSET $2;

-- name: GetGlobalKeywordsCountAscPaginated :many
SELECT
  k.id AS keyword_id,
  k.keyword,
  k.active,
  COALESCE(SUM(um.count), 0) AS total_count
FROM
  keywords k
LEFT JOIN
  user_messages um ON um.keyword_id = k.id
GROUP BY
  k.id, k.keyword, k.active
ORDER BY
  total_count ASC
LIMIT $1 OFFSET $2;

-- name: GetInactiveKeywords :many
SELECT
    id,
    keyword,
    active
FROM
    keywords
WHERE
    active = FALSE;

-- name: IsKeywordActive :one
SELECT
  active
FROM
  keywords
WHERE
  id = $1;


-- name: GetGlobalKeywordsCountTotal :one
SELECT COUNT(*) FROM keywords;
