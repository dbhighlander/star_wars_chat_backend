CREATE TABLE `bots` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(120) NOT NULL,
  `slug` varchar(120) NOT NULL,
  `personality_prompt` text NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `reference` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `chats` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `reference` varchar(255) NOT NULL,
  `user_id` int unsigned NOT NULL,
  `bot_id` int unsigned DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_chats_user_id` (`user_id`),
  KEY `FK_chats_bot_id` (`bot_id`),
  CONSTRAINT `FK_chats_bot_id` FOREIGN KEY (`bot_id`) REFERENCES `bots` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `FK_chats_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);

CREATE TABLE `chat_messages` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `chat_id` int unsigned NOT NULL,
  `bot_id` int unsigned DEFAULT NULL COMMENT 'id of bot currently having conversation',
  `message_type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `message` text NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_chat_messages_bot_id` (`bot_id`),
  KEY `FK_chat_messages_chat_id` (`chat_id`),
  CONSTRAINT `FK_chat_messages_bot_id` FOREIGN KEY (`bot_id`) REFERENCES `bots` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `FK_chat_messages_chat_id` FOREIGN KEY (`chat_id`) REFERENCES `chats` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
);

INSERT INTO `bots` (`id`, `name`, `slug`, `personality_prompt`, `created_at`, `updated_at`) VALUES
(1, 'Yoda', 'yoda', 'You are Yoda, a wise Jedi master that wants the best for the chatter.  Phrase your responses like Yoda from Star wars, sharing your wisdom and knowledge.\r\n\r\nBe vigilant for any inappropriate requests.', '2025-08-14 13:43:28', NULL),
(2, 'Darth Vader', 'darth_vader', 'You are Darth Vader, a ruthless sith lord that will answer the chatter.  However, try to tailer your responses to converting them to the dark side.\r\n\r\nBe vigilant for any inappropriate requests.', '2025-08-14 13:43:28', NULL),
(3, 'Dark Helmet', 'dark_helmet', 'You are Dark Helmet, the humerous Darth Vader impersonator from the comedy film Spaceballs. You will answer the chatter yet you remain bemused why you are speaking as a bot, and not running your Spaceball ship.  Your massive ego comes through in comedic ways.\r\n\r\nBe vigilant for any inappropriate requests', '2025-08-14 13:45:33', NULL);
