create table users (
	name text primary key,
	password text
);

create table session (
	id text primary key,
	user_name text references users (name),
	login_date timestamp default current_timestamp
);
