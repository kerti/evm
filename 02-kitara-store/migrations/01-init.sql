CREATE TABLE IF NOT EXISTS `products` (
    `entity_id` CHAR(36) NOT NULL,
    `sku` VARCHAR(20) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `price` DECIMAL(10,2) NOT NULL,
    PRIMARY KEY (`entity_id`)
);

INSERT INTO `products` (`entity_id`, `sku`, `name`, `price`)
VALUES ('4164d34e-d055-4142-8d66-ff78e19045c5', 'EVM-SAMPLE-PRODUCT-1', 'Sample Product', 1000);

CREATE TABLE IF NOT EXISTS `inventory` (
    `entity_id` CHAR(36) NOT NULL,
    `product_entity_id` CHAR(36) NOT NULL,
    `qty_in_store` INT NOT NULL,
    `qty_reserved` INT NOT NULL,
    `qty_available` INT NOT NULL,
    PRIMARY KEY (`entity_id`)
);

INSERT INTO `inventory` (`entity_id`, `product_entity_id`, `qty_in_store`, `qty_reserved`, `qty_available`)
VALUES ('9e256347-aec6-4e30-9d90-a00e7181aba4', '4164d34e-d055-4142-8d66-ff78e19045c5', 10, 0, 10);

CREATE TABLE IF NOT EXISTS `orders` (
    `entity_id` CHAR(36) NOT NULL,
    `order_code` VARCHAR(255) NOT NULL,
    `total_price` DECIMAL(10,2) NOT NULL,
    `status` ENUM('new', 'processing', 'completed'),
    PRIMARY KEY (`entity_id`)
);

INSERT INTO `orders` (`entity_id`, `order_code`, `total_price`, `status`)
VALUES
('63b2ae3c-040a-4d64-89b2-b5e7e1951e89', 'EVM-TEST-ORDER-1', 10000, 'new'),
('5b27773a-efce-4e21-8474-2694ebdaa084', 'EVM-TEST-ORDER-1', 10000, 'new');

CREATE TABLE IF NOT EXISTS `order_items` (
    `entity_id` CHAR(36) NOT NULL,
    `order_entity_id` CHAR(36) NOT NULL,
    `product_entity_id` CHAR(36) NOT NULL,
    `qty` CHAR(36) NOT NULL,
    `price` DECIMAL(10,2) NOT NULL
);

INSERT INTO `order_items` (`entity_id`, `order_entity_id`, `product_entity_id`, `qty`, `price`)
VALUES
('3b9d6b54-3c96-4610-8b25-8805e86d85fc', '63b2ae3c-040a-4d64-89b2-b5e7e1951e89', '4164d34e-d055-4142-8d66-ff78e19045c5', 10, 10000),
('9854d163-d8ce-4247-86c6-3c4f9ac287bc', '5b27773a-efce-4e21-8474-2694ebdaa084', '4164d34e-d055-4142-8d66-ff78e19045c5', 10, 10000);