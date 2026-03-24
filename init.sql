DROP TABLE IF EXISTS transfers; -- Borramos para recrear con la nueva estructura
USE transfers_db;

-- CREATE TABLE IF NOT EXISTS transfers (
--     id VARCHAR(50) PRIMARY KEY,
--     sender_id VARCHAR(50) NOT NULL,
--     receiver_id VARCHAR(50) NOT NULL,
--     currency INT NOT NULL,
--     amount DECIMAL(15, 2) NOT NULL,
--     state VARCHAR(20) NOT NULL
-- );

CREATE TABLE transfers (
    IDTransaction INT AUTO_INCREMENT PRIMARY KEY, -- La nueva llave primaria
    id VARCHAR(50) NOT NULL,                      
    sender_id VARCHAR(50) NOT NULL,
    receiver_id VARCHAR(50) NOT NULL,
    currency VARCHAR(10) NOT NULL,                
    amount DECIMAL(15, 2) NOT NULL,
    state VARCHAR(20) NOT NULL
);