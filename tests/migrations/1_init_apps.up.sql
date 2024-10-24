INSERT INTO apps (id, name, secret)
VALUES (10, 'test', 'test-secret')
ON CONFLICT DO NOTHING;