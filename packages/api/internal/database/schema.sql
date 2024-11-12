CREATE TABLE servers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  max_players INTEGER,
  start_map VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ports (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  port INTEGER UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  server_id INTEGER,
  FOREIGN KEY(server_id) REFERENCES servers(id)
);
