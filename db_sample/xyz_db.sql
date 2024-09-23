-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               10.4.27-MariaDB - mariadb.org binary distribution
-- Server OS:                    Win64
-- HeidiSQL Version:             12.2.0.6576
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- Dumping database structure for xyz_multifinance
DROP DATABASE IF EXISTS `xyz_multifinance`;
CREATE DATABASE IF NOT EXISTS `xyz_multifinance` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */;
USE `xyz_multifinance`;

-- Dumping structure for table xyz_multifinance.customers
DROP TABLE IF EXISTS `customers`;
CREATE TABLE IF NOT EXISTS `customers` (
  `id` varchar(144) NOT NULL,
  `nik` varchar(191) DEFAULT NULL,
  `email` varchar(191) DEFAULT NULL,
  `full_name` varchar(191) DEFAULT NULL,
  `salt` varchar(191) DEFAULT NULL,
  `password` longtext DEFAULT NULL,
  `legal_name` varchar(191) DEFAULT NULL,
  `birth_place` varchar(191) DEFAULT NULL,
  `birth_date` datetime(3) DEFAULT NULL,
  `salary` int(11) DEFAULT NULL,
  `photo_id` varchar(191) DEFAULT NULL,
  `selfie` varchar(191) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT current_timestamp(3),
  `updated_at` datetime(3) NOT NULL DEFAULT current_timestamp(3) ON UPDATE current_timestamp(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `nik` (`nik`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

-- Dumping structure for table xyz_multifinance.payments
DROP TABLE IF EXISTS `payments`;
CREATE TABLE IF NOT EXISTS `payments` (
  `id` varchar(50) DEFAULT NULL,
  `id_user` varchar(50) DEFAULT NULL,
  `id_transaksi` varchar(50) DEFAULT NULL,
  `nominal_normal` int(11) DEFAULT NULL,
  `is_paid` tinyint(4) DEFAULT NULL,
  `denda` int(11) DEFAULT NULL,
  `sisa_cicilan` int(11) DEFAULT NULL,
  `nominal_bayar` int(11) DEFAULT NULL,
  `interest` int(11) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `jatuh_tempo` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `paid_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `idx_transaction` int(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

-- Dumping structure for table xyz_multifinance.transaksis
DROP TABLE IF EXISTS `transaksis`;
CREATE TABLE IF NOT EXISTS `transaksis` (
  `id` varchar(36) NOT NULL,
  `nomer_kontrak` varchar(50) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `harga_otr` int(11) NOT NULL,
  `type` varchar(32) NOT NULL DEFAULT '',
  `admin_fee` int(11) NOT NULL,
  `tenor` int(11) NOT NULL,
  `jumlah_cicilan` int(11) NOT NULL,
  `cicilan_perbulan` int(11) NOT NULL,
  `nama_asset` varchar(50) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `nomer_kontrak` (`nomer_kontrak`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `transaksis_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `customers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
