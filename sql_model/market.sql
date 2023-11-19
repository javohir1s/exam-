
--table for category
CREATE TABLE "category" (
    "id" UUID NOT NULL PRIMARY KEY,
    "title" VARCHAR(46) NOT NULL,
    "parent_id" UUID REFERENCES "category" ("id"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "image" VARCHAR(255) NOT NULL,
    "updated_at" TIMESTAMP
);

CREATE TABLE "product" (
    "id" UUID NOT NULL PRIMARY KEY,
    "product_id" VARCHAR(255) NOT NULL UNIQUE,
    "title" VARCHAR(46) NOT NULL,
    "description" VARCHAR NOT NULL, 
    "price" NUMERIC NOT NULL,
    "photo" VARCHAR(255) NOT NULL,
    "category_id" UUID NOT NULL REFERENCES "category"("id"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);


CREATE TABLE "client" (
    "id" UUID NOT NULL PRIMARY KEY,
    "first_name" VARCHAR(50) NOT NULL,
    "last_name" VARCHAR(50) NOT NULL,
    "phone" VARCHAR(20) NOT NULL,
    "photo" VARCHAR(255) NOT NULL,
    "date_of_birth" DATE NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);


-- Table for branches
CREATE TABLE "branches" (
    "id" UUID NOT NULL PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "phone" VARCHAR(20) NOT NULL,
    "photo" VARCHAR(255) NOT NULL,
    "work_start_hour" VARCHAR NOT NULL,
    "work_end_hour" VARCHAR NOT NULL,
    "address" VARCHAR(255) NOT NULL,
    "delivery_price" NUMERIC DEFAULT 10000,
    "active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);



CREATE TABLE "order" (
    "id" UUID NOT NULL PRIMARY KEY,
    "order_id"  VARCHAR(255) NOT NULL UNIQUE ,
    "client_id" uuid not null REFERENCES "client"("id"),
    "branch_id" uuid NOT NULL REFERENCES "branches"("id"), 
    "address" VARCHAR(255),
    "delivery_price" NUMERIC, 
    "total_count" INT,
    "total_price" NUMERIC,
    "status" VARCHAR(20) NOT NULL DEFAULT 'new',
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);


CREATE TABLE "order_products" (
    "order_product_id" UUID NOT NULL PRIMARY KEY,
    "order_id" UUID NOT NULL REFERENCES "order"("id"),
    "product_id" UUID NOT NULL REFERENCES "product"("id"),
    "discount_type" VARCHAR(20) ,
    "discount_amount" NUMERIC ,
    "quantity" NUMERIC NOT NULL,
    "price" NUMERIC NOT NULL,
    "sum" NUMERIC NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);
