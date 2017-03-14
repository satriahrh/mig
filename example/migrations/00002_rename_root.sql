-- +mig Up
-- +mig StatementBegin
UPDATE users SET username='admin' WHERE username='root';
-- +mig StatementEnd

-- +mig Down
-- +mig StatementBegin
UPDATE users SET username='root' WHERE username='admin';
-- +mig StatementEnd
