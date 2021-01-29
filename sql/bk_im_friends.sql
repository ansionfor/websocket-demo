CREATE TABLE `bk_im_friends` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `publisher_id` int(11) NOT NULL,
 `friends_ids` varchar(1024) NOT NULL COMMENT '好友id列表，[id,id]',
 PRIMARY KEY (`id`),
 UNIQUE KEY `publisher_id` (`publisher_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8