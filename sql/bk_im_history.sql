CREATE TABLE `bk_im_history` (
 `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
 `f_publish_id` int(11) NOT NULL COMMENT '发消息者',
 `t_publish_id` int(11) NOT NULL COMMENT '收消息者',
 `content` varchar(1024) NOT NULL,
 `c_time` int(11) NOT NULL COMMENT '发布时间',
 `is_readed` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0未读，1已读',
 `is_received` tinyint(2) NOT NULL DEFAULT '0' COMMENT '1已收到，0未收到',
 PRIMARY KEY (`id`),
 KEY `from_to_time` (`f_publish_id`,`t_publish_id`,`c_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天历史'
