create table challenges
(
    challenge   varchar
        constraint challenges_pk
            primary key,
    valid_till  integer,
    difficulty  integer,
    max_allowed integer,
    used        integer,
    hash_algo   varchar,
    hash        varchar
);