CREATE TABLE button
(
    id     UUID PRIMARY KEY,
    action VARCHAR NOT NULL,
    data   JSONB   NOT NULL
);

create table chat_info
(
    chat_id           BIGINT PRIMARY KEY,
    active_chain      VARCHAR NOT NULL,
    active_chain_step VARCHAR NOT NULL,
    chain_data        JSONB   NOT NULL
);

