-- Initial script to create the required tables for Logbook.

create table logbook_user (
	id            bigserial primary key,
	username      text not null,
	password      text not null,
	display_name  text not null,
	active        boolean default true,
	created_at    timestamp with time zone not null default now(),
	last_login_at timestamp with time zone
);

create table logbook_entry (
	id          bigserial primary key,
	title       text not null,
 	url         text,
 	notes       text,
 	private     boolean not null default true,
 	user_id     bigint references logbook_user(id) not null,
 	created_at  timestamp with time zone not null default now(),
 	entry_at    date not null default now(),    
 	tags        text
);

-- drop table logbook_entry 
-- drop table logbook_user
