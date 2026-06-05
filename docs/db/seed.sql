USE go_order_service;

INSERT INTO products (id, name, description, price, status, created_at, updated_at)
VALUES
    (1, 'Go Backend Course', 'A practical Go backend course', 19900, 'ON_SALE', NOW(), NOW()),
    (2, 'Gin Web Framework Guide', 'A hands-on Gin framework guide', 9900, 'ON_SALE', NOW(), NOW()),
    (3, 'Redis Practice Lab', 'Redis cache and idempotency practice', 12900, 'ON_SALE', NOW(), NOW())
    ON DUPLICATE KEY UPDATE
                         name = VALUES(name),
                         description = VALUES(description),
                         price = VALUES(price),
                         status = VALUES(status),
                         updated_at = NOW();

INSERT INTO inventory (product_id, stock, locked_stock, created_at, updated_at)
VALUES
    (1, 100, 0, NOW(), NOW()),
    (2, 200, 0, NOW(), NOW()),
    (3, 150, 0, NOW(), NOW())
    ON DUPLICATE KEY UPDATE
                         stock = VALUES(stock),
                         locked_stock = VALUES(locked_stock),
                         updated_at = NOW();