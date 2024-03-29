CREATE DATABASE IF NOT EXISTS `inventory-db`;

USE `inventory-db`;

CREATE TABLE IF NOT EXISTS `products`(
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(255) NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `product_options`(
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(100) NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `product_option_values`(
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`product_option_id` INT(11) NOT NULL,
	`value` VARCHAR(100) NOT NULL,
	FOREIGN KEY (`product_option_id`)
		REFERENCES `product_options`(`id`),
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `product_options` (`id`, `name`) VALUES
(1, "size"),
(2, "color");

INSERT INTO `product_option_values` (`id`, `product_option_id`, `value`) VALUES
(1, 1, "S"),
(2, 1, "M"),
(3, 1, "L"),
(4, 1, "XL"),
(5, 1, "XXL"),
(6, 1, "XXXL"),
(7, 2, "Broken White"),
(8, 2, "Navy"),
(9, 2, "Black"),
(10, 2, "Salem"),
(11, 2, "Yellow"),
(12, 2, "White"),
(13, 2, "Khaki"),
(14, 2, "Red");

CREATE TABLE IF NOT EXISTS `product_variants`(
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`product_id` INT(11) NOT NULL,
	`sku` VARCHAR(50) NOT NULL,
	`fg_delete` TINYINT(1) DEFAULT 0,
	FOREIGN KEY (`product_id`)
		REFERENCES `products`(`id`),
	FOREIGN KEY (`product_option_value_id`)
		REFERENCES `product_option_values`(`id`),
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `product_variant_options`(
	`product_variant_id` INT(11) NOT NULL,
	`product_option_value_id` INT(11) NOT NULL,
	FOREIGN KEY (`product_variant_id`)
		REFERENCES `product_variants`(`id`),
	FOREIGN KEY (`product_option_value_id`)
		REFERENCES `product_option_values`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `purchase_orders`(
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`po_number` VARCHAR(50)
		COMMENT "null means receipt is lost",
	`time` DATETIME NOT NULL,
	`fg_delete` TINYINT(1) DEFAULT 0,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `purchase_order_details`(
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`purchase_order_id` INT(11) NOT NULL,
	`product_variant_id` INT(11) NOT NULL,
	`qty` INT(11) NOT NULL,
	`price` INT(11) NOT NULL,
	FOREIGN KEY `purchase_order_id`
		REFERENCES `purchase_orders`(`id`),
	FOREIGN KEY `product_variant_id`
		REFERENCES `product_variants`(`id`),
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `stock_in`(
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`purchase_order_detail_id` INT(11) NOT NULL,
	`time` DATETIME NOT NULL,
	`receive_qty` INT(11) NOT NULL,
	`fg_delete` TINYINT(1) DEFAULT 0,
	FOREIGN KEY (`purchase_order_detail_id`)
		REFERENCES `purchase_order_details`(`id`),
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `sales_orders`(
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`so_number` VARCHAR(50) NOT NULL,
	`time` DATETIME NOT NULL,
	`fg_delete` TINYINT(1) DEFAULT 0,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `sales_order_details`(
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`sales_order_id` INT(11) NOT NULL,
	`product_variant_id` INT(11) NOT NULL,
	`qty` INT(11) NOT NULL,
	`price` INT(11) NOT NULL,
	FOREIGN KEY (`sales_order_id`)
		REFERENCES `sales_orders`(`id`),
	FOREIGN KEY (`product_variant_id`)
		REFERENCES `product_variants`(`id`),
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `stock_out`(
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`sales_order_detail_id` INT(11)
		COMMENT "null means not for sale such as lost or damaged or sampling",
	`product_variant_id` INT(11) NOT NULL,
	`time` DATETIME NOT NULL,
	`qty` INT(11) NOT NULL,
	`notes` TEXT,
	`fg_delete` TINYINT(1) DEFAULT 0,
	FOREIGN KEY (`product_variant_id`)
		REFERENCES `product_variants`(`id`),
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;