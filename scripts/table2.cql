--  Created: 2024-Mar-15
-- Modified: 2024-Mar-25

CREATE TABLE IF NOT EXISTS polls.voters (
  vote_id int,
  user_id text,
  created timestamp,
  PRIMARY KEY ((vote_id), user_id)
);

-- CSV
-- timestamp format: yyyy-mm-dd'T'HH:mm:ssZ
-- (where Z is the RFC-822 4-digit time zone like +/-HHmm,
-- Hours must be 00 - 23, always expressed as two digits,
-- and minutes must be 00 - 59).

-- For example:
-- 2023-06-03T08:05:00+0300
