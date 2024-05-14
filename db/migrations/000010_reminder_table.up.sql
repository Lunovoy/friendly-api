CREATE TABLE IF NOT EXISTS "reminder" {
  "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  "minutes_until_event" integer not null,
  "event_id" UUID not null,
  "user_id" UUID not null,
  FOREIGN KEY ("event_id") REFERENCES "event" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
}