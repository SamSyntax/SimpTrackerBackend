-- +goose Up
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL
);

CREATE TABLE keywords (
  id SERIAL PRIMARY KEY,
  keyword TEXT UNIQUE NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE user_messages (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users (id),
  keyword_id INTEGER,
  count INTEGER DEFAULT 1,
  last_message TEXT,
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW (),
  UNIQUE (user_id, keyword_id),
  CONSTRAINT user_messages_keyword_id_fkey FOREIGN KEY (keyword_id) REFERENCES keywords (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE user_messages;

DROP TABLE users;

DROP TABLE keywords;
