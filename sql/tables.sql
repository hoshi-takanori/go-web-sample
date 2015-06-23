create table users (
	name text primary key,
	password text,
	year int,
	year_no int,
	staff_year int,
	unique (year, year_no)
);

create table session (
	id text primary key,
	user_name text references users (name),
	login_date timestamp default current_timestamp
);
