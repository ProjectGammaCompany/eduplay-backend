CREATE TABLE events (
    eventId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL DEFAULT '',
    description text NOT NULL DEFAULT '', 
    tags text[] NOT NULL DEFAULT '{}',
    cover VARCHAR(255) DEFAULT '',
    startDate TIMESTAMP,
    endDate TIMESTAMP, 
    private BOOLEAN DEFAULT false,
    password VARCHAR(255) DEFAULT '',
    ownerId uuid NOT NULL,
    lastEditionDate TIMESTAMP DEFAULT now(), 
    showRating BOOLEAN NOT NULL DEFAULT false,
    allowDownloading BOOLEAN NOT NULL DEFAULT false,
    groupEvent BOOLEAN NOT NULL DEFAULT false
    -- FOREIGN KEY (ownerId) REFERENCES users(userid) ON DELETE CASCADE
);

CREATE TABLE blocks (
    blockId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    eventId uuid NOT NULL,
    name text NOT NULL DEFAULT '',
    blockOrder INTEGER NOT NULL DEFAULT 0,
    isParallel BOOLEAN DEFAULT false,
    showPoints BOOLEAN DEFAULT false,
    showAnswers BOOLEAN DEFAULT false,
    partialPoints BOOLEAN DEFAULT false,
    FOREIGN KEY (eventId) REFERENCES events(eventId) ON DELETE CASCADE
);

CREATE TABLE tasks (
    taskId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    blockId uuid NOT NULL,
    name text NOT NULL DEFAULT '',
    description text NOT NULL DEFAULT '',
    type INTEGER NOT NULL DEFAULT 0, 
    -- 0 - info, 
    -- 1 - single choice, 
    -- 2 - multiple choice, 
    -- 3 - text, 
    -- 4 - qr-code
    -- optionIds uuid[] NOT NULL DEFAULT '{}',
    files text[] NOT NULL DEFAULT '{}',
    time INTEGER NOT NULL DEFAULT 0,
    points INTEGER NOT NULL DEFAULT 0,
    partialPoint BOOLEAN DEFAULT false,
    taskOrder INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (blockId) REFERENCES blocks(blockId) ON DELETE CASCADE
);

CREATE TABLE conditions (
    conditionId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    prevBlockId uuid,
    nextBlockId uuid,
    groupName text[] NOT NULL DEFAULT '{}',
    min INTEGER NOT NULL DEFAULT 0,
    max INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (prevBlockId) REFERENCES blocks(blockId) ON DELETE CASCADE,
    FOREIGN KEY (nextBlockId) REFERENCES blocks(blockId) ON DELETE CASCADE
);

CREATE TABLE options (
    optionId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    taskId uuid NOT NULL,
    value text NOT NULL DEFAULT '',
    isCorrect BOOLEAN DEFAULT false,
    FOREIGN KEY (taskId) REFERENCES tasks(taskId) ON DELETE CASCADE
);

CREATE TABLE tags (
    tagId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL DEFAULT ''
);

CREATE TABLE groups (
    groupId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    eventId uuid NOT NULL,
    login text NOT NULL DEFAULT '',
    password text NOT NULL DEFAULT '',
    FOREIGN KEY (eventId) REFERENCES events(eventId) ON DELETE CASCADE
);

CREATE TABLE ratings (
    ratingId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    userId uuid NOT NULL,
    eventId uuid NOT NULL,
    rating INTEGER NOT NULL DEFAULT 0,
    -- FOREIGN KEY (userId) REFERENCES users(userid) ON DELETE CASCADE,
    FOREIGN KEY (eventId) REFERENCES events(eventId) ON DELETE CASCADE
);

CREATE TABLE userFavorites (
    favoriteId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    userId uuid NOT NULL,
    eventId uuid NOT NULL,
    -- FOREIGN KEY (userId) REFERENCES users(userid) ON DELETE CASCADE,
    FOREIGN KEY (eventId) REFERENCES events(eventId) ON DELETE CASCADE
);

CREATE TABLE answers (
    answerId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    userId uuid NOT NULL,
    taskId uuid NOT NULL,
    optionIds uuid[] NOT NULL DEFAULT '{}',
    values text[] NOT NULL DEFAULT '{}',
    points INTEGER NOT NULL DEFAULT 0,
    -- FOREIGN KEY (userId) REFERENCES users(userid) ON DELETE CASCADE,
    FOREIGN KEY (taskId) REFERENCES tasks(taskId) ON DELETE CASCADE
);

CREATE TABLE userLinks (
    linkId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    userId uuid NOT NULL,
    eventId uuid NOT NULL,
    isParticipant BOOLEAN DEFAULT true,
    currTaskId uuid DEFAULT NULL, 
    currBlockId uuid DEFAULT NULL,
    currTaskStartTime TIMESTAMP NOT NULL DEFAULT '1970-01-01 00:00:00',
    finished BOOLEAN DEFAULT false,
    FOREIGN KEY (eventId) REFERENCES events(eventId) ON DELETE CASCADE,
    -- FOREIGN KEY (userId) REFERENCES users(userid) ON DELETE CASCADE,
    FOREIGN KEY (currTaskId) REFERENCES tasks(taskId) ON DELETE CASCADE
);

CREATE TABLE joinCodes (
    code       VARCHAR(6) PRIMARY KEY,   -- or use a serial ID + unique constraint
    eventId   uuid UNIQUE NOT NULL REFERENCES events(eventId) ON DELETE CASCADE,
    expiresAt TIMESTAMP NOT NULL
);

-- CREATE INDEX idx_join_codes_expires_at ON join_codes(expiresAt);

CREATE TABLE userGroups (
    userGroupId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    userId uuid NOT NULL,
    groupId uuid NOT NULL,
    FOREIGN KEY (groupId) REFERENCES groups(groupId) ON DELETE CASCADE
);

CREATE TABLE complaints (
    complaintId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    userId uuid NOT NULL,
    eventId uuid NOT NULL,
    reason text NOT NULL DEFAULT '',
    FOREIGN KEY (eventId) REFERENCES events(eventId) ON DELETE CASCADE
);