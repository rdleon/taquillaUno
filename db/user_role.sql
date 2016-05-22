CREATE TABLE user_role (
	uid int references users(uid),
	rid int references roles(rid)
);
