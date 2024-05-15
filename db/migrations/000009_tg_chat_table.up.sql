CREATE TABLE IF NOT EXISTS "tg_chat" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "chat_id" bigint not null,
    "user_id" UUID unique not null,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);