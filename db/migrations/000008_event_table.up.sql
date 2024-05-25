CREATE TABLE IF NOT EXISTS "event" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "title" varchar(50) not null,
    "description" text,
    "start_date" timestamp with time zone,
    "end_date" timestamp with time zone,
    "frequency" varchar(50) not null,
    "is_active" boolean DEFAULT true,
    "user_id" UUID not null,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "friends_events"(
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "friend_id" UUID not null,
    "event_id" UUID not null,
    FOREIGN KEY ("friend_id") REFERENCES "friend" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("event_id") REFERENCES "event" ("id") ON DELETE CASCADE
);