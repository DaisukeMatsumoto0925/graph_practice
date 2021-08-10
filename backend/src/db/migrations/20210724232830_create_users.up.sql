CREATE TABLE users (
  id INT NOT NULL AUTO_INCREMENT,
  email varchar(255) DEFAULT NULL,
  name varchar(255) DEFAULT NULL,
  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(id),
  UNIQUE unique_index_email_on_users (`email`),
  UNIQUE unique_index_name_on_users (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
