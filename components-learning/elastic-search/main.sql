-- CREATE  TABLE  Items (
--     item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     item_name varchar(100) not null ,
--     description jsonb not null ,
--     is_deleted bool not null ,
--     created_at timestamp default  current_timestamp,
--     update_at timestamp
-- );
--
-- INSERT INTO Items (item_name, description, is_deleted) VALUES
--                                                            ('Laptop', '{"brand": "Dell", "model": "XPS 13", "specs": {"ram": "16GB", "storage": "512GB SSD"}, "price": 1200}', false),
--                                                            ('Phone', '{"brand": "Apple", "model": "iPhone 14", "specs": {"storage": "128GB", "color": "blue"}, "price": 999}', false),
--                                                            ('Book', '{"title": "PostgreSQL Guide", "author": "John Doe", "pages": 350, "categories": ["database", "programming"]}', false);
--
-- SELECT  item_name, description->'brand' as brandjsonb  FROM items;
-- SELECT item_name, description->>'brand' as brand_text FROM Items;
--
-- SELECT item_name, description->'specs'->>'ram' as ram FROM Items;
--
-- SELECT item_name, jsonb_array_elements_text(description->'categories') as category
-- FROM Items
-- WHERE description ? 'categories';
--
-- SELECT  * FROM  items WHERE  description ? 'author';
--
-- SELECT * FROM Items WHERE description @> '{"brand": "Apple"}';
-- SELECT * FROM Items WHERE '{"brand": "Apple"}' <@ description;
--
-- SELECT  json_object_agg(Items.item_name,Items.description) FROM  Items;
-- SELECT * FROM Items
-- WHERE (description->>'price')::numeric > 1000;

CREATE  TABLE  Products (
    product_id uuid primary key default gen_random_uuid(),
    product_name varchar(100) not null,
    description text not null ,
    created_at timestamp default  current_timestamp
);