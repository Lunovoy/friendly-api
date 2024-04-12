CREATE TABLE IF NOT EXISTS "work_info" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "country" varchar(100) DEFAULT '',
    "city" varchar(100) DEFAULT '',
    "company" varchar(100) DEFAULT '',
    "position" varchar(100) DEFAULT '',
    "messenger" varchar(100) DEFAULT '',
    "communication_method" varchar(100) DEFAULT '',
    "nationality" varchar(50) DEFAULT '',
    "language" varchar(100) DEFAULT '',
    "friend_id" UUID not null,
    FOREIGN KEY ("friend_id") REFERENCES "friend" ("id") ON DELETE CASCADE
);