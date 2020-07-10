CREATE TABLE articles
(
    id      uuid PRIMARY KEY,
    title   varchar(200) NOT NULL,
    url     varchar(255) NOT NULL,
    savedon timestamp    not null,
    readon  timestamp
);

CREATE TABLE articles_tags (
    article_id uuid PRIMARY KEY ,
    tag_id uuid PRIMARY KEY
);

CREATE TABLE tags (
    id uuid PRIMARY KEY ,
    name varchar(200) NOT NULL
);