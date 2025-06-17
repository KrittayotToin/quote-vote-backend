# quote-vote-backend
Pretest Softnix 

users
------
id            INTEGER PRIMARY KEY AUTOINCREMENT
username      TEXT UNIQUE NOT NULL
password_hash TEXT NOT NULL

quotes
-------
id        INTEGER PRIMARY KEY AUTOINCREMENT
text      TEXT NOT NULL
author    TEXT
votes     INTEGER DEFAULT 0
created_by INTEGER REFERENCES users(id)

votes
------
id        INTEGER PRIMARY KEY AUTOINCREMENT
user_id   INTEGER REFERENCES users(id)
quote_id  INTEGER REFERENCES quotes(id)
created_at DATETIME DEFAULT CURRENT_TIMESTAMP

(UNIQUE (user_id))     -- user โหวตได้แค่ครั้งเดียว
