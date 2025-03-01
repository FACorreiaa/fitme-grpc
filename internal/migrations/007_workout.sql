CREATE TABLE IF NOT EXISTS "workout_plan" (
                                              "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                              "user_id" UUID not null,
                                              "description" varchar(255),
                                              "notes" varchar(255),
                                              "rating" integer DEFAULT 10,
                                              "created_at" timestamp DEFAULT (now()),
                                              "updated_at" timestamp DEFAULT null,
                                              FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS "workout_plan_detail" (
                                                     "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                                     "workout_plan_id" UUID,
                                                     "day" varchar(100),
                                                     "exercises" uuid[],
                                                     "created_at" timestamp DEFAULT (now()),
                                                     FOREIGN KEY (workout_plan_id) REFERENCES workout_plan(id)

);




CREATE TABLE IF NOT EXISTS "workout_day" (
                                             "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                             "workout_plan_id" UUID,
                                             "day" varchar(255),
                                             "created_at" timestamp DEFAULT (now()),
                                             "updated_at" timestamp DEFAULT null,
                                             FOREIGN KEY (workout_plan_id) REFERENCES workout_plan(id)
);

CREATE TABLE IF NOT EXISTS workout_day_exercise (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  workout_day_id UUID NOT NULL,
  exercise_id UUID NOT NULL,
  series INTEGER,
  repetitions VARCHAR(50),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT null,
  FOREIGN KEY (workout_day_id) REFERENCES workout_day(id),
  FOREIGN KEY (exercise_id) REFERENCES exercise_list(id)
  );
