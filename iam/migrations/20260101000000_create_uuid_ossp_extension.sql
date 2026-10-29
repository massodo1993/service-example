-- +goose Up
create extension if not exists "uuid-ossp";

-- +goose Down
drop extension if exists "uuid-ossp";
