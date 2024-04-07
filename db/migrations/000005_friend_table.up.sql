CREATE TABLE IF NOT EXISTS "friend" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "first_name" varchar(50) not null,
    "last_name" varchar(50) DEFAULT '',
    "middle_name" varchar(50) DEFAULT '',
    "dob" timestamp with time zone DEFAULT '1970-01-01 00:00:00+00', 
    "user_id" UUID not null,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "friendlists_friends" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "friendlist_id" UUID not null,
    "friend_id" UUID not null,
    FOREIGN KEY ("friendlist_id") REFERENCES "friendlist" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("friend_id") REFERENCES "friend" ("id") ON DELETE CASCADE
);