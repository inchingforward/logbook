create table logbook_user (
    id bigserial not null,
    username text not null,
    password text not null,
    display_name text not null,
    active boolean default true not null,
    created_at timestamp with time zone default now() not null,
    last_login_at timestamp with time zone
);

create table logbook_entry (
    id bigserial not null,
    title text not null,
    url text,
    notes text,
    private boolean default true not null,
    user_id bigint not null references logbook_user(id),
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null,
    tags text
);
