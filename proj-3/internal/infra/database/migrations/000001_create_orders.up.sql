CREATE TABLE orders (
    id varchar(255) NOT NULL,
    price float NOT NULL, 
    tax float NOT NULL, 
    final_price float NOT NULL, PRIMARY KEY (id)
);

INSERT INTO orders (id, price, tax, final_price) VALUES ('1', 100, 10, 110);
INSERT INTO orders (id, price, tax, final_price) VALUES ('2', 200, 20, 220);