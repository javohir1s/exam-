CREATE OR REPLACE FUNCTION calculate_order_totals()
RETURNS TRIGGER 
LANGUAGE plpgsql 
AS $$
BEGIN
    UPDATE "order"
    SET
        "total_count" = COALESCE((SELECT SUM(quantity) FROM "order_products" WHERE "order_id" = NEW."order_id"), 0),
        "total_price" = COALESCE((SELECT SUM(sum) FROM "order_products" WHERE "order_id" = NEW."order_id"), 0)
    WHERE "id" = NEW."order_id";

    RETURN NEW;
END;
$$ ;

CREATE OR REPLACE TRIGGER update_totals AFTER
INSERT OR UPDATE ON "order_products"
FOR EACH ROW
EXECUTE FUNCTION calculate_order_totals();


CREATE OR REPLACE FUNCTION calculate_order_product_sum()
RETURNS TRIGGER
LANGUAGE PLPGSQL
AS $$
DECLARE
    product_price NUMERIC;
BEGIN
    SELECT "price" INTO product_price
    FROM "product"
    WHERE id = NEW.product_id;

    CASE NEW.discount_type
        WHEN 'fix' THEN
            NEW.sum := (product_price - NEW.discount_amount) * NEW.quantity;
        WHEN 'percent' THEN
            NEW.sum := product_price * (NEW.discount_amount / 100) * NEW.quantity;
        ELSE
            NEW.sum := product_price * NEW.quantity;
    END CASE;

    RETURN NEW;
END;
$$;

CREATE TRIGGER calculate_order_product_sum_trigger
BEFORE INSERT OR UPDATE ON "order_products"
FOR EACH ROW
EXECUTE FUNCTION calculate_order_product_sum();




CREATE OR REPLACE FUNCTION update_order_product_price()
RETURNS TRIGGER
LANGUAGE PLPGSQL
AS $$
DECLARE
    product_price NUMERIC;
BEGIN
    SELECT "price" INTO product_price
    FROM "product"
    WHERE id = NEW.product_id;

    
    NEW.price := product_price;

    RETURN NEW;
END;
$$;

CREATE TRIGGER trigger_update_order_product_price
BEFORE INSERT ON "order_products"
FOR EACH ROW
EXECUTE FUNCTION update_order_product_price();

