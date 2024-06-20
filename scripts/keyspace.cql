--  Created: 2023-Nov-13
-- Modified: 2024-Jun-17

-- Create a keyspace;

CREATE KEYSPACE IF NOT EXISTS polls WITH REPLICATION = {
  'class' : 'SimpleStrategy',
  'replication_factor' : '1' };


-- CREATE ROLE webapp WITH PASSWORD = 'password_a' AND LOGIN = true;


-- GRANT ALL PERMISSIONS ON KEYSPACE polls TO webapp;
