-- name: GetAllUsers :many
select *, count(*) over() as total_count from "user"
limit $1
offset $2;

-- name: GetUserByEmail :one
select * from "user"
where email = $1
and deleted_at is null
limit 1;

-- name: CreateUser :exec
insert into "user" (id, email, role, password, name)
values ($1, $2, $3, $4, $5);

-- name: SetUserDeletedStatus :exec
update "user"
set deleted_at = $2
where id = $1;
