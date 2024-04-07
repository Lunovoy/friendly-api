CREATE TABLE IF NOT EXISTS "tag" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "title" varchar(30) unique not null,
    "user_id" UUID not null,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);