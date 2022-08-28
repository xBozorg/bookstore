-- bookstore.orders definition

CREATE TABLE `orders` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `creation_date` datetime NOT NULL,
  `receipt_date` datetime DEFAULT NULL,
  `status` int unsigned NOT NULL,
  `total` int unsigned NOT NULL,
  `stn` varchar(50) DEFAULT NULL,
  `user_id` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `promo_id` int unsigned DEFAULT NULL,
  `phone_id` int unsigned DEFAULT NULL,
  `address_id` int unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `order_FK` (`user_id`),
  KEY `orders_FK` (`promo_id`),
  KEY `orders_FK_1` (`address_id`),
  KEY `orders_FK_2` (`phone_id`),
  CONSTRAINT `order_FK` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `orders_FK` FOREIGN KEY (`promo_id`) REFERENCES `promo` (`id`),
  CONSTRAINT `orders_FK_1` FOREIGN KEY (`address_id`) REFERENCES `address` (`id`),
  CONSTRAINT `orders_FK_2` FOREIGN KEY (`phone_id`) REFERENCES `phone` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;