CREATE TABLE IF NOT EXISTS tenders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    kod_bidang TEXT,
    kebenaran_khas TEXT,
    hari_lawat_tapak TEXT,
    link TEXT NOT NULL,
    tarikh_iklan TEXT,
    taraf TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_notified BOOLEAN DEFAULT FALSE
);
