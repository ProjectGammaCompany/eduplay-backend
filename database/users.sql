CREATE TABLE users (
    userId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) DEFAULT '',
    surname VARCHAR(255) DEFAULT '',
    passwordHash VARCHAR(255) NOT NULL,
    jobTitle VARCHAR(255) DEFAULT '',
    phone VARCHAR(255) DEFAULT '',
    email VARCHAR(255) DEFAULT '',
    city VARCHAR(255) DEFAULT '',
    avatar VARCHAR(255) DEFAULT ''
);

CREATE TABLE sessions (
    sessionId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    userid uuid NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    refresh_expires TIMESTAMP NOT NULL,
    active BOOLEAN NOT NULL DEFAULT false,
    FOREIGN KEY (userid) REFERENCES users(userid) ON DELETE CASCADE
);

