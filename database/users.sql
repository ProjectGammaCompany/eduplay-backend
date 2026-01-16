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
    FOREIGN KEY (userid) REFERENCES users(userid) ON DELETE CASCADE
);

CREATE TABLE userLinks (
    linkId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    userId uuid NOT NULL,
    eventId uuid NOT NULL,
    isParticipant BOOLEAN DEFAULT true,
    currTaskId uuid DEFAULT NULL, 
    currBlockId uuid DEFAULT NULL,
    FOREIGN KEY (eventId) REFERENCES events(eventId) ON DELETE CASCADE,
    FOREIGN KEY (userId) REFERENCES users(userid) ON DELETE CASCADE, 
    FOREIGN KEY (currTaskId) REFERENCES tasks(taskId) ON DELETE CASCADE
);