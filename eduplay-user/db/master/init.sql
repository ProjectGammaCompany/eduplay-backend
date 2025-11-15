CREATE TABLE users (
                       userId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
                       name VARCHAR(255) DEFAULT '',
                       surname VARCHAR(255) DEFAULT '',
                       login VARCHAR(255) NOT NULL,
                       passwordHash VARCHAR(255) NOT NULL,
                       jobTitle VARCHAR(255) DEFAULT '',
                       phone VARCHAR(255) DEFAULT '',
                       email VARCHAR(255) DEFAULT '',
                       city VARCHAR(255) DEFAULT '',
                       shortOrganisationTitle VARCHAR(255) DEFAULT '',
                       INN VARCHAR(255) DEFAULT '',
                       position VARCHAR(255) DEFAULT '',
                       organisationType VARCHAR(255) DEFAULT '',
                       currentTarrif VARCHAR(255) DEFAULT '',
                       role VARCHAR(255) CHECK (role IN ('user', 'operator', 'moderator', 'jurist')),
                       organisation VARCHAR(255) DEFAULT ''
);

CREATE TABLE sessions (
    sessionId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    userid uuid NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    refresh_expires TIMESTAMP NOT NULL,
    FOREIGN KEY (userid) REFERENCES users(userid) ON DELETE CASCADE
);

CREATE TABLE userSubscriptions (
    subscriptionId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    userId uuid,
    subscriptionLevel INT DEFAULT 0,
    sessions INT DEFAULT 1,
    date TIMESTAMP NOT NULL,
    expiresAt TIMESTAMP NOT NULL,
    FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE
)