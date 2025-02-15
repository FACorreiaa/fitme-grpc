CREATE TABLE conversations (
                             id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                             created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE conversation_participants (
                                         conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
                                         user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                         PRIMARY KEY (conversation_id, user_id)
);

CREATE TABLE messages (
                        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                        conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
                        sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                        content TEXT NOT NULL,
                        created_at TIMESTAMP DEFAULT now()
);


