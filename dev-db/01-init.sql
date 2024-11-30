CREATE TABLE IF NOT EXISTS "games" ("id" varchar(36),"title" text NOT NULL,"game_order" bigint NOT NULL,"story" text NOT NULL,"cheatsheet" text NOT NULL,"max_score" bigint NOT NULL,"description" text NOT NULL,"background" text NOT NULL,"winning_message" text NOT NULL,"wrong_attempt_cost" bigint NOT NULL,"perfect_timeslot" bigint NOT NULL,"great_timeslot" bigint NOT NULL,"medium_timeslot" bigint NOT NULL,"not_so_good_timeslot" bigint NOT NULL,"textual_hint_price" bigint NOT NULL,"textual_hint" text NOT NULL,"hint_solution_price" bigint NOT NULL,"time_freeze_price" bigint NOT NULL,"time_freeze_duration" bigint NOT NULL,PRIMARY KEY ("id"));

CREATE TABLE IF NOT EXISTS "players" ("id" varchar(36),"username" text NOT NULL UNIQUE,"password" text NOT NULL,"email" text NOT NULL UNIQUE,"icon_id" text,PRIMARY KEY ("id"));

CREATE TABLE IF NOT EXISTS "blocks" ("id" varchar(36),"content" text NOT NULL,"order" bigint,"skeleton" boolean,"game_id" varchar(36) NOT NULL,PRIMARY KEY ("id"),CONSTRAINT "fk_games_blocks" FOREIGN KEY ("game_id") REFERENCES "games"("id"));

CREATE TABLE IF NOT EXISTS "player_games" ("player_id" varchar(36) NOT NULL,"game_id" varchar(36) NOT NULL,"score" bigint,"attempts" bigint NOT NULL DEFAULT 0,"start_time" timestamptz NOT NULL,"end_time" timestamptz,"textual_hint_points_used" bigint NOT NULL DEFAULT 0,"hint_solution_points_used" bigint NOT NULL DEFAULT 0,"time_freeze_points_used" bigint NOT NULL DEFAULT 0,CONSTRAINT "fk_games_player_games" FOREIGN KEY ("game_id") REFERENCES "games"("id"),CONSTRAINT "fk_players_player_games" FOREIGN KEY ("player_id") REFERENCES "players"("id"));

CREATE TABLE IF NOT EXISTS "icons" ("id" varchar(36),"svg" text NOT NULL,PRIMARY KEY ("id"));

