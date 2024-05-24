  ALTER TABLE IF EXISTS "user" 
    ADD COLUMN IF NOT EXISTS "country" varchar(100) DEFAULT '',
    ADD COLUMN IF NOT EXISTS   "city" varchar(100) DEFAULT '',
    ADD COLUMN IF NOT EXISTS   "company" varchar(100) DEFAULT '',
    ADD COLUMN IF NOT EXISTS   "profession" varchar(100) DEFAULT '',
    ADD COLUMN IF NOT EXISTS  "position" varchar(100) DEFAULT '',
    ADD COLUMN IF NOT EXISTS  "messenger" varchar(100) DEFAULT '',
    ADD COLUMN IF NOT EXISTS  "communication_method" varchar(100) DEFAULT '',
    ADD COLUMN IF NOT EXISTS   "nationality" varchar(50) DEFAULT '',
    ADD COLUMN IF NOT EXISTS   "resident" boolean DEFAULT false,
    ADD COLUMN IF NOT EXISTS   "language" varchar(100) DEFAULT '';
