create table memo
(
    id   serial constraint memo_pk primary key,
    name varchar not null,
    body varchar not null
);

alter table memo owner to postgres;
