CREATE TABLE tasks (
  id INT NOT NULL AUTO_INCREMENT,
  user_id INT NOT NULL,
  title varchar(255) DEFAULT NULL,
  note text DEFAULT NULL,
  completed integer DEFAULT 0,
  created_at TIMESTAMP DEFAULT NULL,
  updated_at TIMESTAMP DEFAULT NULL,
  PRIMARY KEY(id),
  CONSTRAINT tasks_fk_user_id
    FOREIGN KEY (`user_id`)
    REFERENCES users (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
