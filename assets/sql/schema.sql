USE tmanfredo;

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS customer;
DROP TABLE IF EXISTS product;

CREATE TABLE product ( 
    id int  NOT NULL AUTO_INCREMENT PRIMARY KEY, 
    product_name varchar(255), 
    image_name varchar(255), 
    price decimal(4,2), 
    in_stock int 
);

CREATE TABLE customer ( 
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY, 
    first_name varchar(255), 
    last_name varchar(255), 
    email varchar(255) 
);

CREATE TABLE orders ( 
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY, 
    product_id int REFERENCES product(id),
    customer_id int REFERENCES customer(id), 
    quantity int, 
    price decimal(2,2), 
    tax decimal(2,2), 
    donation decimal(2,2), 
    timestamp bigint 
);

INSERT INTO customer (first_name, last_name, email)
    VALUES ("Mickey", "Mouse", "mmouse@mines.edu"),
           ("Baiza", "Mand", "mand@mines.edu");

INSERT INTO product (product_name, image_name, price, in_stock)
VALUES ("Super Mario Odyssey", "assets/images/odyssey.png", 39.99, 4),
       ("Undertale", "assets/images/undertale.png", 19.99, 10),
       ("Mario Kart 8 Deluxe", "assets/images/kart.png", 39.99, 0);