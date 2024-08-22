PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;

CREATE TABLE bookId(
id integer primary key autoincrement
);

CREATE TABLE bookInfo(
bookID integer,
typeField text,
content text,
foreign key (bookID) references bookId(id)
);

CREATE TABLE exerciseData(
exID integer,
imageName text,
foreign key (exID) references exerciseId(id)
);

CREATE TABLE chapters(
chapterID integer primary key autoincrement,
bookID integer,
number integer,
name text,
foreign key (bookID) references bookId(id)
);

CREATE TABLE exerciseId(
id integer primary key autoincrement,
exNum integer,
chapterID integer,
foreign key (chapterID) references chapters(chaptersID)
);

COMMIT;
