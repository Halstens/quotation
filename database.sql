CREATE TABLE quotes (
    id     SERIAL PRIMARY KEY,
    author VARCHAR(100) NOT NULL,
    text   TEXT NOT NULL
);

INSERT INTO quotes (author, text) VALUES ('Confucius', 'poel, popil, i snova poel');
INSERT INTO quotes (author, text) VALUES ('Kto to tam', 'pospal, vstal, poel i snova pospal');