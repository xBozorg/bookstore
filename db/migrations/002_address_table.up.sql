-- bookstore.address definition

CREATE TABLE `address` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `country` varchar(50) CHARACTER SET utf8mb3 NOT NULL,
  `province` varchar(50) CHARACTER SET utf8mb3 NOT NULL,
  `city` varchar(50) CHARACTER SET utf8mb3 NOT NULL,
  `street` varchar(50) CHARACTER SET utf8mb3 NOT NULL,
  `postalcode` varchar(20) CHARACTER SET utf8mb3 NOT NULL,
  `no` varchar(5) CHARACTER SET utf8mb3 NOT NULL,
  `description` varchar(50) CHARACTER SET utf8mb3 DEFAULT NULL,
  `userID` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `address_FK` (`userID`),
  CONSTRAINT `address_FK` FOREIGN KEY (`userID`) REFERENCES `user` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;