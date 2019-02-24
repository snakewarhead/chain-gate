SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for coin
-- ----------------------------
DROP TABLE IF EXISTS `coin`;
CREATE TABLE `coin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `enable` int(11) NOT NULL DEFAULT '0',
  `main_address` varchar(256) NOT NULL DEFAULT '',
  `public_key` varchar(256) NOT NULL DEFAULT '',
  `password` varchar(256) NOT NULL DEFAULT '',
  `api_url` varchar(256) NOT NULL DEFAULT '',
  `api_wallet_url` varchar(256) NOT NULL DEFAULT '',
  `confirm_num` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of coin
-- ----------------------------
INSERT INTO `coin` VALUES ('1', 'EOS', '1', 'test.account', 'EOS4xQtqisWfSXApejE9t2vnTfbBPUoE8wjF1SuvvsfV7LHXBgKNX', 'PW5Kc1cFRaREjX7m5scd4MmmHCRAbW6Q1DhqSKk39x4dM32zzr7qt', 'http://127.0.0.1:8888', 'http://127.0.0.1:7011', '5');

-- ----------------------------
-- Table structure for transaction_history
-- ----------------------------
DROP TABLE IF EXISTS `transaction_history`;
CREATE TABLE `transaction_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `coin_id` int(11) NOT NULL,
  `contract` varchar(128) NOT NULL DEFAULT '',
  `tx_id` varchar(128) NOT NULL,
  `is_main` tinyint(2) NOT NULL DEFAULT '0',
  `symbol` varchar(256) NOT NULL DEFAULT '',
  `direction` tinyint(2) NOT NULL,
  `status` tinyint(2) NOT NULL DEFAULT '0',
  `from_address` varchar(256) NOT NULL,
  `to_address` varchar(256) NOT NULL,
  `amount` varchar(64) NOT NULL,
  `fee` varchar(64) NOT NULL DEFAULT '',
  `memo` varchar(512) NOT NULL DEFAULT '',
  `create_time` int(11) NOT NULL,
  `update_time` int(11) NOT NULL,
  `block_num` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_cid_tid` (`coin_id`,`tx_id`),
  KEY `idx_from_address` (`from_address`(255)),
  KEY `idx_to_address` (`to_address`(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of transaction_history
-- ----------------------------
