-- bookstore.`language` definition

CREATE TABLE `language` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `language_UN` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;