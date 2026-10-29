-- +goose Up
create table users (
    uuid uuid primary key default uuid_generate_v4(),
    login text not null unique,
    email text not null unique,
    password_hash text not null,
    notification_methods jsonb not null default '[]',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +goose Down
drop table users;
