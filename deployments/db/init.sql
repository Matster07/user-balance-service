create table accounts
(
    id           bigserial   not null
        constraint accounts_pk
            primary key,
    account_type varchar(50) not null,
    balance      numeric default 0
        constraint accounts_balance_check
            check (balance >= (0)::numeric)
);

alter table accounts
    owner to postgres;

create unique index accounts_id_uindex
    on accounts (id);

create table transactions
(
    id                bigserial          not null
        constraint transactions_pk
            primary key,
    type              varchar(50)        not null,
    amount            numeric            not null
        check (amount >= (0)::numeric),
    sender_id         bigint
        constraint transactions_accounts_id_fk_2
            references accounts,
    receiver_id       bigint
        constraint transactions_accounts_id_fk
            references accounts,
    last_time_updated date default now() not null,
    comment           varchar(255)
);

alter table transactions
    owner to postgres;

create table categories
(
    id            bigserial    not null
        constraint service_categories_pk
            primary key,
    category_name varchar(255) not null,
    account_id    bigint       not null
        constraint service_categories_accounts_id_fk
            references accounts
);

alter table categories
    owner to postgres;

create unique index service_categories_id_uindex
    on categories (id);

create table orders
(
    id              bigserial          not null
        constraint orders_pk
            primary key,
    category_id     bigint             not null
        constraint orders_service_categories_id_fk
            references categories,
    price           numeric            not null
        constraint orders_price_check
            check (price > (0)::numeric),
    user_account_id bigint             not null
        constraint orders_accounts_id_fk
            references accounts,
    status          varchar(255)       not null,
    creation_date   date default now() not null
);

alter table orders
    owner to postgres;

create unique index orders_id_uindex
    on orders (id);

insert into accounts (id, balance, account_type) VALUES (10, 0, 'CATEGORY_RESERVATION');
insert into accounts (id, balance, account_type) VALUES (11, 0, 'CATEGORY_RESERVATION');
insert into accounts (id, balance, account_type) VALUES (12, 0, 'CATEGORY_RESERVATION');
insert into categories (id, category_name, account_id) VALUES (1, 'RENT_CAR', 10);
insert into categories (id, category_name, account_id) VALUES (2, 'BOOK_FLAT', 11);
insert into categories (id, category_name, account_id) VALUES (3, 'WALK_THE DOG', 12);






