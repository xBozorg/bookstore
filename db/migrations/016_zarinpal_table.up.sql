-- bookstore.zarinpal definition

CREATE TABLE `zarinpal` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `order_id` int unsigned NOT NULL,
  `authority` varchar(36) NOT NULL,
  `ref_id` int DEFAULT NULL,
  `code` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `zarinpal_FK` (`order_id`),
  CONSTRAINT `zarinpal_FK` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;