USE tmanfredo;

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS customer;
DROP TABLE IF EXISTS product;

CREATE TABLE product ( 
    id SERIAL PRIMARY KEY, 
    product_name varchar(255), 
    image_name varchar(255), 
    price decimal(6,2), 
    in_stock int 
);

CREATE TABLE customer ( 
    id SERIAL PRIMARY KEY, 
    first_name varchar(255), 
    last_name varchar(255), 
    email varchar(255) 
);

CREATE TABLE orders ( 
    id SERIAL PRIMARY KEY, 
    product_id int REFERENCES product(id),
    customer_id int REFERENCES customer(id), 
    quantity int CHECK (quantity >= 0), 
    price decimal(6,2), 
    tax decimal(6,2), 
    donation decimal(6,2), 
    timestamp bigint 
);

INSERT INTO customer (first_name, last_name, email)
    VALUES ("Mickey", "Mouse", "mmouse@mines.edu"),
           ("Baiza", "Mand", "mand@mines.edu");
           

INSERT INTO product (product_name, image_name, price, in_stock)
VALUES ("Super Mario Odyssey", "assets/images/odyssey.png", 39.99, 4),
       ("Undertale", "assets/images/undertale.png", 19.99, 10),
       ("Mario Kart 8 Deluxe", "assets/images/kart.png", 39.99, 0);

/*INSERT INTO orders (product_id, customer_id, quantity, ))*/
