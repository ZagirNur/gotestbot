create table product
(
    chat_id         BIGINT PRIMARY KEY,
    id              UUID      NOT NULL,
    name            TEXT      NOT NULL,
    expiration_date TIMESTAMP NOT NULL,
    created_at      TIMESTAMP NOT NULL
);



