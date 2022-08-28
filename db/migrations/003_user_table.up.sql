-- bookstore.`user` definition

CREATE TABLE `user` (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;