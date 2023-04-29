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

create table quotes
(
    q_id  serial primary key,
    quote varchar,
    by    varchar
);

INSERT INTO quotes (quote, by)
VALUES ('Believe you can and you''re halfway there.', 'Theodore Roosevelt'),
       ('The only way to do great work is to love what you do.', 'Steve Jobs'),
       ('Be the change you wish to see in the world.', 'Mahatma Gandhi'),
       ('Success is not final, failure is not fatal: it is the courage to continue that counts.', 'Winston Churchill'),
       ('The best way to predict your future is to create it.', 'Abraham Lincoln'),
       ('If you want to go fast, go alone. If you want to go far, go together.', 'African proverb'),
       ('Happiness is not something ready made. It comes from your own actions.', 'Dalai Lama'),
       ('The future belongs to those who believe in the beauty of their dreams.', 'Eleanor Roosevelt'),
       ('You miss 100% of the shots you don''t take.', 'Wayne Gretzky'),
       ('The greatest glory in living lies not in never falling, but in rising every time we fall.', 'Nelson Mandela');
