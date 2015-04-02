-- up
CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `display_name` varchar(255) NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_email_index` (`email`)
);
ALTER TABLE `user` ADD UNIQUE INDEX `email_unique` (`email` ASC);

-- down
DROP TABLE `user`;