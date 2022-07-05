/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE='NO_AUTO_VALUE_ON_ZERO', SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


#
转储表 gateway_admin
# ------------------------------------------------------------

DROP TABLE IF EXISTS `gateway_admin`;

CREATE TABLE `gateway_admin`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_name`  varchar(20) NOT NULL,
    `password`   longtext    NOT NULL,
    `salt`       longtext DEFAULT NULL,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `is_delete`  bigint(20) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY          `idx_gateway_admin_username` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK
    TABLES `gateway_admin` WRITE;
/*!40000 ALTER TABLE `gateway_admin` DISABLE KEYS */;

INSERT INTO `gateway_admin` (`id`, `user_name`, `password`, `salt`, `created_at`, `updated_at`, `is_delete`)
VALUES (1, 'admin', '2823d896e9822c0833d41d4904f0c00756d718570fce49b9a379a62c804689d3', 'admin',
        '2022-05-11 12:33:45.000', '2022-05-11 12:33:45.000', 0);

/*!40000 ALTER TABLE `gateway_admin` ENABLE KEYS */;
UNLOCK
    TABLES;


#
转储表 gateway_service_access_control
# ------------------------------------------------------------

DROP TABLE IF EXISTS `gateway_service_access_control`;

CREATE TABLE `gateway_service_access_control`
(
    `id`                  bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `service_id`          bigint(20) unsigned DEFAULT NULL,
    `enable_auth`         bigint(20) DEFAULT NULL,
    `black_list`          longtext DEFAULT NULL,
    `white_list`          longtext DEFAULT NULL,
    `clientip_flow_limit` bigint(20) DEFAULT NULL,
    `service_flow_limit`  bigint(20) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY                   `idx_gateway_service_access_control_service_id` (`service_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK
    TABLES `gateway_service_access_control` WRITE;
/*!40000 ALTER TABLE `gateway_service_access_control` DISABLE KEYS */;

INSERT INTO `gateway_service_access_control` (`id`, `service_id`, `enable_auth`, `black_list`, `white_list`,
                                              `clientip_flow_limit`, `service_flow_limit`)
VALUES (1, 1, 0, '', '', 86, 0),
       (2, 2, 0, '', '', 0, 0),
       (3, 3, 0, '', '', 0, 0),
       (4, 4, 0, '', '', 0, 0),
       (5, 5, 0, '', '', 0, 0);

/*!40000 ALTER TABLE `gateway_service_access_control` ENABLE KEYS */;
UNLOCK
    TABLES;


#
转储表 gateway_service_grpc_rule
# ------------------------------------------------------------

DROP TABLE IF EXISTS `gateway_service_grpc_rule`;

CREATE TABLE `gateway_service_grpc_rule`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `service_id`      bigint(20) unsigned DEFAULT NULL,
    `port`            bigint(20) DEFAULT NULL,
    `header_transfer` longtext DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY               `idx_gateway_service_grpc_rule_service_id` (`service_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK
    TABLES `gateway_service_grpc_rule` WRITE;
/*!40000 ALTER TABLE `gateway_service_grpc_rule` DISABLE KEYS */;

INSERT INTO `gateway_service_grpc_rule` (`id`, `service_id`, `port`, `header_transfer`)
VALUES (1, 5, 8666, '');

/*!40000 ALTER TABLE `gateway_service_grpc_rule` ENABLE KEYS */;
UNLOCK
    TABLES;


#
转储表 gateway_service_http_rule
# ------------------------------------------------------------

DROP TABLE IF EXISTS `gateway_service_http_rule`;

CREATE TABLE `gateway_service_http_rule`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `service_id`      bigint(20) unsigned DEFAULT NULL,
    `rule_type`       bigint(20) DEFAULT NULL,
    `rule`            longtext DEFAULT NULL,
    `need_https`      bigint(20) DEFAULT NULL,
    `need_strip_uri`  bigint(20) DEFAULT NULL,
    `url_rewrite`     longtext DEFAULT NULL,
    `header_transfer` longtext DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY               `idx_gateway_service_http_rule_service_id` (`service_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK
    TABLES `gateway_service_http_rule` WRITE;
/*!40000 ALTER TABLE `gateway_service_http_rule` DISABLE KEYS */;

INSERT INTO `gateway_service_http_rule` (`id`, `service_id`, `rule_type`, `rule`, `need_https`, `need_strip_uri`,
                                         `url_rewrite`, `header_transfer`)
VALUES (1, 1, 0, '/test_http', 0, 0, '', ''),
       (2, 3, 0, '/test_http2', 0, 0, '', ''),
       (3, 4, 0, '/test_https', 1, 0, '', '');

/*!40000 ALTER TABLE `gateway_service_http_rule` ENABLE KEYS */;
UNLOCK
    TABLES;


#
转储表 gateway_service_info
# ------------------------------------------------------------

DROP TABLE IF EXISTS `gateway_service_info`;

CREATE TABLE `gateway_service_info`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `load_type`    bigint(20) DEFAULT NULL,
    `service_name` varchar(130) NOT NULL,
    `service_desc` longtext DEFAULT NULL,
    `updated_at`   datetime(3) DEFAULT NULL,
    `created_at`   datetime(3) DEFAULT NULL,
    `is_delete`    tinyint(4) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY            `idx_gateway_service_info_service_name` (`service_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK
    TABLES `gateway_service_info` WRITE;
/*!40000 ALTER TABLE `gateway_service_info` DISABLE KEYS */;

INSERT INTO `gateway_service_info` (`id`, `load_type`, `service_name`, `service_desc`, `updated_at`, `created_at`,
                                    `is_delete`)
VALUES (1, 0, 'test_http', '测试HTTP服务', '2022-07-05 12:05:48.337', '2022-07-05 12:05:48.337', 0),
       (2, 1, 'test_tcp', '测试MySQL 8555', '2022-07-05 12:11:50.581', '2022-07-05 12:11:50.581', 0),
       (3, 0, 'test_http2', '测试HTTP', '2022-07-05 12:22:46.071', '2022-07-05 12:22:46.071', 0),
       (4, 0, 'test_https', '测试HTTPS', '2022-07-05 12:23:46.270', '2022-07-05 12:23:46.270', 0),
       (5, 2, 'test_grpc', '测试GRPC', '2022-07-05 12:45:13.763', '2022-07-05 12:45:13.763', 0);

/*!40000 ALTER TABLE `gateway_service_info` ENABLE KEYS */;
UNLOCK
    TABLES;


#
转储表 gateway_service_load_balance
# ------------------------------------------------------------

DROP TABLE IF EXISTS `gateway_service_load_balance`;

CREATE TABLE `gateway_service_load_balance`
(
    `id`                       bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `service_id`               bigint(20) unsigned DEFAULT NULL,
    `round_type`               bigint(20) DEFAULT NULL,
    `ip_list`                  longtext DEFAULT NULL,
    `weight_list`              longtext DEFAULT NULL,
    `forbid_list`              longtext DEFAULT NULL,
    `upstream_connect_timeout` bigint(20) DEFAULT NULL,
    `upstream_header_timeout`  bigint(20) DEFAULT NULL,
    `upstream_idle_timeout`    bigint(20) DEFAULT NULL,
    `upstream_max_idle`        bigint(20) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY                        `idx_gateway_service_load_balance_service_id` (`service_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK
    TABLES `gateway_service_load_balance` WRITE;
/*!40000 ALTER TABLE `gateway_service_load_balance` DISABLE KEYS */;

INSERT INTO `gateway_service_load_balance` (`id`, `service_id`, `round_type`, `ip_list`, `weight_list`, `forbid_list`,
                                            `upstream_connect_timeout`, `upstream_header_timeout`,
                                            `upstream_idle_timeout`, `upstream_max_idle`)
VALUES (1, 1, 0, '127.0.0.1:2003', '50', NULL, 0, 0, 0, 0),
       (2, 2, 0, '127.0.0.1:3306', '10', NULL, 0, 0, 0, 0),
       (3, 3, 2, '127.0.0.1:2003,127.0.0.1:2004', '10,20', NULL, 0, 0, 0, 0),
       (4, 4, 2, '127.0.0.1:2003,127.0.0.1:2004', '2,3', NULL, 0, 0, 0, 0),
       (5, 5, 2, '127.0.0.1:2003,127.0.0.1:2004', '2,3', NULL, 0, 0, 0, 0);

/*!40000 ALTER TABLE `gateway_service_load_balance` ENABLE KEYS */;
UNLOCK
    TABLES;


#
转储表 gateway_service_tcp_rule
# ------------------------------------------------------------

DROP TABLE IF EXISTS `gateway_service_tcp_rule`;

CREATE TABLE `gateway_service_tcp_rule`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `service_id` bigint(20) unsigned DEFAULT NULL,
    `port`       bigint(20) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY          `idx_gateway_service_tcp_rule_service_id` (`service_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK
    TABLES `gateway_service_tcp_rule` WRITE;
/*!40000 ALTER TABLE `gateway_service_tcp_rule` DISABLE KEYS */;

INSERT INTO `gateway_service_tcp_rule` (`id`, `service_id`, `port`)
VALUES (1, 2, 8555);

/*!40000 ALTER TABLE `gateway_service_tcp_rule` ENABLE KEYS */;
UNLOCK
    TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
