CREATE TABLE notifications (
    notifId uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    notifType VARCHAR(255) NOT NULL DEFAULT '',
    notifDate TIMESTAMP NOT NULL DEFAULT now(),
    userId uuid NOT NULL,
    eventId uuid NOT NULL,
    timeLeft VARCHAR(255) DEFAULT '',
    eventName text DEFAULT '',
    notStartedFavorite boolean DEFAULT false,
    isRead boolean DEFAULT false
)