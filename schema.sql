CREATE TABLE books(
id integer primary key autoincrement,
title text NOT NULL,
author text NOT NULL,
volume text,
edition text,
publisher text,
year text NOT NULL
);

CREATE TABLE chapters(
id integer primary key autoincrement,
book_id integer NOT NULL,
number integer NOT NULL,
name text NOT NULL,
foreign key (book_id) references books(id)
);

CREATE TABLE images(
ex_id integer NOT NULL,
file_name text NOT NULL,
sequence integer NOT NULL,
foreign key (ex_id) references exercises(id),
UNIQUE (ex_id, sequence)
);

CREATE TABLE exercises(
id integer primary key autoincrement,
number integer NOT NULL,
chapter_id integer NOT NULL,
foreign key (chapter_id) references chapters(id)
);


