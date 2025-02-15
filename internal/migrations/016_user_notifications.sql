CREATE TABLE notifications (
                             id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                             user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                             type VARCHAR(50) NOT NULL,     -- e.g., "FRIEND_REQUEST", "NEW_MESSAGE", ...
                             message TEXT NOT NULL,
                             is_read BOOLEAN NOT NULL DEFAULT FALSE,
                             created_at TIMESTAMP DEFAULT now(),
                             updated_at TIMESTAMP DEFAULT now()
);


