
CREATE TABLE calprice
(
    t_id int PRIMARY KEY,
    t_price FLOAT NOT NULL,
    user_select JSONB NOT NULL,
    address VARCHAR(13)
);

CREATE TABLE orders
(
    o_id INT PRIMARY KEY,
    t_id INT,
    t_price FLOAT,
    status VARCHAR(10),
    create_at TIMESTAMP,
    last_edit TIMESTAMP,
    FOREIGN KEY(t_id) REFERENCES calprice(t_id)
);

CREATE TABLE stocks
(
    s_id     INT PRIMARY KEY,
    quantity INT
);

CREATE TABLE products
(
    p_id    INT PRIMARY KEY,
    p_name  VARCHAR(30),
    p_desc  VARCHAR(100),
    p_price FLOAT,
    s_id    INT,
    FOREIGN KEY (s_id) REFERENCES stocks (s_id)
);
