CREATE TABLE IF NOT EXISTS Carts(
    ID serial PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS Items(
  ID serial PRIMARY KEY,
  cartID serial,
  product_name varchar(255),
  quantity integer,
  CONSTRAINT fk_cart FOREIGN KEY(cartID) REFERENCES Carts(ID)
);