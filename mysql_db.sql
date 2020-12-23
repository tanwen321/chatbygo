-- MySQL dump 10.17  Distrib 10.3.22-MariaDB, for debian-linux-gnu (aarch64)
--
-- Host: localhost    Database: chat
-- ------------------------------------------------------
-- Server version       10.3.22-MariaDB-0+deb10u1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `chat_groups`
--

DROP TABLE IF EXISTS `chat_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `chat_groups` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `nickname` varchar(50) DEFAULT '' COMMENT '群组名',
  `usinge` varchar(128) DEFAULT '' COMMENT '群组说明',
  `status` tinyint(3) unsigned DEFAULT 1 COMMENT '状态 0为禁用、1为启用',
  `gimage` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20004 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chat_groups`
--

LOCK TABLES `chat_groups` WRITE;
/*!40000 ALTER TABLE `chat_groups` DISABLE KEYS */;
INSERT INTO `chat_groups` VALUES (20000,'管理员小组','别总是自卑，你永远比自己想象得更好',1,'/touxiang/tx_1365.jpg'),(20001,'车友吧','人成熟的过程，其实是学会与自我相处的过程',1,'/touxiang/tx_1705.jpg'),(20002,'钓鱼的人','愿你历尽千帆，归来贼特么有钱',1,'/touxiang/tx_1415.jpg'),(20003,'打 工 的 人','打工人,打工魂,网上办公要留神!',1,'/touxiang/tx_1753.jpg');
/*!40000 ALTER TABLE `chat_groups` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `chat_user_friends`
--

DROP TABLE IF EXISTS `chat_user_friends`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `chat_user_friends` (
  `userid` int(10) unsigned NOT NULL,
  `frendid` int(10) unsigned NOT NULL,
  `status` tinyint(3) unsigned DEFAULT 1 COMMENT '状态 0为禁用1为启用',
  PRIMARY KEY (`userid`,`frendid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户好友';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chat_user_friends`
--

LOCK TABLES `chat_user_friends` WRITE;
/*!40000 ALTER TABLE `chat_user_friends` DISABLE KEYS */;
INSERT INTO `chat_user_friends` VALUES (10000,10001,1),(10001,10000,1),(10001,10002,1),(10001,10003,1),(10001,10004,1),(10002,10001,1),(10002,10003,1),(10003,10001,1),(10003,10002,1),(10004,10001,1);
/*!40000 ALTER TABLE `chat_user_friends` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `chat_user_groups`
--

DROP TABLE IF EXISTS `chat_user_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `chat_user_groups` (
  `userid` int(10) unsigned NOT NULL,
  `groupid` int(10) unsigned NOT NULL,
  `status` tinyint(3) unsigned DEFAULT 1 COMMENT '状态 0为禁用1为启用',
  PRIMARY KEY (`userid`,`groupid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户群组';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chat_user_groups`
--

LOCK TABLES `chat_user_groups` WRITE;
/*!40000 ALTER TABLE `chat_user_groups` DISABLE KEYS */;
INSERT INTO `chat_user_groups` VALUES (10000,20000,1),(10001,20000,1),(10001,20001,1),(10001,20002,1),(10002,20001,1),(10002,20002,1),(10003,20001,1),(10003,20002,1),(10004,20003,1);
/*!40000 ALTER TABLE `chat_user_groups` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `chat_user_tokens`
--

DROP TABLE IF EXISTS `chat_user_tokens`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `chat_user_tokens` (
  `userid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `utoken` varchar(300) DEFAULT NULL,
  `status` tinyint(3) unsigned DEFAULT 1 COMMENT '状态 0为禁用1为启用',
  PRIMARY KEY (`userid`)
) ENGINE=InnoDB AUTO_INCREMENT=10005 DEFAULT CHARSET=utf8 COMMENT='token管理';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chat_user_tokens`
--

--
-- Table structure for table `chat_users`
--

DROP TABLE IF EXISTS `chat_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `chat_users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) DEFAULT '' COMMENT '账号',
  `password` varchar(50) DEFAULT '' COMMENT '密码',
  `urool` int(10) DEFAULT 2000 COMMENT '角色',
  `nickname` varchar(50) DEFAULT '' COMMENT '昵称',
  `usinge` varchar(128) DEFAULT '' COMMENT '签名',
  `status` tinyint(3) unsigned DEFAULT 1 COMMENT '状态 0为禁用、1为启用',
  `uimage` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10005 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chat_users`
--

LOCK TABLES `chat_users` WRITE;
/*!40000 ALTER TABLE `chat_users` DISABLE KEYS */;
INSERT INTO `chat_users` VALUES (10000,'admin','pawd@123',0,'管理员','我们在一起',1,'/touxiang/tx005.jpg'),(10001,'zikuang','pawd@123',10,'紫色的矿','如果你够独立,谁都可以失去',1,'/touxiang/tx003.jpg'),(10002,'maomao','pawd@123',10,'毛毛虫','要偷偷的努力，希望自己也能成为别人的梦想',1,'/touxiang/tx006.jpg'),(10003,'zhangsan','pawd@123',10,'张三','以后的你，会为自己现在所做的努力；而感到庆幸，不要在最好的年纪选择安逸',1,'/touxiang/tx008.jpg'),(10004,'iphone14','test123456',10,'我是最新的','这家伙很懒，什么都没留下',1,'/touxiang/tx_1484.jpg');
/*!40000 ALTER TABLE `chat_users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-12-23 10:59:59