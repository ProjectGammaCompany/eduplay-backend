CREATE TABLE files (
    fileId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    fileKey text NOT NULL DEFAULT '', 
    fileName text NOT NULL DEFAULT '',
    count int NOT NULL DEFAULT 0
);