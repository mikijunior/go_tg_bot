CREATE TABLE users (
    id bigserial not null primary key,
    username varchar not null unique,
    chat_id varchar not null
)