CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR(255) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    company TEXT,
    position TEXT,
    phone TEXT,
    g7Company TEXT,
    existingG7Project TEXT,
    range TEXT,
    upcomingG7Project TEXT,
    description TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
