

#DROP USER 'user'@'172.17.0.1' ;

CREATE USER 'user'@'172.17.0.1' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON * . * TO 'user'@'172.17.0.1';
FLUSH PRIVILEGES;
SHOW GRANTS FOR 'user'@'172.17.0.1' ;

CREATE USER 'user'@'127.0.0.1' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON * . * TO 'user'@'127.0.0.1';
FLUSH PRIVILEGES;
SHOW GRANTS FOR 'user'@'127.0.0.1' ;



/*
SQLyog Community v12.4.0 (64 bit)
MySQL - 5.7.17-log : Database - goschool
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`goschool` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `goschool`;

/*Table structure for table `courses` */

DROP TABLE IF EXISTS `courses`;

CREATE TABLE `courses` (
  `ID` varchar(12) NOT NULL,
  `Title` varchar(30) DEFAULT NULL,
  `Description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Data for the table `courses` */

insert  into `courses`(`ID`,`Title`,`Description`) values 
('GOACTION01','GO IN ACTION I','Explore the practical aspect of Go software development. This course covers creating servers and clients in Go with networking.'),
('GOACTION02','GO IN ACTION II','Dive deeper and examine some of the practices and patterns that are generally adopted such as idiomatic Go, security and documentation.'),
('GOADVANCE01','GO ADVANCE','Learn advance concepts in Go programming such as packages, creation of data structures and error handling mechanism.'),
('GOBASIC01','GO BASICS','Gain fundamental knowledge and Go skills with an introduction to Golang programming.'),
('GOMS01','GO MICROSERVICES I','Learn the fundamental of microservice architecture, how to encode/decode in JSON, and RESTful communication.'),
('GOMS02','GO MICROSERVICES II','Accelerate the development of Go Projects which focuses on testing, monitoring and Continuous Integration/ Continuous Delivery.'),
('YODA123','GO GURU I','Learn guru\'s level technique of crafting Go\'s program.');

/*Table structure for table `trainers` */

DROP TABLE IF EXISTS `trainers`;

CREATE TABLE `trainers` (
  `ID` varchar(5) NOT NULL,
  `FirstName` varchar(30) DEFAULT NULL,
  `LastName` varchar(30) DEFAULT NULL,
  `Age` int(11) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Data for the table `trainers` */

insert  into `trainers`(`ID`,`FirstName`,`LastName`,`Age`) values 
('001','Adam','Smith',25),('002','Bertrand','Russel',33),
('003','Charlie','Munger',43);
('004','Dwight','Eisenhower',43);

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
