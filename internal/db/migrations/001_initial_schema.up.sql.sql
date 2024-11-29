-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL -- Chatters' usernames
);


CREATE TABLE IF NOT EXISTS streamers (
  id SERIAL PRIMARY KEY,
  twitch_id TEXT UNIQUE NOT NULL, -- Twitch's unique user ID
  username TEXT NOT NULL,         -- Display name of the streamer
  access_token TEXT NOT NULL,     -- Twitch access token
  refresh_token TEXT NOT NULL,    -- Twitch refresh token
  expires_at TIMESTAMP NOT NULL,  -- Token expiry
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS keywords (
  id SERIAL PRIMARY KEY,
  streamer_id INTEGER NOT NULL REFERENCES streamers (id) ON DELETE CASCADE,
  keyword TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  UNIQUE (streamer_id, keyword) -- Ensures unique keywords per streamer
);

CREATE TABLE IF NOT EXISTS user_messages (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users (id),
  keyword_id INTEGER REFERENCES keywords (id) ON DELETE CASCADE,
  streamer_id INTEGER REFERENCES streamers (id) ON DELETE CASCADE,
  count INTEGER DEFAULT 0,
  last_message TEXT,
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
  message_date DATE NOT NULL DEFAULT CURRENT_DATE,
  UNIQUE (user_id, keyword_id, streamer_id, message_date)
);

-- +goose Down
DROP TABLE IF EXISTS user_messages;
DROP TABLE IF EXISTS keywords;
DROP TABLE IF EXISTS streamers;
DROP TABLE IF EXISTS users;
