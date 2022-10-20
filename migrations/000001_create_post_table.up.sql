CREATE Table post(
  id serial PRIMARY KEY,
  name VARCHAR(200),
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  delete_at TIMESTAMP
);
CREATE Table media(
  post_id int REFERENCES post(id),
  name VARCHAR(100),
  link TEXT,
  type VARCHAR(100)
);