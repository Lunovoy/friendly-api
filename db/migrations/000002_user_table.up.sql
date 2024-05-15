CREATE TABLE IF NOT EXISTS "user" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "username" varchar(50) DEFAULT '',
    "first_name" varchar(50) DEFAULT '',
    "last_name" varchar(50) DEFAULT '',
    "middle_name" varchar(50) DEFAULT '',
    "tg_username" varchar(50) DEFAULT '',
    "mail" varchar(100) unique not null,
    "password_hash" varchar(255) not null,
    "salt" varchar(255) not null
);