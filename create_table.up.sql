create table exchange_rate (
    code varchar(3) not null,
    code_in varchar(3) not null,
    name varchar(255) not null,
    high decimal(10, 2) not null,
    low decimal(10, 2) not null,
    var_bid decimal(10, 2) not null,
    pct_change decimal(10, 2) not null,
    bid decimal(10, 2) not null,
    ask decimal(10, 2) not null,
    timestamp_str varchar(255) not null,
    created_date timestamp not null default current_timestamp,
);