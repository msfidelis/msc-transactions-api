CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    id_client UUID,
    amount INTEGER NOT NULL,
    type CHAR(1) NOT NULL CHECK (type IN ('c', 'd')),
    description TEXT NOT NULL,
    date TEXT NOT NULL
);