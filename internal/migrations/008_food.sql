CREATE TABLE "ingredients" (
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
                        "updated_at" timestamp DEFAULT null,
                        "user_id" UUID DEFAULT NULL, -- Associates ingredient with a user
                        CONSTRAINT fk_user FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

--
CREATE TABLE "meals" (
                       "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                       "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
                       "meal_number" INTEGER NOT NULL, -- E.g., breakfast, lunch, dinner
                       "meal_description" VARCHAR(255),
                       "created_at" TIMESTAMP DEFAULT NOW(),
                       "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "meal_ingredients" (
                                  "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                  "meal_id" UUID NOT NULL REFERENCES "meals" ("id") ON DELETE CASCADE,
                                  "ingredient_id" UUID NOT NULL REFERENCES "ingredients" ("id") ON DELETE CASCADE,
                                  "quantity" FLOAT(8) NOT NULL, -- Quantity of the ingredient in grams
                                  "created_at" TIMESTAMP DEFAULT NOW(),
                                  "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "meal_plans" (
                            "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                            "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
                            "description" VARCHAR(255),
                            "notes" VARCHAR(255),
                            "total_calories" FLOAT(8),
                            "created_at" TIMESTAMP DEFAULT NOW(),
                            "updated_at" TIMESTAMP DEFAULT NULL,
                            "rating" INTEGER DEFAULT 10
);

-- Meal Plan Meals Table: A many-to-many relationship table linking meal plans and meals.
CREATE TABLE "meal_plan_meals" (
                                 "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                 "meal_plan_id" UUID NOT NULL REFERENCES "meal_plans" ("id") ON DELETE CASCADE,
                                 "meal_id" UUID NOT NULL REFERENCES "meals" ("id") ON DELETE CASCADE,
                                 "created_at" TIMESTAMP DEFAULT NOW(),
                                 "updated_at" TIMESTAMP DEFAULT NULL
);


-- CREATE TABLE "meal_type" (
--                              "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
--                              "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
--                              "ingredient_id" UUID UNIQUE,
--                              "meal_number" integer,
--                              "meal_description" varchar(255),
--                              "created_at" timestamp DEFAULT (now()),
--                              "updated_at" timestamp DEFAULT null
-- );
-- CREATE TABLE "meal_plan" (
--                              "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
--                              "user_id" UUID UNIQUE,
--                              "meal_type_id" UUID UNIQUE,
--                              "description" varchar(255),
--                              "notes" varchar(255),
--                              "total_calories" float(8),
--                              "created_at" timestamp DEFAULT (now()),
--                              "updated_at" timestamp DEFAULT null,
--                              "rating" integer DEFAULT 10
-- );
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
                          "ingredient_id" UUID UNIQUE,
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



CREATE TABLE IF NOT EXISTS "user_macro_distribution" (
                                                       "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                                       "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
                                                       "age" INTEGER NOT NULL,
                                                       "height" FLOAT NOT NULL,
                                                       "weight" FLOAT NOT NULL,
                                                       "gender" VARCHAR(10) NOT NULL,
                                                       "system" VARCHAR(10) NOT NULL,
                                                       "activity" VARCHAR(20) NOT NULL,
                                                       "activity_description" VARCHAR NOT NULL,
                                                       "objective" VARCHAR NOT NULL,
                                                       "objective_description" VARCHAR NOT NULL,
                                                       "calories_distribution" VARCHAR NOT NULL,
                                                       "calories_distribution_description" VARCHAR NOT NULL,
                                                       "protein" INTEGER NOT NULL,
                                                       "fats" INTEGER NOT NULL,
                                                       "carbs" INTEGER NOT NULL,
                                                       "bmr" INTEGER NOT NULL,
                                                       "tdee" INTEGER NOT NULL,
                                                       "goal" INTEGER NOT NULL,
                                                       "created_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE "food_logs" (
                           "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                           "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
                           "meal_id" UUID REFERENCES "meals" ("id") ON DELETE CASCADE,
                           "quantity" FLOAT(8) NOT NULL, -- Quantity of food logged
                           "log_date" TIMESTAMP NOT NULL, -- When the food was logged
                           "created_at" TIMESTAMP DEFAULT NOW(),
                           "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "diet_preferences" (
                                  "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                  "name" VARCHAR(255) NOT NULL,
                                  "description" VARCHAR(255) NOT NULL,
                                  "created_at" TIMESTAMP DEFAULT NOW()
);

CREATE TABLE "user_diet_preferences" (
                                       "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                       "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
                                       "diet_preference_id" UUID NOT NULL REFERENCES "diet_preferences" ("id") ON DELETE CASCADE,
                                       "created_at" TIMESTAMP DEFAULT NOW()
);

CREATE TABLE "diseases" (
                           "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                           "name" VARCHAR(255) NOT NULL,
                           "description" VARCHAR(255) NOT NULL,
                           "created_at" TIMESTAMP DEFAULT NOW()
);

CREATE TABLE "allergies" (
                               "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                               "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
                               "disease_id" UUID NOT NULL REFERENCES "diseases" ("id") ON DELETE CASCADE,
                               "created_at" TIMESTAMP DEFAULT NOW()
);

-- to do later
CREATE TABLE "user_allergies" (
                                "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
                                "allergy_id" UUID NOT NULL REFERENCES "allergies" ("id") ON DELETE CASCADE,
                                "created_at" TIMESTAMP DEFAULT NOW()
);

CREATE TABLE "ingredient_categories" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "description" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);

ALTER TABLE "ingredients" ADD COLUMN "category_id" UUID REFERENCES "ingredient_categories" ("id");

CREATE TABLE "meal_tags" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP DEFAULT NOW()
);

CREATE TABLE "meal_meal_tags" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "meal_id" UUID NOT NULL REFERENCES "meals" ("id") ON DELETE CASCADE,
    "tag_id" UUID NOT NULL REFERENCES "meal_tags" ("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP DEFAULT NOW()
);

CREATE TABLE "user_meal_history" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "meal_id" UUID NOT NULL REFERENCES "meals" ("id") ON DELETE CASCADE,
    "log_date" TIMESTAMP NOT NULL DEFAULT NOW(),
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "user_food_preferences" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "ingredient_id" UUID REFERENCES "ingredients" ("id") ON DELETE CASCADE,
    "category_id" UUID REFERENCES "ingredient_categories" ("id"),
    "preference" VARCHAR(255) NOT NULL CHECK (preference IN ('like', 'dislike', 'intolerance')),
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "meal_nutritional_goals" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "meal_id" UUID NOT NULL REFERENCES "meals" ("id") ON DELETE CASCADE,
    "calories" FLOAT(8),
    "protein" FLOAT(8),
    "fat" FLOAT(8),
    "carbs" FLOAT(8),
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "shopping_lists" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "meal_plan_id" UUID REFERENCES "meal_plans" ("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "shopping_list_items" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "shopping_list_id" UUID NOT NULL REFERENCES "shopping_lists" ("id") ON DELETE CASCADE,
    "ingredient_id" UUID NOT NULL REFERENCES "ingredients" ("id"),
    "quantity" FLOAT(8) NOT NULL,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "meal_feedback" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "meal_id" UUID NOT NULL REFERENCES "meals" ("id") ON DELETE CASCADE,
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "rating" INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    "comments" TEXT,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "meal_plan_feedback" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "meal_plan_id" UUID NOT NULL REFERENCES "meal_plans" ("id") ON DELETE CASCADE,
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "rating" INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    "comments" TEXT,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "activity_logs" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "activity_type" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "activity_time" TIMESTAMP DEFAULT NOW(),
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "meal_schedules" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "meal_id" UUID NOT NULL REFERENCES "meals" ("id") ON DELETE CASCADE,
    "scheduled_date" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE TABLE "custom_meal_plans" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "trainer_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE, -- Trainer or dietician
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "meal_plan_id" UUID NOT NULL REFERENCES "meal_plans" ("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);


COMMENT ON COLUMN "ingredients"."serving_size" IS 'grams';
COMMENT ON COLUMN "ingredients"."protein" IS 'grams';
COMMENT ON COLUMN "ingredients"."fat_total" IS 'grams';
COMMENT ON COLUMN "ingredients"."fat_saturated" IS 'grams';
COMMENT ON COLUMN "ingredients"."carbohydrates_total" IS 'grams';
COMMENT ON COLUMN "ingredients"."fiber" IS 'grams';
COMMENT ON COLUMN "ingredients"."sugar" IS 'grams';
COMMENT ON COLUMN "ingredients"."sodium" IS 'miligrams';
COMMENT ON COLUMN "ingredients"."potassium" IS 'miligrams';
COMMENT ON COLUMN "ingredients"."cholesterol" IS 'miligrams';
CREATE INDEX idx_ingredients_user_id ON "ingredients" ("user_id");
