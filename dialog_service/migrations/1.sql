CREATE TABLE dialogs (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW(),

    user_id_1 BIGSERIAL NOT NULL,
    dialog_name_1 VARCHAR(255) NOT NULL,
    unread_messages_count_1 BIGSERIAL NOT NULL DEFAULT 0,

    user_id_2 BIGSERIAL NOT NULL,
    dialog_name_2 VARCHAR(255) NOT NULL,
    unread_messages_count_2 BIGSERIAL NOT NULL DEFAULT 0
);


CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    dialog_id BIGSERIAL NOT NULL,
    sender_id BIGSERIAL NOT NULL,
    text VARCHAR(1000) NOT NULL
);
