CREATE TABLE user_points (
                           user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
                           total_points INT NOT NULL DEFAULT 0,
                           updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE achievements (
                            id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                            name VARCHAR(255) NOT NULL,
                            description TEXT,
                            points_required INT NOT NULL,
                            created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE user_achievements (
                                 user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                 achievement_id UUID NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
                                 earned_at TIMESTAMP DEFAULT now(),
                                 PRIMARY KEY (user_id, achievement_id)
);

