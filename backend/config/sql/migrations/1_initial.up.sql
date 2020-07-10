-- +migrate Up
CREATE TABLE articles
(
    id      uuid PRIMARY KEY,
    title   varchar(200) NOT NULL,
    url     varchar(255) NOT NULL,
    savedon timestamp    not null,
    readon  timestamp
);