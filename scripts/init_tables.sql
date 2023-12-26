CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS frog (
    id uuid default uuid_generate_v4() primary key,
    photoId varchar unique not null
);

CREATE TABLE IF NOT EXISTS turbo (
    id uuid default uuid_generate_v4() primary key,
    userId int not null,
    linerId int not null,
    createdAt timestamp default current_timestamp
);

CREATE TABLE IF NOT EXISTS donate (
    id uuid default uuid_generate_v4() primary key,
    username varchar(16) not null,
    sum float not null,
    createdAt timestamp default current_timestamp
);

CREATE TABLE IF NOT EXISTS anime (
    id uuid default uuid_generate_v4() primary key,
    photoId varchar unique not null
);

create table if not exists words (
    word varchar(32) not null unique,
    len int not null,
    count int not null default 1
);
create or replace function calc_word_len() returns trigger as $cwl$
    begin
        new.len = length(new.word);
        return new;
    end;
$cwl$ LANGUAGE plpgsql;
create or replace trigger set_word_len before insert or update on words
    for each row execute procedure calc_word_len();

drop table if exists quote;
create table if not exists quote (
    id uuid default uuid_generate_v4() primary key,
    userId int not null,
    userName varchar(32) not null,
    text varchar(255) not null,
    createdAt timestamp default current_timestamp
);
