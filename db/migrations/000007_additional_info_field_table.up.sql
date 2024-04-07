CREATE TABLE IF NOT EXISTS "additional_info_field" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "title" varchar(50) unique not null,
    "user_id" UUID not null,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "friends_additional_info_fields" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "friend_id" UUID not null,
    "additional_info_field_id" UUID not null,
    FOREIGN KEY ("friend_id") REFERENCES "friend" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("additional_info_field_id") REFERENCES "additional_info_field" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "additional_info_field_text" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "content" varchar(255) DEFAULT '',
    "additional_info_field_id" UUID not null,
    "friend_id" UUID not null,
    FOREIGN KEY ("additional_info_field_id") REFERENCES "additional_info_field" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("friend_id") REFERENCES "friend" ("id") ON DELETE CASCADE
);