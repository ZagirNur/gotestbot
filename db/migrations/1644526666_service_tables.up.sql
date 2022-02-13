CREATE TABLE fridge
(
    id UUID PRIMARY KEY
);

CREATE TABLE chat_fridge
(
    chat_id   BIGINT PRIMARY KEY,
    fridge_id UUID NOT NULL,
    FOREIGN KEY (fridge_id) REFERENCES fridge (id)
);


CREATE TABLE product
(
    id              UUID PRIMARY KEY,
    fridge_id       UUID,
    name            TEXT      NOT NULL,
    expiration_date TIMESTAMP NOT NULL,
    created_at      TIMESTAMP NOT NULL,
    FOREIGN KEY (fridge_id) REFERENCES fridge (id)
);



