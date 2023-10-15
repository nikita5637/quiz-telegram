-- +goose Up

DROP TABLE IF EXISTS `request`;
CREATE TABLE IF NOT EXISTS `request` (
  `id` int(11) NOT NULL,
  `group_uuid` varchar(36) NOT NULL,
  `uuid` varchar(36) NULL,
  `body` text NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

ALTER TABLE `request`
  ADD PRIMARY KEY (`id`), ADD UNIQUE KEY `uuid` (`uuid`);

ALTER TABLE `request`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=1;

-- +goose StatementBegin

CREATE TRIGGER `request_before_insert` 
BEFORE INSERT ON `request` FOR EACH ROW 
BEGIN
	IF new.uuid IS NULL THEN
		SET new.uuid = uuid();
	END IF;
END

-- +goose StatementEnd
-- +goose Down

DROP TRIGGER `request_before_insert`;

DROP TABLE IF EXISTS `request`;
