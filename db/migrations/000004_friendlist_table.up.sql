CREATE TABLE IF NOT EXISTS "friendlist" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "title" varchar(50) unique not null,
    "description" varchar(255) DEFAULT '',
    "image_id" UUID,
    "user_id" UUID not null,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "friendlists_tags" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "friendlist_id" UUID not null,
    "tag_id" UUID not null,
    FOREIGN KEY ("friendlist_id") REFERENCES "friendlist" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("tag_id") REFERENCES "tag" ("id") ON DELETE CASCADE
);