--  Created: 2024-Mar-15
-- Modified: 2024-Jun-17

-- This table stores votes' data and results;

CREATE TABLE IF NOT EXISTS polls.votes (
  vote_id int,
  co_id smallint,
  header text,
  message text,
  resources text,
  deadline timestamp,
  authenticate boolean,
  allowresults boolean,
  co_name text,
  co_alias text,
  co_info text,
  co_picture text,
  co_count bigint,
  co_updated timestamp,
  PRIMARY KEY ((vote_id), co_id)
);

-- CSV
-- timestamp format: yyyy-mm-dd'T'HH:mm:ssZ
-- (where Z is the RFC-822 4-digit time zone like +/-HHmm,
-- Hours must be 00 - 23, always expressed as two digits,
-- and minutes must be 00 - 59).

-- For example:
-- 2023-06-03T08:05:00+0300
