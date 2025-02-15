CREATE TABLE "activity" (
                            "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                            "user_id" UUID,
                            "name" varchar(255),
                            "calories_per_hour" float(8),
                            "duration_minutes" float(8),
                            "total_calories" float(8),
                            "created_at" timestamp DEFAULT (now()),
                            "updated_at" timestamp DEFAULT null
);

CREATE TABLE "activity_user" (
                               "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                               "user_id" UUID NOT NULL
                                 REFERENCES "users"(id) ON DELETE CASCADE,
                               "activity_id" UUID NOT NULL
                                 REFERENCES "activity"(id) ON DELETE CASCADE,
                               "created_at" timestamp DEFAULT (now()),
                               "updated_at" timestamp DEFAULT null,
                               UNIQUE ("user_id", "activity_id")
);



CREATE INDEX idx_activity_id ON activity (id);
CREATE INDEX idx_activity_user_id ON activity (user_id);
CREATE INDEX idx_activity_name ON activity (name);

CREATE INDEX idx_activity_u_id ON activity_user (id);
CREATE INDEX idx_activity_u_user_id ON activity_user (user_id);

CREATE TABLE "exercise_stats" (
                                "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                "user_id" UUID NOT NULL REFERENCES users(id),
                                "activity_id" UUID NOT NULL REFERENCES activity(id),
                                "session_name" varchar(255),
                                "number_of_times" integer,
                                "total_duration_seconds" integer,  -- store total durations in seconds, simpler
                                "total_calories_burned" integer,
                                "created_at" timestamp DEFAULT now(),
                                "updated_at" timestamp DEFAULT now(),
                                UNIQUE (user_id, activity_id)
);
