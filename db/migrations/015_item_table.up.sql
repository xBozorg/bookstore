-- bookstore.item definition

CREATE TABLE `item` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `book_id` int unsigned NOT NULL,
  `type` int unsigned NOT NULL,
  `quantity` int unsigned NOT NULL DEFAULT '1',
  `order_id` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `item_UN` (`book_id`,`type`,`quantity`),
  KEY `item_FK1` (`order_id`),
  CONSTRAINT `item_FK1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;