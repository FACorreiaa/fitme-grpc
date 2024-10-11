CREATE TABLE IF NOT EXISTS "weight_measure" (
                                                "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                                "user_id" UUID,
                                                "weight_value" decimal(10,2),
                                                "created_at" timestamp DEFAULT (now()),
                                                "updated_at" timestamp DEFAULT null
);

CREATE TABLE IF NOT EXISTS "water_intake" (
                                              "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                              "user_id" UUID,
                                              "quantity" decimal(10,2),
                                              "created_at" timestamp DEFAULT (now()),
                                              "updated_at" timestamp DEFAULT null
);

CREATE TABLE IF NOT EXISTS "waist_line" (
                                            "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                            "user_id" UUID,
                                            "quantity" decimal(10,2),
                                            "created_at" timestamp DEFAULT (now()),
                                            "updated_at" timestamp DEFAULT null
);
