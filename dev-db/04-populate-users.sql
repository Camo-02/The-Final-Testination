INSERT INTO "players" ("id", "username", "email", "password", "icon_id") VALUES
    ('userId1', 'username1', 'email@user1.com', 'password1', '1'),
    ('userId2', 'username2', 'email@user2.com', 'password2','1'),
    ('userId3', 'username3', 'email@user3.com', 'password3','1'),
    ('userId4', 'username4', 'email@user4.com', 'password4','1'),
    ('userId5', 'username5', 'email@user5.com', 'password5','1'),
    ('userId6', 'username6', 'email@user6.com', 'password6','1'),
    ('userId7', 'username7', 'email@user7.com', 'password7','1'),
    ('userId8', 'username8', 'email@user8.com', 'password8','1'),
    ('userId9', 'username9', 'email@user9.com', 'password9','1'),
    ('userId10', 'username10', 'email@user10.com', 'password10','1'),
    ('userId11', 'username11', 'email@user11.com', 'password11','1'),
    ('userId12', 'username12', 'email@user12.com', 'password12','1'),
    ('userId13', 'username13', 'email@user13.com', 'password13','1'),
    ('userId14', 'username14', 'email@user14.com', 'password14','1'),
    ('userId15', 'username15', 'email@user15.com', 'password15','1'),
    ('userId31', 'username31', 'email@user31.com', 'password31','1'),
    ('userId32', 'username32', 'email@user32.com', 'password32','1'),
    ('userId33', 'username33', 'email@user33.com', 'password33','1'),
    ('userId34', 'username34', 'email@user34.com', 'password34','1'),
    ('userId35', 'username35', 'email@user35.com', 'password35','1'),
    ('userId36', 'username36', 'email@user36.com', 'password36','1'),
    ('userId37', 'username37', 'email@user37.com', 'password37','1'),
    ('userId38', 'username38', 'email@user38.com', 'password38','1'),
    ('userId39', 'username39', 'email@user39.com', 'password39','1'),
    ('userId40', 'username40', 'email@user40.com', 'password40','1'),
    ('userId41', 'username41', 'email@user41.com', 'password41','1'),
    ('userId42', 'username42', 'email@user42.com', 'password42','1'),
    ('userId43', 'username43', 'email@user43.com', 'password43','1'),
    ('userId44', 'username44', 'email@user44.com', 'password44','1'),
    ('userId45', 'username45', 'email@user45.com', 'password45','1'),
    ('userId46', 'username46', 'email@user46.com', 'password46','1'),
    ('userId47', 'username47', 'email@user47.com', 'password47','1'),
    ('userId48', 'username48', 'email@user48.com', 'password48','1'),
    ('userId49', 'username49', 'email@user49.com', 'password49','1'),
    ('userId50', 'username50', 'email@user50.com', 'password50','1'),
    ('userId51', 'username51', 'email@user51.com', 'password51','1'),
    ('userId52', 'username52', 'email@user52.com', 'password52','1'),
    ('userId53', 'username53', 'email@user53.com', 'password53','1'),
    ('userId54', 'username54', 'email@user54.com', 'password54','1'),
    ('userId55', 'username55', 'email@user55.com', 'password55','1'),
    ('userId56', 'username56', 'email@user56.com', 'password56','1'),
    ('userId57', 'username57', 'email@user57.com', 'password57','1'),
    ('userId58', 'username58', 'email@user58.com', 'password58','1'),
    ('userId59', 'username59', 'email@user59.com', 'password59','1'),
    ('userId60', 'username60', 'email@user60.com', 'password60','1'),
    ('a977b9b6-00dd-43de-b9ad-bd1c41ff20be', 'admin', 'admin@testination.com', '$2a$10$JLcdkCOtHsXQ9IR1uNtMtu3w..glvlUpeK4hZ9.E0QK89y2sNvjaS','1'),
    ('7f57e6ae-c66c-42e7-b878-72c79ba131e9','newUser','newUser@email.com', '$2a$10$ykMkEEFbz/Pi/vGEjGMdseDJqXtlqBxoIg56GzHOZS5mdu16u09De','1'),
    ('6d4c437b-5803-4b08-890b-44383af74ab3', 'test', 'test@testination.com', '$2a$10$3NittdjoXds0rKX/52.zKuL6ytbpGh2K9wDD0l1evnAy0jpM2FbG.','1'),
    ('33c89d92-5a6d-40ca-9784-2d41a2e5feac', 'empty', 'empty@testination.com', '$2a$10$Uvl15seNqoeor5eKEIQbfucXWsO5lKTxK0fxM4in1YdpwL9xABRgy','1'),
    ('fbfbe74f-4d5b-4a45-9590-5d72e8696f71','profilePageTest','eee@Task.com','$2a$10$XnMdZ/TVT4.zF2DWuNAEU.iGLtRJ3RK0eTZc5MR50B9JU0SaWN0o.',''),
    ('aebc3093-dcb4-4cad-a569-1204470617aa','PleaseRunTests','EachTimeYouFinishA@Task.com','$2a$10$m4ZEEWU9pxntC5H0TpfHo.gK/ScBJKbz8LyvAGVyi62PFBi3juPmO',''),
    ('c977b9b6-10dd-43de-b2ad-bd3c41ff20be', 'luca', 'luca@gmail.com', '$2a$10$JLcdkCOtHsXQ9IR1uNtMtu3w..glvlUpeK4hZ9.E0QK89y2sNvjaS','1');

-- admin : rootroot , test : rootroot , newUser : rootroot, PleaseRunTests : rootroot



INSERT INTO "player_games" ("game_id", "player_id","score", "start_time") VALUES
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId3', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId4', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId5', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId6', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId7', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId8', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId9', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId10', 200, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId11', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId12', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId13', 200, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId15', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId31', 150, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId32', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId33', 200, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId34', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId35', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId36', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId37', 150, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId38', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId39', 200, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId40', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId41', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId42', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId43', 150, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId44', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId45', 200, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId46', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId47', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId48', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId49', 150, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId50', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId51', 200, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId52', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId53', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId54', 100, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId55', 150, '2024-02-24 10:00:00 UTC'),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId56', 100, '2024-02-24 10:00:00 UTC');

--playerGames with hints used
INSERT INTO "player_games" ("game_id", "player_id","score", "attempts", "start_time", "end_time", "textual_hint_points_used", "hint_solution_points_used", "time_freeze_points_used") VALUES
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId1', 100, 1, '2024-02-24 10:00:00', '2024-02-24 10:01:00', 20, 20, 20),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'userId2', 100, 1, '2024-02-24 10:00:00', '2024-02-24 10:01:00', 20, 0, 0),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'a977b9b6-00dd-43de-b9ad-bd1c41ff20be', 100, 1, '2024-02-24 10:00:00', '2024-02-24 10:01:00', 10, 0, 0),
    ('05732286-9fa5-45d4-bef3-13ae0d481afa', 'a977b9b6-00dd-43de-b9ad-bd1c41ff20be', 200, 1, '2024-02-24 10:00:00', '2024-02-24 10:01:00', 0, 0, 0),
    ('a76db50b-ee98-4dd4-9d63-4c0ab695ad5f', 'a977b9b6-00dd-43de-b9ad-bd1c41ff20be', 300, 1, '2024-02-24 10:00:00', '2024-02-24 10:01:00', 0, 0, 0),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', '6d4c437b-5803-4b08-890b-44383af74ab3', 100, 1, '2024-02-24 10:00:00', '2024-02-24 10:01:00', 0, 0, 0),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'fbfbe74f-4d5b-4a45-9590-5d72e8696f71', 100, 1, '2024-02-24 10:00:00', null, 20, 20, 20),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'aebc3093-dcb4-4cad-a569-1204470617aa', 0, 1, '2024-02-24 10:00:00', '2024-02-24 10:01:00', 0, 0, 0),
    ('af8e4754-1b84-4fec-bec4-154a3f894b8f', 'c977b9b6-10dd-43de-b2ad-bd3c41ff20be', 100, 1, '2024-02-24 10:00:00', '2024-02-24 10:01:00', 0, 0, 0),
    ('05732286-9fa5-45d4-bef3-13ae0d481afa', 'c977b9b6-10dd-43de-b2ad-bd3c41ff20be', 200, 1, '2024-02-24 10:00:00', '2024-02-24 10:01:00', 0, 0, 0);


