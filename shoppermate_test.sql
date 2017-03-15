-- phpMyAdmin SQL Dump
-- version 4.6.4
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Mar 07, 2017 at 08:42 AM
-- Server version: 10.1.17-MariaDB
-- PHP Version: 7.1.1

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `shoppermate_test`
--

-- --------------------------------------------------------

--
-- Table structure for table `admin`
--

CREATE TABLE `admin` (
  `id` int(11) NOT NULL,
  `fullname` varchar(150) NOT NULL,
  `email` varchar(150) NOT NULL,
  `password` varchar(150) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `ads`
--

CREATE TABLE `ads` (
  `id` int(11) NOT NULL,
  `guid` varchar(100) NOT NULL,
  `img` varchar(255) NOT NULL,
  `front_name` varchar(150) NOT NULL,
  `name` varchar(255) NOT NULL,
  `body` text NOT NULL,
  `advertiser_id` int(11) NOT NULL,
  `campaign_id` int(11) NOT NULL,
  `item_id` int(11) NOT NULL,
  `category_id` varchar(255) NOT NULL,
  `subcategory_id` int(11) NOT NULL DEFAULT '0',
  `positive_tag` text NOT NULL,
  `negative_tag` text NOT NULL,
  `type` varchar(50) NOT NULL,
  `start_date` date DEFAULT NULL,
  `end_date` date DEFAULT NULL,
  `time` varchar(255) NOT NULL,
  `refresh_period` int(11) NOT NULL,
  `perlimit` int(11) NOT NULL,
  `cashback_amount` decimal(10,2) NOT NULL,
  `quota` int(11) NOT NULL DEFAULT '0',
  `status` varchar(60) NOT NULL DEFAULT 'draft',
  `grocer_exclusive` int(11) DEFAULT NULL,
  `terms` text NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `ads_grocer`
--

CREATE TABLE `ads_grocer` (
  `ads_id` int(11) NOT NULL,
  `grocer_id` int(11) NOT NULL,
  `grocer_location_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `advertiser`
--

CREATE TABLE `advertiser` (
  `id` int(11) NOT NULL,
  `guid` varchar(50) NOT NULL,
  `fullname` varchar(150) NOT NULL,
  `email` varchar(150) NOT NULL,
  `password` varchar(150) NOT NULL,
  `mobile` varchar(20) NOT NULL,
  `company` varchar(150) NOT NULL,
  `address` varchar(200) NOT NULL,
  `postcode` varchar(10) NOT NULL,
  `city` varchar(100) NOT NULL,
  `state` varchar(50) NOT NULL,
  `total_credits` int(11) NOT NULL,
  `remaining_credits` int(11) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `campaign`
--

CREATE TABLE `campaign` (
  `id` int(11) NOT NULL,
  `guid` varchar(150) NOT NULL,
  `advertiser_id` int(11) NOT NULL,
  `name` varchar(150) NOT NULL,
  `description` text NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `cashout_transactions`
--

CREATE TABLE `cashout_transactions` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `user_guid` varchar(40) NOT NULL,
  `transaction_guid` varchar(40) NOT NULL,
  `bank_account_holder_name` varchar(100) DEFAULT NULL,
  `bank_account_number` varchar(50) DEFAULT NULL,
  `bank_name` varchar(50) DEFAULT NULL,
  `bank_country` varchar(50) DEFAULT NULL,
  `remark_title` varchar(255) DEFAULT NULL,
  `remark_body` text,
  `transfer_date` datetime DEFAULT NULL,
  `receipt_image` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `category`
--

CREATE TABLE `category` (
  `id` int(11) NOT NULL,
  `guid` varchar(100) NOT NULL,
  `img` varchar(255) NOT NULL,
  `name` varchar(150) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `deal_cashbacks`
--

CREATE TABLE `deal_cashbacks` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `user_guid` varchar(40) NOT NULL,
  `shopping_list_guid` varchar(40) NOT NULL,
  `deal_guid` varchar(40) NOT NULL,
  `deal_cashback_transaction_guid` varchar(40) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `deal_cashback_status`
--

CREATE TABLE `deal_cashback_status` (
  `id` int(11) NOT NULL,
  `guid` varchar(100) NOT NULL,
  `deal_cashback_guid` varchar(100) NOT NULL,
  `status` varchar(100) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `deal_cashback_transactions`
--

CREATE TABLE `deal_cashback_transactions` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `user_guid` varchar(40) NOT NULL,
  `transaction_guid` varchar(40) NOT NULL,
  `receipt_url` varchar(255) NOT NULL,
  `verification_date` timestamp NULL DEFAULT NULL,
  `remark_title` varchar(255) DEFAULT NULL,
  `remark_body` text,
  `status` varchar(100) NOT NULL DEFAULT 'pending cleaning',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `default_shopping_lists`
--

CREATE TABLE `default_shopping_lists` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `occasion_guid` varchar(40) NOT NULL,
  `name` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `default_shopping_list_items`
--

CREATE TABLE `default_shopping_list_items` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `shopping_list_guid` varchar(40) NOT NULL,
  `name` varchar(255) NOT NULL,
  `category` varchar(255) NOT NULL,
  `subcategory` varchar(255) NOT NULL,
  `quantity` int(6) NOT NULL,
  `remark` text NOT NULL,
  `added_to_cart` int(1) UNSIGNED NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `default_shopping_list_item_images`
--

CREATE TABLE `default_shopping_list_item_images` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `shopping_list_guid` varchar(40) NOT NULL,
  `shopping_list_item_guid` varchar(40) NOT NULL,
  `url` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `devices`
--

CREATE TABLE `devices` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(255) NOT NULL,
  `user_guid` varchar(255) DEFAULT NULL,
  `uuid` varchar(255) NOT NULL,
  `os` varchar(100) DEFAULT NULL,
  `model` varchar(255) DEFAULT NULL,
  `push_token` text NOT NULL,
  `app_version` varchar(20) NOT NULL,
  `token_expired` int(1) UNSIGNED NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `event`
--

CREATE TABLE `event` (
  `id` int(11) NOT NULL,
  `guid` varchar(100) NOT NULL,
  `name` varchar(255) NOT NULL,
  `color` varchar(7) NOT NULL,
  `img` varchar(255) NOT NULL,
  `start_date` date NOT NULL,
  `end_date` date NOT NULL,
  `status` varchar(50) NOT NULL DEFAULT 'draft',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `event_deal`
--

CREATE TABLE `event_deal` (
  `event_id` int(11) NOT NULL,
  `deal_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `generic`
--

CREATE TABLE `generic` (
  `id` int(11) NOT NULL,
  `guid` varchar(100) NOT NULL,
  `category_id` int(11) NOT NULL,
  `subcategory_id` int(11) NOT NULL,
  `name` varchar(150) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `goose_db_version`
--

CREATE TABLE `goose_db_version` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `version_id` bigint(20) NOT NULL,
  `is_applied` tinyint(1) NOT NULL,
  `tstamp` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `grocer`
--

CREATE TABLE `grocer` (
  `id` int(11) NOT NULL,
  `guid` varchar(255) NOT NULL,
  `img` varchar(255) NOT NULL,
  `name` varchar(150) NOT NULL,
  `email` varchar(255) NOT NULL,
  `status` varchar(50) NOT NULL DEFAULT 'draft',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `grocer_location`
--

CREATE TABLE `grocer_location` (
  `id` int(11) NOT NULL,
  `guid` varchar(150) NOT NULL,
  `grocer_id` int(11) NOT NULL,
  `name` varchar(150) NOT NULL,
  `lat` double NOT NULL,
  `lng` double NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `item`
--

CREATE TABLE `item` (
  `id` int(11) NOT NULL,
  `guid` varchar(100) NOT NULL,
  `generic_id` int(11) DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `category_id` int(11) NOT NULL,
  `subcategory_id` int(11) NOT NULL,
  `remarks` varchar(255) NOT NULL,
  `status` int(11) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `occasions`
--

CREATE TABLE `occasions` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `slug` varchar(50) NOT NULL,
  `name` varchar(255) NOT NULL,
  `image` varchar(255) NOT NULL,
  `active` int(1) UNSIGNED DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `receipt`
--

CREATE TABLE `receipt` (
  `id` int(11) NOT NULL,
  `guid` varchar(100) NOT NULL,
  `deal_cashback_transaction_guid` varchar(100) NOT NULL,
  `name` varchar(100) NOT NULL,
  `outlet` varchar(100) NOT NULL,
  `cashier` varchar(100) NOT NULL,
  `receipt_no` varchar(100) NOT NULL,
  `date` date NOT NULL,
  `time` varchar(100) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `receipt_item`
--

CREATE TABLE `receipt_item` (
  `id` int(11) NOT NULL,
  `guid` varchar(100) NOT NULL,
  `receipt_guid` varchar(100) NOT NULL,
  `name` varchar(255) NOT NULL,
  `quantity` int(11) NOT NULL,
  `price` decimal(10,0) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `referral_cashback_transactions`
--

CREATE TABLE `referral_cashback_transactions` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL DEFAULT '',
  `user_guid` varchar(40) NOT NULL DEFAULT '',
  `referrer_guid` varchar(40) NOT NULL DEFAULT '',
  `transaction_guid` varchar(40) NOT NULL DEFAULT '',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `refreshperiod`
--

CREATE TABLE `refreshperiod` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `settings`
--

CREATE TABLE `settings` (
  `id` int(11) NOT NULL,
  `guid` varchar(100) NOT NULL,
  `name` varchar(255) NOT NULL,
  `slug` varchar(100) NOT NULL,
  `value` text NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `shopping_lists`
--

CREATE TABLE `shopping_lists` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `user_guid` varchar(40) NOT NULL,
  `occasion_guid` varchar(40) NOT NULL,
  `name` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `shopping_list_items`
--

CREATE TABLE `shopping_list_items` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `user_guid` varchar(40) NOT NULL,
  `shopping_list_guid` varchar(40) NOT NULL,
  `name` varchar(255) NOT NULL,
  `category` varchar(255) NOT NULL,
  `sub_category` varchar(255) NOT NULL,
  `quantity` int(6) NOT NULL,
  `remark` text NOT NULL,
  `added_from_deal` int(1) NOT NULL DEFAULT '0',
  `deal_guid` varchar(40) DEFAULT NULL,
  `cashback_amount` decimal(4,2) DEFAULT NULL,
  `deal_expired` int(1) DEFAULT NULL,
  `added_to_cart` int(1) UNSIGNED NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `shopping_list_item_images`
--

CREATE TABLE `shopping_list_item_images` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `user_guid` varchar(40) NOT NULL,
  `shopping_list_guid` varchar(40) NOT NULL,
  `shopping_list_item_guid` varchar(40) NOT NULL,
  `url` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `sms_histories`
--

CREATE TABLE `sms_histories` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(255) NOT NULL,
  `user_guid` varchar(255) NOT NULL,
  `provider` varchar(255) NOT NULL,
  `sms_id` varchar(100) NOT NULL,
  `text` varchar(255) NOT NULL,
  `recipient_no` varchar(20) NOT NULL,
  `verification_code` varchar(255) NOT NULL,
  `event` varchar(10) NOT NULL,
  `status` int(1) UNSIGNED NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `subcategory`
--

CREATE TABLE `subcategory` (
  `id` int(11) NOT NULL,
  `guid` varchar(100) NOT NULL,
  `name` varchar(150) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `transactions`
--

CREATE TABLE `transactions` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `user_guid` varchar(40) NOT NULL,
  `transaction_type_guid` varchar(40) NOT NULL,
  `transaction_status_guid` varchar(40) NOT NULL,
  `read_status` int(1) UNSIGNED NOT NULL DEFAULT '0',
  `reference_id` varchar(20) NOT NULL,
  `total_amount` decimal(9,2) NOT NULL,
  `approved_amount` decimal(9,2) DEFAULT NULL,
  `rejected_amount` decimal(9,2) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `transaction_statuses`
--

CREATE TABLE `transaction_statuses` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `slug` varchar(40) NOT NULL,
  `name` varchar(40) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `transaction_types`
--

CREATE TABLE `transaction_types` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(40) NOT NULL,
  `slug` varchar(40) NOT NULL,
  `name` varchar(40) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` int(10) UNSIGNED NOT NULL,
  `guid` varchar(255) NOT NULL,
  `facebook_id` varchar(100) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `email` varchar(255) NOT NULL,
  `phone_no` varchar(20) NOT NULL,
  `profile_picture` varchar(255) DEFAULT NULL,
  `referral_code` varchar(20) NOT NULL DEFAULT '',
  `bank_country` varchar(50) DEFAULT NULL,
  `bank_name` varchar(50) DEFAULT NULL,
  `bank_account_name` varchar(50) DEFAULT NULL,
  `bank_account_number` varchar(50) DEFAULT NULL,
  `wallet` decimal(9,2) NOT NULL DEFAULT '0.00',
  `register_by` varchar(20) NOT NULL,
  `verified` int(1) UNSIGNED NOT NULL DEFAULT '0',
  `blacklist` int(1) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `wallet_transaction_log`
--

CREATE TABLE `wallet_transaction_log` (
  `id` int(11) NOT NULL,
  `guid` varchar(100) NOT NULL,
  `user_guid` varchar(100) NOT NULL,
  `transaction_guid` varchar(100) NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `type` varchar(255) NOT NULL,
  `approve_by` varchar(100) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `admin`
--
ALTER TABLE `admin`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `ads`
--
ALTER TABLE `ads`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `status_start_end_date` (`status`,`start_date`,`end_date`),
  ADD KEY `status` (`status`,`start_date`,`end_date`,`quota`,`perlimit`);

--
-- Indexes for table `ads_grocer`
--
ALTER TABLE `ads_grocer`
  ADD KEY `ads_id` (`ads_id`),
  ADD KEY `grocer_id` (`grocer_id`),
  ADD KEY `grocer_location_id` (`grocer_location_id`);

--
-- Indexes for table `advertiser`
--
ALTER TABLE `advertiser`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `campaign`
--
ALTER TABLE `campaign`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `cashout_transactions`
--
ALTER TABLE `cashout_transactions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `user_guid` (`user_guid`),
  ADD KEY `transaction_guid` (`transaction_guid`);

--
-- Indexes for table `category`
--
ALTER TABLE `category`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`);

--
-- Indexes for table `deal_cashbacks`
--
ALTER TABLE `deal_cashbacks`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `deal_guid` (`deal_guid`),
  ADD KEY `deal_cashback_transaction_guid` (`deal_cashback_transaction_guid`),
  ADD KEY `user_guid_deal_guid` (`user_guid`,`deal_guid`) USING BTREE,
  ADD KEY `user_deal_shopping_list_guid` (`user_guid`,`shopping_list_guid`,`deal_guid`) USING BTREE,
  ADD KEY `	deal_cashback_transaction_shopping_list_guid` (`deal_cashback_transaction_guid`,`shopping_list_guid`);

--
-- Indexes for table `deal_cashback_status`
--
ALTER TABLE `deal_cashback_status`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `deal_cashback_transactions`
--
ALTER TABLE `deal_cashback_transactions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `user_guid` (`user_guid`),
  ADD KEY `transaction_guid` (`transaction_guid`),
  ADD KEY `status` (`status`),
  ADD KEY `guid_transaction_guid` (`guid`,`transaction_guid`);

--
-- Indexes for table `default_shopping_lists`
--
ALTER TABLE `default_shopping_lists`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `occasion_guid` (`occasion_guid`);

--
-- Indexes for table `default_shopping_list_items`
--
ALTER TABLE `default_shopping_list_items`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `shopping_list_guid` (`shopping_list_guid`);

--
-- Indexes for table `default_shopping_list_item_images`
--
ALTER TABLE `default_shopping_list_item_images`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `shopping_list_guid` (`shopping_list_guid`),
  ADD KEY `shopping_list_item_guid` (`shopping_list_item_guid`);

--
-- Indexes for table `devices`
--
ALTER TABLE `devices`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD UNIQUE KEY `uuid` (`uuid`),
  ADD KEY `uuid_deleted_at` (`uuid`,`deleted_at`) USING BTREE,
  ADD KEY `uuid_user_guid_deleted_at` (`user_guid`,`uuid`,`deleted_at`);

--
-- Indexes for table `event`
--
ALTER TABLE `event`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `status` (`status`),
  ADD KEY `start_date` (`start_date`),
  ADD KEY `end_date` (`end_date`);

--
-- Indexes for table `event_deal`
--
ALTER TABLE `event_deal`
  ADD KEY `event_id` (`event_id`),
  ADD KEY `deal_id` (`deal_id`);

--
-- Indexes for table `generic`
--
ALTER TABLE `generic`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `category_id` (`category_id`),
  ADD KEY `subcategory_id` (`subcategory_id`),
  ADD KEY `name` (`name`);

--
-- Indexes for table `goose_db_version`
--
ALTER TABLE `goose_db_version`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `id` (`id`);

--
-- Indexes for table `grocer`
--
ALTER TABLE `grocer`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`);

--
-- Indexes for table `grocer_location`
--
ALTER TABLE `grocer_location`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `grocer_id` (`grocer_id`);

--
-- Indexes for table `item`
--
ALTER TABLE `item`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`);

--
-- Indexes for table `occasions`
--
ALTER TABLE `occasions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `slug` (`slug`),
  ADD UNIQUE KEY `guid` (`guid`,`deleted_at`) USING BTREE;

--
-- Indexes for table `receipt`
--
ALTER TABLE `receipt`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `receipt_item`
--
ALTER TABLE `receipt_item`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `referral_cashback_transactions`
--
ALTER TABLE `referral_cashback_transactions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD UNIQUE KEY `transaction_guid` (`transaction_guid`),
  ADD KEY `idx_referral_cashbacks_referrer_guid` (`user_guid`),
  ADD KEY `idx_referral_cashbacks_referent_guid` (`referrer_guid`),
  ADD KEY `user_guid` (`user_guid`);

--
-- Indexes for table `refreshperiod`
--
ALTER TABLE `refreshperiod`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `settings`
--
ALTER TABLE `settings`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `shopping_lists`
--
ALTER TABLE `shopping_lists`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `occasion_guid` (`occasion_guid`),
  ADD KEY `user_guid_deleted_at` (`user_guid`,`deleted_at`) USING BTREE;

--
-- Indexes for table `shopping_list_items`
--
ALTER TABLE `shopping_list_items`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `user_guid` (`user_guid`),
  ADD KEY `category` (`category`),
  ADD KEY `sub_category` (`sub_category`),
  ADD KEY `added_to_cart` (`added_to_cart`),
  ADD KEY `shopping_list_guid_deleted_at` (`shopping_list_guid`,`deleted_at`) USING BTREE;

--
-- Indexes for table `shopping_list_item_images`
--
ALTER TABLE `shopping_list_item_images`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `shopping_list_guid` (`shopping_list_guid`),
  ADD KEY `shopping_list_item_guid` (`shopping_list_item_guid`);

--
-- Indexes for table `sms_histories`
--
ALTER TABLE `sms_histories`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid_idx` (`guid`) USING BTREE,
  ADD UNIQUE KEY `recipient_no_verification_code_event` (`recipient_no`,`verification_code`,`event`),
  ADD KEY `idx_devices_deleted_at` (`deleted_at`),
  ADD KEY `recipient_no_verification_code_idx` (`recipient_no`,`verification_code`) USING BTREE,
  ADD KEY `recipient_no_idx` (`deleted_at`) USING BTREE,
  ADD KEY `recipient_no_created_at_event_idx` (`recipient_no`,`created_at`,`event`);

--
-- Indexes for table `subcategory`
--
ALTER TABLE `subcategory`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `name` (`name`);

--
-- Indexes for table `transactions`
--
ALTER TABLE `transactions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `guid_type_status` (`transaction_type_guid`,`transaction_status_guid`,`user_guid`) USING BTREE,
  ADD KEY `user_guid_deleted_at` (`user_guid`,`deleted_at`) USING BTREE;

--
-- Indexes for table `transaction_statuses`
--
ALTER TABLE `transaction_statuses`
  ADD PRIMARY KEY (`id`),
  ADD KEY `guid` (`guid`,`deleted_at`) USING BTREE,
  ADD KEY `slug_deleted_at` (`slug`) USING BTREE;

--
-- Indexes for table `transaction_types`
--
ALTER TABLE `transaction_types`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `guid` (`guid`),
  ADD KEY `slug` (`slug`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `phone_no` (`phone_no`) USING HASH,
  ADD UNIQUE KEY `guid_deleted_at` (`guid`) USING BTREE,
  ADD KEY `facebook_id_deleted_at` (`facebook_id`,`deleted_at`) USING BTREE,
  ADD KEY `deleted_at` (`deleted_at`) USING BTREE,
  ADD KEY `referral_code_deleted_at` (`referral_code`) USING BTREE;

--
-- Indexes for table `wallet_transaction_log`
--
ALTER TABLE `wallet_transaction_log`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `admin`
--
ALTER TABLE `admin`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `ads`
--
ALTER TABLE `ads`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `advertiser`
--
ALTER TABLE `advertiser`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `campaign`
--
ALTER TABLE `campaign`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `cashout_transactions`
--
ALTER TABLE `cashout_transactions`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `category`
--
ALTER TABLE `category`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `deal_cashbacks`
--
ALTER TABLE `deal_cashbacks`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `deal_cashback_status`
--
ALTER TABLE `deal_cashback_status`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `deal_cashback_transactions`
--
ALTER TABLE `deal_cashback_transactions`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `default_shopping_lists`
--
ALTER TABLE `default_shopping_lists`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `default_shopping_list_items`
--
ALTER TABLE `default_shopping_list_items`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `default_shopping_list_item_images`
--
ALTER TABLE `default_shopping_list_item_images`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `devices`
--
ALTER TABLE `devices`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `event`
--
ALTER TABLE `event`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `generic`
--
ALTER TABLE `generic`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `goose_db_version`
--
ALTER TABLE `goose_db_version`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `grocer`
--
ALTER TABLE `grocer`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `grocer_location`
--
ALTER TABLE `grocer_location`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `item`
--
ALTER TABLE `item`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16357;
--
-- AUTO_INCREMENT for table `occasions`
--
ALTER TABLE `occasions`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `receipt`
--
ALTER TABLE `receipt`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `receipt_item`
--
ALTER TABLE `receipt_item`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `referral_cashback_transactions`
--
ALTER TABLE `referral_cashback_transactions`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `refreshperiod`
--
ALTER TABLE `refreshperiod`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `settings`
--
ALTER TABLE `settings`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `shopping_lists`
--
ALTER TABLE `shopping_lists`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `shopping_list_items`
--
ALTER TABLE `shopping_list_items`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `shopping_list_item_images`
--
ALTER TABLE `shopping_list_item_images`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `sms_histories`
--
ALTER TABLE `sms_histories`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `subcategory`
--
ALTER TABLE `subcategory`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `transactions`
--
ALTER TABLE `transactions`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `transaction_statuses`
--
ALTER TABLE `transaction_statuses`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `transaction_types`
--
ALTER TABLE `transaction_types`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `wallet_transaction_log`
--
ALTER TABLE `wallet_transaction_log`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
