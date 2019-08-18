CREATE TABLE IF NOT EXISTS `products`(
	`id` INTEGER PRIMARY KEY,
	`name` VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS `product_options`(
	`id` INTEGER PRIMARY KEY,
	`name` VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS `product_option_values`(
	`id` INTEGER PRIMARY KEY,
	`product_option_id` INT(11) NOT NULL,
	`value` VARCHAR(100) NOT NULL,
	FOREIGN KEY (`product_option_id`)
		REFERENCES `product_options`(`id`)
);

DELETE FROM `product_options`;
INSERT INTO `product_options` (`id`, `name`) VALUES
(1, "size"),
(2, "color");

DELETE FROM `product_option_values`;
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
	`id` INTEGER PRIMARY KEY,
	`product_id` INT(11) NOT NULL,
	`sku` VARCHAR(50) NOT NULL,
	`fg_delete` TINYINT(1) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS `product_variant_options`(
	`product_variant_id` INT(11) NOT NULL,
	`product_option_value_id` INT(11) NOT NULL,
	FOREIGN KEY (`product_variant_id`)
		REFERENCES `product_variants`(`id`),
	FOREIGN KEY (`product_option_value_id`)
		REFERENCES `product_option_values`(`id`)
);

CREATE TABLE IF NOT EXISTS `purchase_orders`(
	`id` INTEGER PRIMARY KEY,
	`po_number` VARCHAR(50),
	`time` DATETIME NOT NULL,
	`fg_delete` TINYINT(1) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS `purchase_order_details`(
	`id` INTEGER PRIMARY KEY,
	`purchase_order_id` INT(11) NOT NULL,
	`product_variant_id` INT(11) NOT NULL,
	`qty` INT(11) NOT NULL,
	`price` INT(11) NOT NULL
);

CREATE TABLE IF NOT EXISTS `stock_in`(
	`id` INTEGER PRIMARY KEY,
	`purchase_order_detail_id` INT(11) NOT NULL,
	`time` DATETIME NOT NULL,
	`receive_qty` INT(11) NOT NULL,
	`fg_delete` TINYINT(1) DEFAULT 0,
	FOREIGN KEY (`purchase_order_detail_id`)
		REFERENCES `purchase_order_details`(`id`)
);

CREATE TABLE IF NOT EXISTS `sales_orders`(
	`id` INTEGER PRIMARY KEY,
	`so_number` VARCHAR(50) NOT NULL,
	`time` DATETIME NOT NULL,
	`fg_delete` TINYINT(1) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS `sales_order_details`(
	`id` INTEGER PRIMARY KEY,
	`sales_order_id` INT(11) NOT NULL,
	`product_variant_id` INT(11) NOT NULL,
	`qty` INT(11) NOT NULL,
	`price` INT(11) NOT NULL,
	FOREIGN KEY (`sales_order_id`)
		REFERENCES `sales_orders`(`id`),
	FOREIGN KEY (`product_variant_id`)
		REFERENCES `product_variants`(`id`)
);

CREATE TABLE IF NOT EXISTS `stock_out`(
	`id` INTEGER PRIMARY KEY,
	`sales_order_detail_id` INT(11),
	`product_variant_id` INT(11),
	`time` DATETIME,
	`qty` INT(11) NOT NULL,
	`notes` TEXT,
	`fg_delete` TINYINT(1) DEFAULT 0,
	FOREIGN KEY (`product_variant_id`)
		REFERENCES `product_variants`(`id`)
);