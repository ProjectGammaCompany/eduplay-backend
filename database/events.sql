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
    FOREIGN KEY (ownerId) REFERENCES users(userid) ON DELETE CASCADE
)

CREATE TABLE blocks (
    blockId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    eventId uuid NOT NULL,
    name text NOT NULL DEFAULT '',
    order INTEGER NOT NULL DEFAULT 0,
    isParallel BOOLEAN DEFAULT false,
    FOREIGN KEY (eventId) REFERENCES events(eventId) ON DELETE CASCADE
)

CREATE TABLE tasks (
    taskId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    blockId uuid NOT NULL,
    name text NOT NULL DEFAULT '',
    type INTEGER NOT NULL DEFAULT 0, 
    -- 0 - info, 
    -- 1 - single choice, 
    -- 2 - multiple choice, 
    -- 3 - text, 
    -- 4 - qr-code
    -- optionIds uuid[] NOT NULL DEFAULT '{}',
    files text[] NOT NULL DEFAULT '{}',
    time INTEGER NOT NULL DEFAULT 0,
    partialPoint BOOLEAN DEFAULT false,
    FOREIGN KEY (blockId) REFERENCES blocks(blockId) ON DELETE CASCADE
)

CREATE TABLE conditions (
    conditionId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    prevBlockId uuid NOT NULL,
    nextBlockId uuid NOT NULL,
    group text[] NOT NULL DEFAULT '{}',
    min INTEGER NOT NULL DEFAULT 0,
    max INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (prevBlockId) REFERENCES blocks(blockId) ON DELETE CASCADE,
    FOREIGN KEY (nextBlockId) REFERENCES blocks(blockId) ON DELETE CASCADE
)

CREATE TABLE options (
    optionId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    taskId uuid NOT NULL,
    value text NOT NULL DEFAULT '',
    isCorrect BOOLEAN DEFAULT false,
    FOREIGN KEY (taskId) REFERENCES tasks(taskId) ON DELETE CASCADE
)