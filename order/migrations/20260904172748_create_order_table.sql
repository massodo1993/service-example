-- +goose Up
create table orders (
    id serial primary key,
    order_uuid uuid not null,
    user_uuid uuid not null,
    parts_uuids uuid[] not null default '{}',
    total_price numeric not null,
    transaction_uuid uuid,
    payment_method smallint,
    status smallint not null
);

-- +goose Down
drop table orders;
