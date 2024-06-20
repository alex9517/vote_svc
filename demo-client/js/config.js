//  Created : 2024-Apr-29
// Modified : 2024-May-15

const VOTE_ID = 1;
const SERVICE_URL = "https://s7026:8443";
const RESOURCES_URL = "https://ws4/votes/1/images";
const VOTE_EXPIRE = 120; // 120 sec - for tests; real value can be like 2592000 ~ one month;
const MSG_EXPIRE = 30;
const ONE_VOTE = true; // 'false' allows this user to vote many times;
const ALLOW_RESULTS = true;

const VOTED_OK = "voted_ok.html";
const VOTED_BEFORE = "voted_before.html";
const RESULTS_PAGE = "results.html";

const MSG_CURRENT_RESULTS = "Current Results (voting ends at ";
const MSG_FINAL_RESULTS = "Final Results (voting is over)";

const KEY_VOTE = "vote0928";
const KEY_APP_MESSAGE = 'vote-svc-message';

// -END-
