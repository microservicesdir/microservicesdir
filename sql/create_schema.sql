# Creates projects table

CREATE TABLE IF NOT EXISTS `projects` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) DEFAULT NULL,
  `gitUrl` varchar(200) DEFAULT NULL,
  `owner` varchar(200) DEFAULT NULL,
  `language` varchar(60) DEFAULT NULL,
  `imported_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `manifest` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_project_name` (`name`),
  KEY `index_projects_on_names` (`name`),
  KEY `index_projects_on_owner` (`owner`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
