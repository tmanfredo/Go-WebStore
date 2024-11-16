USE tmanfredo;

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS customer;
DROP TABLE IF EXISTS product;
DROP TABLE IF EXISTS users;

CREATE TABLE product ( 
    id SERIAL PRIMARY KEY, 
    product_name varchar(255), 
    image_name varchar(255), 
    price decimal(6,2), 
    in_stock int CHECK (in_stock >= 0),
    inactive boolean 
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

CREATE TABLE users ( 
    id SERIAL PRIMARY KEY,
    first_name varchar(255),
    last_name varchar(255),
    password varchar(30),
    email varchar(255),
    role int
);

INSERT INTO customer (first_name, last_name, email)
    VALUES ("Mickey", "Mouse", "mmouse@mines.edu"),
           ("Baiza", "Mand", "mand@mines.edu");
           

INSERT INTO product (product_name, image_name, price, in_stock, inactive)
VALUES ("Super Mario Odyssey", "assets/images/odyssey.png", 39.99, 4, 0),
       ("Undertale", "assets/images/undertale.png", 19.99, 10, 0),
       ("Mario Kart 8 Deluxe", "assets/images/kart.png", 39.99, 0, 1);

INSERT INTO users (first_name, last_name, password, email, role)
VALUES ("Frodo", "Baggins", "fb", "fb@mines.edu", "1"),
       ("Harry", "Potter", "hp", "hp@mines.edu", "2");