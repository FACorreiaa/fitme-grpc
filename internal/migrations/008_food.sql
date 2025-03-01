-- DO $$
--   BEGIN
--     -- Create Objective ENUM
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'objective_enum') THEN
--       CREATE TYPE objective_enum AS ENUM ('MAINTENANCE', 'BULKING', 'CUTTING');
--     END IF;
--
--     -- Create Activity ENUM
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'activity_enum') THEN
--       CREATE TYPE activity_enum AS ENUM ('SEDENTARY', 'LIGHT', 'MODERATE', 'HEAVY', 'EXTRA_HEAVY');
--     END IF;
--
--     -- Create Gender ENUM
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'gender_enum') THEN
--       CREATE TYPE gender_enum AS ENUM ('MALE', 'FEMALE');
--     END IF;
--
--     -- Create QuantityUnit ENUM
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'quantity_unit_enum') THEN
--       CREATE TYPE quantity_unit_enum AS ENUM ('GRAM', 'KILOGRAM', 'MILLILITER', 'LITER', 'OUNCE', 'POUND', 'CUP', 'TEASPOON', 'TABLESPOON');
--     END IF;
--   END $$;

CREATE TYPE objective_enum AS ENUM ('MAINTENANCE', 'BULKING', 'CUTTING');
CREATE TYPE activity_enum AS ENUM ('SEDENTARY', 'LIGHT', 'MODERATE', 'HEAVY', 'EXTRA_HEAVY');
CREATE TYPE gender_enum AS ENUM ('MALE', 'FEMALE');
CREATE TYPE quantity_unit_enum AS ENUM ('GRAM', 'KILOGRAM', 'MILLILITER', 'LITER', 'OUNCE', 'POUND', 'CUP', 'TEASPOON', 'TABLESPOON');

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

CREATE TABLE "meal_plans" (
                            "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                            "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
                            "name" VARCHAR(255) NOT NULL,
                            "description" TEXT,
                            "notes" TEXT,
                            "rating" INTEGER DEFAULT 0,
                            "objective" objective_enum,    -- Enum for the objective field
                            "activity" activity_enum,      -- Enum for the activity field
                            "gender" gender_enum,          -- Enum for the gender field
                            "quantity_unit" quantity_unit_enum, -- Enum for the quantity unit field
                            "total_macros" JSONB DEFAULT '{}'::jsonb, -- Stores nutrients in JSON format
                            "created_at" TIMESTAMP DEFAULT NOW(),
                            "updated_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_total_macros ON "meal_plans" USING GIN ("total_macros");

--
CREATE TABLE "meals" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "meal_plan_id" UUID REFERENCES "meal_plans" (id) ON DELETE CASCADE,
    "meal_number" INTEGER NOT NULL, -- E.g., breakfast, lunch, dinner
    "meal_description" VARCHAR(255),
     -- "meal_ingredients" uuid[], -- Array of ingredient IDs
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);


CREATE TABLE "meal_ingredients" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "meal_id" UUID NOT NULL REFERENCES "meals" ("id") ON DELETE CASCADE,
    "ingredient_id" UUID NOT NULL REFERENCES "ingredients" ("id") ON DELETE CASCADE,
    "quantity" FLOAT(8) NOT NULL, -- Quantity of the ingredient in grams
    "calories" FLOAT(8) NOT NULL, -- Computed calories for the quantity
    "protein" FLOAT(8) NOT NULL, -- Computed protein for the quantity
    "carbohydrates_total" FLOAT(8) NOT NULL, -- Computed carbs for the quantity
    "fat_total" FLOAT(8) NOT NULL, -- Computed fat for the quantity
    "fat_saturated" FLOAT(8) NOT NULL, -- Computed saturated fat
    "fiber" FLOAT(8) NOT NULL, -- Computed fiber
    "sugar" FLOAT(8) NOT NULL, -- Computed sugar
    "sodium" FLOAT(8) NOT NULL, -- Computed sodium
    "potassium" FLOAT(8) NOT NULL, -- Computed potassium
    "cholesterol" FLOAT(8) NOT NULL, -- Computed cholesterol
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NULL
);

-- Meal Plan Meals Table: A many-to-many relationship table linking meal plans and meals.
CREATE TABLE "meal_plan_meals" (
                                 "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                 "meal_plan_id" UUID NOT NULL REFERENCES "meal_plans" ("id") ON DELETE CASCADE,
                                 "meal_id" UUID NOT NULL REFERENCES "meals" ("id") ON DELETE CASCADE,
                                 "meal_order" INTEGER NOT NULL,
                                 "created_at" TIMESTAMP DEFAULT NOW(),
                                 "updated_at" TIMESTAMP DEFAULT NULL,
                                 UNIQUE(meal_plan_id, meal_id)
);

ALTER TABLE "meals" ADD COLUMN "total_macros" JSONB DEFAULT NULL;
CREATE UNIQUE INDEX unique_meal ON meals (id, meal_number);
CREATE UNIQUE INDEX unique_meal_plan ON meal_plans (id, name);
CREATE UNIQUE INDEX unique_ingredient ON ingredients (id, name);

-- CREATE TABLE "meal_plan_meals" (
--                                  "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
--                                  "meal_plan_id" UUID NOT NULL REFERENCES "meal_plans" ("id") ON DELETE CASCADE,
--                                  "meal_id" UUID NOT NULL REFERENCES "meals" ("id") ON DELETE CASCADE,
--                                  "created_at" TIMESTAMP DEFAULT NOW(),
--                                  "updated_at" TIMESTAMP DEFAULT NULL
-- );


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
CREATE TABLE "favourite_exercises" (
                                     id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                     user_id    UUID NOT NULL REFERENCES users(id),
                                     exercise_id UUID NOT NULL REFERENCES exercise_list(id),
                                     created_at timestamp DEFAULT now(),
                                     UNIQUE (user_id, exercise_id)
);

CREATE TABLE "favourite_activities" (
                                      id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                      user_id    UUID NOT NULL REFERENCES users(id),
                                      activity_id UUID NOT NULL REFERENCES activity(id),
                                      created_at timestamp DEFAULT now(),
                                      UNIQUE (user_id, activity_id)
);

CREATE TABLE "favourite_meals" (
                                 id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                 user_id    UUID NOT NULL REFERENCES users(id),
                                 meal_id UUID NOT NULL REFERENCES meals(id),
                                 created_at timestamp DEFAULT now(),
                                 UNIQUE (user_id, meal_id)
);

CREATE TABLE recipes (
                       id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                       user_id UUID NOT NULL REFERENCES users(id),
                       description TEXT,
                       created_at timestamp DEFAULT now(),
                       updated_at timestamp DEFAULT now()
  -- no unique constraints that block multiple
);

CREATE TABLE recipe_ingredients (
                                  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                  recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
                                  ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
                                  quantity FLOAT(8),
                                  created_at TIMESTAMP DEFAULT now(),
                                  updated_at TIMESTAMP DEFAULT now(),
                                  UNIQUE(recipe_id, ingredient_id)
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
                                     "meal_plan_id" UUID NOT NULL
                                       REFERENCES "meal_plans"(id) ON DELETE CASCADE,
                                     --"meal_type_id" UUID NOT NULL
                                     --  REFERENCES "meal_type"(id) ON DELETE CASCADE,
                                     "created_at" timestamp DEFAULT (now()),
                                     "updated_at" timestamp DEFAULT null,
                                     UNIQUE ("meal_plan_id")
);


CREATE TABLE "meal_plan_user" (
                                "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                "meal_plan_id" UUID NOT NULL
                                  REFERENCES "meal_plans"(id) ON DELETE CASCADE,
                                "user_id" UUID NOT NULL
                                  REFERENCES "users"(id) ON DELETE CASCADE,
                                "created_at" timestamp DEFAULT (now()),
                                "updated_at" timestamp DEFAULT null,
                                UNIQUE ("meal_plan_id", "user_id")
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
ALTER TABLE user_macro_distribution
  ADD COLUMN is_current BOOLEAN NOT NULL DEFAULT FALSE;

CREATE UNIQUE INDEX idx_one_current_macro_per_user
  ON user_macro_distribution (user_id)
  WHERE is_current = TRUE;



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
