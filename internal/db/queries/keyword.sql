-- name: UpsertKeyword :one
INSERT INTO keywords (keyword)
VALUES ($1)
ON CONFLICT (keyword) DO UPDATE SET keyword = EXCLUDED.keyword
RETURNING id;

-- name: GetGlobalKeywordsCountDesc :many
SELECT
  k.id as keyword_id,
  k.keyword,
  SUM(um.count) AS total_count
FROM user_messages um
JOIN keywords k ON um.keyword_id = k.id
GROUP BY k.id, k.keyword
ORDER BY total_count DESC;

-- name: GetGlobalKeywordsCountAsc :many
SELECT
  k.id as keyword_id,
  k.keyword,
  SUM(um.count) AS total_count
FROM user_messages um
JOIN keywords k ON um.keyword_id = k.id
GROUP BY k.id, k.keyword
ORDER BY total_count ASC;


-- name: GetGlobalKeywordsCount :many
SELECT
  k.id as keyword_id,
  k.keyword,
  SUM(um.count) AS total_count
FROM user_messages um
JOIN keywords k ON um.keyword_id = k.id
GROUP BY k.id, k.keyword
ORDER BY k.id ASC;

-- name: GetKeywordById :one
SELECT 
  k.id as keyword_id,
  k.keyword,
  SUM(um.count) AS total_count
FROM user_messages um
JOIN keywords k ON um.keyword_id = k.id
WHERE k.id = $1
GROUP BY k.id, k.keyword;


-- name: GetActiveKeywords :many
SELECT
    id,
    keyword,
    active
FROM keywords
WHERE active = TRUE;
