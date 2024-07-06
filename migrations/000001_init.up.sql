CREATE TABLE users (
  id INT PRIMARY KEY,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  link VARCHAR(255) UNIQUE
);

CREATE TABLE money (
  user_id INT,
  gold INT DEFAULT 0,
  silver INT DEFAULT 0,
  CONSTRAINT money_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT money_user_id_unique UNIQUE (user_id)
);

CREATE TABLE horse (
  user_id INT,
  level SMALLINT DEFAULT 1,
  distance INT DEFAULT 0,
  CONSTRAINT horse_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT horse_user_id_unique UNIQUE (user_id)
);

CREATE TABLE gunfight_queue (
  user_id INT NOT NULL,
  gold INT NOT NULL,
  CONSTRAINT gunfight_queue_user_id_unique UNIQUE (user_id)
);

CREATE TABLE gunfight (
  id SERIAL PRIMARY KEY,
  user_1_id INT NOT NULL,
  user_2_id INT NOT NULL,
  winner_id INT,
  start_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  end_date TIMESTAMP,
  CONSTRAINT fk_user_1 FOREIGN KEY (user_1_id) REFERENCES users(id),
  CONSTRAINT fk_user_2 FOREIGN KEY (user_2_id) REFERENCES users(id),
  CONSTRAINT chk_winner_id CHECK (winner_id IS NULL OR winner_id = user_1_id OR winner_id = user_2_id)
);

CREATE TABLE gunfight_health (
  gunfight_id INT,
  user_id INT,
  health SMALLINT DEFAULT 3,
  CONSTRAINT gunfight_health_gunfight_id_fk FOREIGN KEY (gunfight_id) REFERENCES gunfight(id),
  CONSTRAINT gunfight_health_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT gunfight_health_gunfight_id_unique UNIQUE (gunfight_id),
  CONSTRAINT gunfight_health_user_id_unique UNIQUE (user_id)
);
