-- Create the products table
CREATE TABLE IF NOT EXISTS products (
                                        product_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_name varchar(100) NOT NULL,
    description text NOT NULL,
    created_at timestamp DEFAULT current_timestamp
    );

-- Create an index for better performance
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);
CREATE INDEX IF NOT EXISTS idx_products_name ON products(product_name);

-- Insert some sample data (optional)
INSERT INTO products (product_name, description) VALUES
                                                     ('iPhone 15 Pro', 'Latest iPhone with A17 Pro chip, titanium design, and advanced camera system'),
                                                     ('MacBook Pro 14-inch', 'Powerful laptop with M3 Pro chip, perfect for professional work'),
                                                     ('AirPods Pro', 'Premium wireless earbuds with active noise cancellation')
    ON CONFLICT DO NOTHING;