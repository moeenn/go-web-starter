-- +goose Up
-- +goose StatementBegin
create type
  user_role as enum('ADMIN', 'CLIENT', 'EMPLOYEE');

create table
  "user" (
    id uuid not null,
    email varchar(255) not null,
    role user_role not null,
    password text not null,
    name varchar(255),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp,
    primary key (id),
    constraint email_unique unique (email)
  );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table "user";
drop type user_role;
-- +goose StatementEnd
