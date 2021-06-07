CREATE TABLE tasks (
  id INT NOT NULL AUTO_INCREMENT,
  title varchar(255) DEFAULT NULL,
  note text DEFAULT NULL,
  completed integer DEFAULT 0,
  created_at TIMESTAMP DEFAULT NULL,
  updated_at TIMESTAMP DEFAULT NULL,
  PRIMARY KEY(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
