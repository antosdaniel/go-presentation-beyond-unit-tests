begin;

create table expenses
(
    id       UUID PRIMARY KEY,
    amount   FLOAT,
    category VARCHAR(255),
    date     DATE,
    notes    TEXT
);

commit;
