-- 用户表
CREATE TABLE `sys_user`
(
    `user_id`     int(20) NOT NULL AUTO_INCREMENT,
    `username`    varchar(50) NOT NULL COMMENT '用户名',
    `password`    varchar(50) NOT NULL COMMENT '密码',
    `salt`        varchar(16) NOT NULL COMMENT '盐',
    `email`       varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
    `create_time` varchar(50) NOT NULL COMMENT '创建时间',
    `update_time` varchar(50)          DEFAULT NULL COMMENT '更新时间',
    `is_admin`    int(11) NOT NULL DEFAULT '0' COMMENT '是否是 admin 用户',
    `status`      int(11) NOT NULL DEFAULT '1' COMMENT '启用状态',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `UQE_user_name` (`username`),
    UNIQUE KEY `UQE_user_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户表';

INSERT INTO sys_user(username, password, salt, email, create_time, update_time, is_admin, status)
VALUES ('admin', 'f3e251e3242469a361cdf2653a75f70e', 'JWK1tU', '123@123.com', '2022-01-28 12:22:18',
        '2020-06-20 14:49:36', 1, 1);

CREATE TABLE `sys_user_token`
(
    `user_id`     int(20) NOT NULL COMMENT '用户 id',
    `token`       varchar(200) NOT NULL COMMENT 'token',
    `expire_time` varchar(50) DEFAULT '' COMMENT '过期时间',
    `update_time` varchar(50) DEFAULT '' COMMENT '更新时间',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `UQE_token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户Token';

CREATE TABLE `sys_login_log`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT,
    `username`    varchar(100) NOT NULL COMMENT '用户名',
    `ip`          varchar(64)  NOT NULL COMMENT 'IP地址',
    `create_time` varchar(50)  NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登陆日志管理';

-- 网站配置表
CREATE TABLE `sys_website`
(
    `website_id`          int(20) NOT NULL AUTO_INCREMENT,
    `file_name`           varchar(50) NOT NULL DEFAULT '' COMMENT '配置文件名称',
    `server_name`         varchar(50) NOT NULL DEFAULT '' COMMENT '网站域名',
    `root_directory`      varchar(50) NOT NULL DEFAULT '' COMMENT '网站根目录',
    `home_page`           varchar(50) NOT NULL DEFAULT '' COMMENT '网站主页',
    `http_port`           int(11) NOT NULL DEFAULT 80 COMMENT 'http 监听端口',
    `support_ssl`         int(11) NOT NULL DEFAULT '0' COMMENT '启用 TLS 状态 1: 启用 0: 禁用',
    `https_port`          int(11) NOT NULL DEFAULT 443 COMMENT 'https 监听端口',
    `ssl_certificate`     varchar(50) NOT NULL DEFAULT '' COMMENT 'TLS 证书路径',
    `ssl_certificate_key` varchar(50) NOT NULL DEFAULT '' COMMENT '私钥路径',
    `status`              int(11) NOT NULL DEFAULT '0' COMMENT '启用状态 1: 正常 0:禁用',
    `create_time`         varchar(50) NOT NULL COMMENT '创建时间',
    `update_time`         varchar(50)          DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`website_id`),
    UNIQUE KEY `UQE_file_name` (`file_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '网站配置表';
