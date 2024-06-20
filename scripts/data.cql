--  Created: 2024-Mar-22
-- Modified: 2024-Jun-14

INSERT INTO polls.votes (vote_id, header, message, resources, deadline, authenticate, allowresults,
    co_id, co_name, co_alias, co_info, co_picture, co_count, co_updated)
  VALUES(1, 'The Most Hated Dictators', 'Current results (the voting ends 2024-Dec-31)', 'https://ws4/votes/1/images/',
    '2024-12-31T22:00:00+0300', false, true, 1, 'Joseph Stalin', 'Stalin',
    'Joseph Stalin | General Secretary (1924-1953) of Soviet Communist Party, aka Russian Social Democratic Labour Party (Bolsheviks)',
    'Joseph_Stalin.jpg', 240, toTimeStamp(toDate(now())));


INSERT INTO polls.votes (vote_id, header, message, resources, deadline, authenticate, allowresults,
    co_id, co_name, co_alias, co_info, co_picture, co_count, co_updated)
  VALUES(1, 'The Most Hated Dictators', 'Current results (the voting ends 2024-Dec-31)', 'https://ws4/votes/1/images/',
    '2024-12-31T22:00:00+0300', false, true, 2, 'Adolf Hitler', 'Hitler',
    'Adolf Hitler | German Chancellor (1932-1945), leader of NSDAP (the Nazi Party, officially the National Socialist German Workers Party, 1920-1945)',
    'Adolf_Hitler.jpg', 300, toTimeStamp(toDate(now())));


INSERT INTO polls.votes (vote_id, header, message, resources, deadline, authenticate, allowresults,
    co_id, co_name, co_alias, co_info, co_picture, co_count, co_updated)
  VALUES(1, 'The Most Hated Dictators', 'Current results (the voting ends 2024-Dec-31)', 'https://ws4/votes/1/images/',
    '2024-12-31T22:00:00+0300', false, true, 3, 'Kim Jong Un', 'Kim',
    'Kim Jong Un | President (?) of North Korea, General Secretary of Nort Korean Communist Party',
    'Kim_Jong_Un.jpg', 190, toTimeStamp(toDate(now())));


INSERT INTO polls.votes (vote_id, header, message, resources, deadline, authenticate, allowresults,
    co_id, co_name, co_alias, co_info, co_picture, co_count, co_updated)
  VALUES(1, 'The Most Hated Dictators', 'Current results (the voting ends 2024-Dec-31)', 'https://ws4/votes/1/images/',
    '2024-12-31T22:00:00+0300', false, true, 4, 'Augusto Pinochet', 'Pinochet',
    'Augusto Pinochet | A Chilean army officer and military dictator who ruled Chile from 1973 to 1990',
    'Augusto_Pinochet.jpg', 100, toTimeStamp(toDate(now())));
