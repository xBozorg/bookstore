-- bookstore.promo_user definition

CREATE TABLE `promo_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `promo_id` int unsigned NOT NULL,
  `user_id` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `promo_user_FK_1` (`promo_id`),
  KEY `promo_user_FK` (`user_id`),
  CONSTRAINT `promo_user_FK` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `promo_user_FK_1` FOREIGN KEY (`promo_id`) REFERENCES `promo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;