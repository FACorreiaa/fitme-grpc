CREATE TABLE gyms (
                    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                    name VARCHAR(255) NOT NULL,
                    address TEXT,
                    phone VARCHAR(50),
                    created_at TIMESTAMP DEFAULT now(),
                    updated_at TIMESTAMP DEFAULT now()
);

-- Users who work at or manage the gym:
CREATE TABLE gym_staff (
                         id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                         gym_id UUID NOT NULL REFERENCES gyms(id) ON DELETE CASCADE,
                         user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                         role VARCHAR(255) NOT NULL DEFAULT 'STAFF',
                         created_at TIMESTAMP DEFAULT now(),
                         updated_at TIMESTAMP DEFAULT now(),
                         UNIQUE(gym_id, user_id)
);

-- Gym classes, schedules, participants:
CREATE TABLE gym_classes (
                           id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                           gym_id UUID NOT NULL REFERENCES gyms(id) ON DELETE CASCADE,
                           name VARCHAR(255),
                           description TEXT,
                           start_time TIMESTAMP,
                           end_time TIMESTAMP,
                           created_at TIMESTAMP DEFAULT now(),
                           updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE gym_class_participants (
                                      class_id UUID NOT NULL REFERENCES gym_classes(id) ON DELETE CASCADE,
                                      user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                      PRIMARY KEY (class_id, user_id)
);

CREATE TABLE trainer_clients (
                               trainer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               client_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               created_at TIMESTAMP DEFAULT now(),
                               PRIMARY KEY(trainer_id, client_id)
);

