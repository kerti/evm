CREATE TABLE IF NOT EXISTS `players` (
    `entity_id` CHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `ready_to_play` TINYINT(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (`entity_id`)
);

INSERT INTO `players` (`entity_id`, `name`, `ready_to_play`)
VALUES ('4164d34e-d055-4142-8d66-ff78e19045c5', 'Rahman', FALSE);

CREATE TABLE IF NOT EXISTS `containers` (
    `entity_id` CHAR(36) NOT NULL,
    `player_entity_id` CHAR(36) NOT NULL,
    `capacity` INT NOT NULL,
    `ball_count` INT NOT NULL,
    PRIMARY KEY (`entity_id`)
);

INSERT INTO `containers` (`entity_id`, `player_entity_id`, `capacity`, `ball_count`)
VALUES
('9e256347-aec6-4e30-9d90-a00e7181aba4', '4164d34e-d055-4142-8d66-ff78e19045c5', 10, 0),
('40e1fa5b-01a8-4255-8334-8a961ce74c57', '4164d34e-d055-4142-8d66-ff78e19045c5', 10, 0),
('3136c84f-9892-4b49-a474-357d58aeac90', '4164d34e-d055-4142-8d66-ff78e19045c5', 10, 0);
