


CREATE TABLE "total_exercise_session" (
                                        "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                        "user_id" UUID UNIQUE,
                                        "activity_id" UUID,
                                        "total_duration_hours" integer,
                                        "total_duration_minutes" integer,
                                        "total_duration_seconds" integer,
                                        "total_calories_burned" integer,
                                        "session_name" varchar(255),
                                        "created_at" timestamp DEFAULT (now()),
                                        "updated_at" timestamp
);

CREATE TABLE "exercise_session" (
                                  "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                  "user_id" UUID,
                                  "activity_id" UUID,
                                  "session_name" varchar(255),
                                  "start_time" timestamp,
                                  "end_time" timestamp,
                                  "duration_seconds" integer,
                                  "duration_minutes" integer,
                                  "duration_hours" integer,
                                  calories_burned float(8),
                                  "created_at" timestamp DEFAULT (now())
);
