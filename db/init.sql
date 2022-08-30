CREATE TABLE IF NOT EXISTS `admin` (
  `id` varchar(60) NOT NULL,
  `password` binary(60) NOT NULL,
  `phonenumber` varchar(20) NOT NULL,
  `email` varchar(150) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT IGNORE INTO admin (id,password,phonenumber,email) VALUES
	 ('2a25bcf0-8c72-4404-9ae7-2f94de398bae',0x243261243134243365486E3871343735735A473632766C64694255342E2E62396347504D50786A3875757033624D6276316866756148474152485A4B,'09112223333','admin@admin.com');

CREATE TABLE IF NOT EXISTS `user` (
  `id` varchar(60) NOT NULL,
  `email` varchar(150) NOT NULL,
  `password` binary(60) NOT NULL,
  `username` varchar(40) DEFAULT NULL,
  `firstname` varchar(80) NOT NULL,
  `lastname` varchar(80) NOT NULL,
  `regdate` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `address` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `country` varchar(50) NOT NULL,
  `province` varchar(50) NOT NULL,
  `city` varchar(50) NOT NULL,
  `street` varchar(50) NOT NULL,
  `postalcode` varchar(20) NOT NULL,
  `no` varchar(5) NOT NULL,
  `description` varchar(50) DEFAULT NULL,
  `userID` varchar(60) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `address_FK` (`userID`),
  CONSTRAINT `address_FK` FOREIGN KEY (`userID`) REFERENCES `user` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `phone` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(2) NOT NULL,
  `phonenumber` varchar(20) NOT NULL,
  `userID` varchar(60) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `phone_phonenumber_uindex` (`phonenumber`),
  KEY `userID` (`userID`),
  CONSTRAINT `phone_ibfk_1` FOREIGN KEY (`userID`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `language` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `language_UN` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `publisher` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `publisher_UN` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `author` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `author_UN` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `topic` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `topic_UN` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `book` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL,
  `isbn` varchar(13) NOT NULL,
  `pages` int unsigned NOT NULL,
  `description` varchar(500) DEFAULT NULL,
  `year` year NOT NULL,
  `date` datetime NOT NULL,
  `digital_price` int unsigned NOT NULL,
  `digital_discount` int unsigned DEFAULT '0',
  `physical_price` int unsigned NOT NULL,
  `physical_discount` int unsigned DEFAULT '0',
  `physical_stock` int unsigned NOT NULL,
  `pdf` varchar(150) DEFAULT NULL,
  `epub` varchar(150) DEFAULT NULL,
  `djvu` varchar(150) DEFAULT NULL,
  `azw` varchar(150) DEFAULT NULL,
  `txt` varchar(150) DEFAULT NULL,
  `docx` varchar(150) DEFAULT NULL,
  `lang_id` int unsigned NOT NULL,
  `cover_front` varchar(150) NOT NULL,
  `cover_back` varchar(150) NOT NULL,
  `publisher` int unsigned NOT NULL DEFAULT '0',
  `availability` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `book_UN` (`isbn`),
  KEY `book_FK` (`lang_id`),
  KEY `book_FK_1` (`publisher`),
  CONSTRAINT `book_FK` FOREIGN KEY (`lang_id`) REFERENCES `language` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `book_FK_1` FOREIGN KEY (`publisher`) REFERENCES `publisher` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `book_author` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `book_id` int unsigned NOT NULL,
  `author_id` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `book_author_FK_1` (`author_id`),
  KEY `book_author_FK` (`book_id`),
  CONSTRAINT `book_author_FK` FOREIGN KEY (`book_id`) REFERENCES `book` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `book_author_FK_1` FOREIGN KEY (`author_id`) REFERENCES `author` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `book_topic` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `book_id` int unsigned NOT NULL,
  `topic_id` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `book_topic_FK` (`book_id`),
  KEY `book_topic_FK_1` (`topic_id`),
  CONSTRAINT `book_topic_FK` FOREIGN KEY (`book_id`) REFERENCES `book` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `book_topic_FK_1` FOREIGN KEY (`topic_id`) REFERENCES `topic` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `promo` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(20) NOT NULL,
  `expiration` datetime NOT NULL DEFAULT '2200-01-02 15:04:05',
  `limit` int unsigned NOT NULL DEFAULT '1',
  `percentage` int unsigned NOT NULL,
  `max_price` int unsigned DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `promo_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `promo_id` int unsigned NOT NULL,
  `user_id` varchar(60)   NOT NULL,
  PRIMARY KEY (`id`),
  KEY `promo_user_FK_1` (`promo_id`),
  KEY `promo_user_FK` (`user_id`),
  CONSTRAINT `promo_user_FK` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `promo_user_FK_1` FOREIGN KEY (`promo_id`) REFERENCES `promo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `orders` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `creation_date` datetime NOT NULL,
  `receipt_date` datetime DEFAULT NULL,
  `status` int unsigned NOT NULL,
  `total` int unsigned NOT NULL,
  `stn` varchar(50) DEFAULT NULL,
  `user_id` varchar(60)   NOT NULL,
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
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `item` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `book_id` int unsigned NOT NULL,
  `type` int unsigned NOT NULL,
  `quantity` int unsigned NOT NULL DEFAULT '1',
  `order_id` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `item_UN` (`book_id`,`type`,`quantity`),
  KEY `item_FK1` (`order_id`),
  CONSTRAINT `item_FK1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `zarinpal` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `order_id` int unsigned NOT NULL,
  `authority` varchar(36) NOT NULL,
  `ref_id` int DEFAULT NULL,
  `code` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `zarinpal_FK` (`order_id`),
  CONSTRAINT `zarinpal_FK` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;