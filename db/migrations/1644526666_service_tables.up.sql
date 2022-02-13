create table fridge
(
    id UUID PRIMARY KEY
);

create table chat_fridge
(
    chat_id   BIGINT PRIMARY KEY,
    fridge_id UUID NOT NULL,
    foreign key (fridge_id) references fridge (id)
);


create table product
(
    id              UUID PRIMARY KEY,
    fridge_id       UUID,
    name            TEXT      NOT NULL,
    expiration_date TIMESTAMP NOT NULL,
    created_at      TIMESTAMP NOT NULL,
    foreign key (fridge_id) references fridge (id)
);



