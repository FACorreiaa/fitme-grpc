CREATE TYPE friend_request_status AS ENUM ('PENDING', 'ACCEPTED', 'REJECTED', 'BLOCKED');

CREATE TABLE friend_requests (
                               id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                               from_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               to_user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               status friend_request_status NOT NULL DEFAULT 'PENDING',
                               created_at TIMESTAMP DEFAULT now(),
                               updated_at TIMESTAMP DEFAULT now(),
                               UNIQUE(from_user_id, to_user_id)
);

