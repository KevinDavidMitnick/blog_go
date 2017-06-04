SET NAMES utf8;

DROP DATABASE IF EXISTS blog;
CREATE DATABASE blog;

USE blog;

CREATE TABLE `article` (
  `id`      INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title`   VARCHAR(255)     NOT NULL,
  `url`     VARCHAR(255)     NOT NULL,
  `content` LONGTEXT         NOT NULL,
  `tags`    VARCHAR(255),
  `status`  INT(1),
  `created` INT(10),
  PRIMARY KEY (`id`),
  UNIQUE KEY (`url`)
)
  ENGINE =INNODB
  DEFAULT CHARSET =utf8;

CREATE TABLE `page` (
  `name`    VARCHAR(32) NOT NULL PRIMARY KEY,
  `content` LONGTEXT    NOT NULL
)
  ENGINE =INNODB
  DEFAULT CHARSET =utf8;
