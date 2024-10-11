CREATE TABLE "food" (
                        "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                        "name" varchar(255),
                        "calories" float(8),
                        "serving_size" float(8),
                        "protein" float(8),
                        "fat_total" float(8),
                        "fat_saturated" float(8),
                        "carbohydrates_total" float(8),
                        "fiber" float(8),
                        "sugar" float(8),
                        "sodium" float(8),
                        "potassium" float(8),
                        "cholesterol" float(8),
                        "created_at" timestamp DEFAULT (now()),
                        "updated_at" timestamp DEFAULT null
);
CREATE TABLE "meal_type" (
                             "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                             "user_id" UUID UNIQUE,
                             "food_id" UUID UNIQUE,
                             "meal_number" integer,
                             "meal_description" varchar(255),
                             "created_at" timestamp DEFAULT (now()),
                             "updated_at" timestamp DEFAULT null
);
CREATE TABLE "meal_plan" (
                             "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                             "user_id" UUID UNIQUE,
                             "meal_type_id" UUID UNIQUE,
                             "description" varchar(255),
                             "notes" varchar(255),
                             "total_calories" float(8),
                             "created_at" timestamp DEFAULT (now()),
                             "updated_at" timestamp DEFAULT null,
                             "rating" integer DEFAULT 10
);
CREATE TABLE "favourite" (
                             "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                             "user_id" UUID UNIQUE,
                             "exercise_id" UUID UNIQUE,
                             "activity_id" UUID UNIQUE,
                             "food_id" UUID UNIQUE,
                             "created_at" timestamp DEFAULT (now()),
                             "updated_at" timestamp DEFAULT null
);
CREATE TABLE "recipe" (
                          "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                          "user_id" UUID UNIQUE,
                          "food_id" UUID UNIQUE,
                          "description" varchar(255),
                          "created_at" timestamp DEFAULT (now()),
                          "updated_at" timestamp DEFAULT null
);
CREATE TABLE "recipe_user" (
                               "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                               "recipe_id" UUID UNIQUE,
                               "user_id" UUID UNIQUE,
                               "created_at" timestamp DEFAULT (now()),
                               "updated_at" timestamp DEFAULT null
);

CREATE TABLE "meal_plan_meal_type" (
                                       "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                       "meal_plan_id" UUID UNIQUE,
                                       "meal_type_id" UUID UNIQUE,
                                       "created_at" timestamp DEFAULT (now()),
                                       "updated_at" timestamp DEFAULT null
);
CREATE TABLE "meal_plan_user" (
                                  "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                  "meal_plan_id" UUID UNIQUE,
                                  "user_id" UUID UNIQUE,
                                  "created_at" timestamp DEFAULT (now()),
                                  "updated_at" timestamp DEFAULT null
);

CREATE TABLE "water_intake_user" (
                                     "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                     "water_intake_id" UUID UNIQUE,
                                     "user_id" integer UNIQUE
);

CREATE TABLE IF NOT EXISTS user_macro_distribution (
                                                       id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                                       user_id UUID NOT NULL,
                                                       age INTEGER NOT NULL,
                                                       height FLOAT NOT NULL,
                                                       weight FLOAT NOT NULL,
                                                       gender VARCHAR(10) NOT NULL,
                                                       system VARCHAR(10) NOT NULL,
                                                       activity VARCHAR(20) NOT NULL,
                                                       activity_description VARCHAR NOT NULL,
                                                       objective VARCHAR NOT NULL,
                                                       objective_description VARCHAR NOT NULL,
                                                       calories_distribution VARCHAR NOT NULL,
                                                       calories_distribution_description VARCHAR NOT NULL,
                                                       protein INTEGER NOT NULL,
                                                       fats INTEGER NOT NULL,
                                                       carbs INTEGER NOT NULL,
                                                       bmr INTEGER NOT NULL,
                                                       tdee INTEGER NOT NULL,
                                                       goal INTEGER NOT NULL,
                                                       created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


COMMENT ON COLUMN "food"."serving_size" IS 'grams';
COMMENT ON COLUMN "food"."protein" IS 'grams';
COMMENT ON COLUMN "food"."fat_total" IS 'grams';
COMMENT ON COLUMN "food"."fat_saturated" IS 'grams';
COMMENT ON COLUMN "food"."carbohydrates_total" IS 'grams';
COMMENT ON COLUMN "food"."fiber" IS 'grams';
COMMENT ON COLUMN "food"."sugar" IS 'grams';
COMMENT ON COLUMN "food"."sodium" IS 'miligrams';
COMMENT ON COLUMN "food"."potassium" IS 'miligrams';
COMMENT ON COLUMN "food"."cholesterol" IS 'miligrams';
