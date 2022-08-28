-- bookstore.phone definition

CREATE TABLE `phone` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(2) NOT NULL,
  `phonenumber` varchar(20) NOT NULL,
  `userID` varchar(60) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `phone_phonenumber_uindex` (`phonenumber`),
  KEY `userID` (`userID`),
  CONSTRAINT `phone_ibfk_1` FOREIGN KEY (`userID`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb3;