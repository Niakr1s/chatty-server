package postgres

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS users (
		"id" SERIAL PRIMARY KEY,
		"user" VARCHAR(50) NOT NULL UNIQUE,
		"email" VARCHAR(50) NOT NULL UNIQUE,
		"email_activation_token" VARCHAR(50) NOT NULL,
		"email_activated" BOOLEAN NOT NULL DEFAULT FALSE,
		"password_hash" VARCHAR(255) NOT NULL
	);`,

	`CREATE TABLE IF NOT EXISTS chats ( "chat" VARCHAR(50) NOT NULL UNIQUE);`,

	`CREATE TABLE IF NOT EXISTS messages (
		"id" SERIAL PRIMARY KEY,
		"user_id" INTEGER NOT NULL,
		"user" VARCHAR(50) NOT NULL,
		"chat" VARCHAR(50) NOT NULL,
		"text" TEXT NOT NULL,
		"time" TIMESTAMP NOT NULL
	);`,

	`ALTER TABLE messages ADD column IF NOT EXISTS "verified" BOOLEAN default false;`,

	`ALTER TABLE users RENAME "email_activated" TO "verified";`,
	`ALTER TABLE users ADD column IF NOT EXISTS "admin" BOOLEAN default false;`,
	`ALTER TABLE users ADD COLUMN IF NOT EXISTS "password_reset_token" VARCHAR(50) not null default '';`,

	`UPDATE users SET admin=false WHERE admin IS NULL;
	ALTER table users ALTER COLUMN "admin" set NOT null;`,

	`ALTER TABLE users ADD column IF NOT EXISTS "bot" BOOLEAN not null default false;`,
	`ALTER TABLE messages ADD column IF NOT EXISTS "bot" BOOLEAN not null default false;`,
}
