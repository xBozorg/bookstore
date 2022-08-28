-- bookstore.promo definition

CREATE TABLE `promo` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(20) NOT NULL,
  `expiration` datetime NOT NULL DEFAULT '2200-01-02 15:04:05',
  `limit` int unsigned NOT NULL DEFAULT '1',
  `percentage` int unsigned NOT NULL,
  `max_price` int unsigned DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;