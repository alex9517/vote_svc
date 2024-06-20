//  Created : 2024-Jan-11
// Modified : 2024-May-01

const month3char = [
  'Jan','Feb','Mar','Apr','May','Jun',
  'Jul','Aug','Sep','Oct','Nov','Dec',];

const RE_ISO_DATE = /^\(?([0-9]{1,4})\)?[-. ]?([0-9]{1,2})[-. ]?([0-9]{1,2})$/;

const getCurrDate = () => {
  let date = new Date();
  let day = date.getDate();

  if (day < 10) {
    day = String('0' + day);
  }
  return(date.getFullYear() + '-' +
      month3char[date.getMonth()] + '-' + day);
};

const getCurrTime = () => {
  let date = new Date();
  let hours = date.getHours();
  let minutes = date.getMinutes();

  if (hours < 10) {
    hours = String('0' + hours);
  }
  if (minutes < 10) {
    minutes = String('0' + minutes);
  }
  return(hours + ':' + minutes);
};

const getCharMonthDateString = (s) => {
  // Note! 's' must be ISO date strings "YYYY-MM-DD"

  // Note that new Date(null) is equiv to new Date(0)
  // which is "Thu Jan 01 1970 04:00:00 GMT+0400 (+04)"

  if (!s || s.trim().length === 0) return '?';

  try {
    let d = new Date(s);
    console.log('Converted to : ' + d);
    let day = d.getDate();

    if (day < 10) {
      day = String('0' + day);
    }
    let ds = d.getFullYear() + '-' + month3char[d.getMonth()] + '-' + day;
      // console.log("New date string : " + ds);
    return ds;
  }
  catch(err) {
    console.log('Exception in getCharMonthDateString() : ' + err);
    return '?';
  }
};

const thisIsDate = (s) => {
  // 1. 's' must be an ISO date string like "YYYY-MM-DD",
  // 2. 's' must be a valid date.

  if (!s || s.trim().length == 0) return false;

  const dateStr = s.trim();

  if (dateStr.match(RE_ISO_DATE)) {
    const d = dateStr.split('-');
    // console.log('Converted to : ' + d);
    const year = d[0];
    const month = d[1];
    const day = d[2];

    let leapYear = false;
    if ((!(year % 4) && year % 100) || !(year % 400)) leapYear = true;

    if (month < 1 || month > 12) return false;
    if (day < 1 || day > 31) return false;
    if (month == 2) {
      if (leapYear && day > 29) return false;
      if ((!leapYear) && day > 28) return false;
    } else if (month == 4 || month == 6 || month == 9 || month == 11) {
      if (day > 30) return false;
    }
    return true;
  }
  return false;
};

const getBetterDate = (s) => {
  // It takes the ISO date strings "YYYY-MM-DD"
  // and returns "YYYY-MON-DD" like "2020-Nay-22";

  // Note that new Date(null) is equiv to new Date(0)
  // which is "Thu Jan 01 1970 04:00:00 GMT+0400 (+04)"

  if (thisIsDate(s)) {
    const dateStr = s.trim();
    const d = new Date(dateStr);
    // console.log('Converted to : ' + d);

    try {
      let day = d.getDate();

      if (day < 10) {
        day = String('0' + day);
      }
      const ds = d.getFullYear() +
          '-' + month3char[d.getMonth()] + '-' + day;
      // console.log("New date string : " + ds);
      return ds;
    }
    catch(err) {
      console.log('Exception in getBetterDate() : ' + err);
    }
  }
  return '';
};

const getAge = (s1, s2) => {
  // 's1' - birth date,
  // 's2' - pass date,
  //
  // 's2' can be null or empty. If not then:
  //    1. 's1' and 's2' must be strings like "YYYY-MM-DD",
  //    2. 's1' and 's2' must be valid dates,
  //    3. 's1' must be a date before 's2'.

  // Note! Date(null) is equiv to Date(0) which
  //   is "Thu Jan 01 1970 04:00:00 GMT+0400 (+04)".

  if (!thisIsDate(s1)) return -1;
  const dateStr1 = s1.trim();
  const d1 = new Date(dateStr1);

  let d2 = new Date();

  if (thisIsDate(s2)) {
    const dateStr2 = s2.trim();
    d2 = new Date(dateStr2);
  }

  try {
    let numYears = d2.getFullYear() - d1.getFullYear();
    if (numYears <= 0) return -1;
    if (d2.getMonth() < d1.getMonth()) {
      numYears--;
    } else if ((d2.getMonth() == d1.getMonth()) &&
        (d2.getDate() < d1.getDate())) {
      numYears--;
    }
    return numYears;
  }
  catch(err) {
    console.log('Exception in getAge() : ' + err);
  }
  return -1;  // Probably date string format is bad;
};

const birthDateIsGood = (s) => {
  // 1. 's' must be an ISO date string like "YYYY-MM-DD",
  // 2. 's' must be a valid date,
  // 3. 's' must be a date in the past.

  if (thisIsDate(s)) {
    const dateStr = s.trim();
    const d1 = new Date(dateStr);
    const d2 = new Date();
    if (d2.getTime() > d1.getTime()) return true;
  }
  return false;
};

const passDateIsGood = (s1, s2) => {
  // 's1' is birthdate,
  //    1. 's1' and 's2' must be strings like "YYYY-MM-DD",
  //    2. 's1' and 's2' must be valid dates,
  //    3. 's1' must be a date before 's2'.

  if (thisIsDate(s1) && thisIsDate(s2)) {
    const dateStr1 = s1.trim();
    const dateStr2 = s2.trim();
    const d1 = new Date(dateStr1);  // Birthdate
    const d2 = new Date(dateStr2);  // Passdate
    if (d2.getTime() > d1.getTime()) return true;
  }
  return false;
};

const addSeconds = (numOfSeconds, date) => {
  // Note! It takes care of rolling over minutes, hours,
  // days, etc. if adding seconds changes their values.
  date.setSeconds(date.getSeconds() + numOfSeconds);
  return date;
};

const timeDifferenceInMinutes = (date1, date2) => {
  const millisec = date2 - date1;
  return millisec / (1000 * 60);
};

const timeDifferenceInSeconds = (date1, date2) => {
  const millisec = date2 - date1;
  return millisec / 1000;
};

const base64ToArrayBuffer = (base64str) => {
  if (base64str) {
    let binary_str = window.atob(base64str);
    const len = binary_str.length;
    let bytes = new Uint8Array(len);

    for (let i = 0; i < len; i++) {
      bytes[i] = binary_str.charCodeAt(i);
    }
    return bytes.buffer;
  }
  return '';
};

const arrayBufferToBase64 = (buff) => {
  let binary = '';
  let bytes = new Uint8Array(buff);
  const len = bytes.byteLength;

  for (let i = 0; i < len; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return window.btoa(binary);
};

const moveBack = () => {
  history.back();
};
/*
export const getErrPage = (httpStatus) => {
  if (httpStatus == 401 || httpStatus == 403 ||
      httpStatus == 404 || httpStatus == 400 ||
      httpStatus == 500 || httpStatus == 503 ||
      httpStatus == 412 || httpStatus == 204) {
    return '' + httpStatus + '.html'
  }
};
*/

const getErrPage = (st) => {
  if (st == 401 || st == 403 || st == 404 ||
      st == 400 || st == 412 || st == 204) {
    return "" + st + ".html";
  } else {
    return "500.html";
  }
};

const getCookie = (name) => {
  let matches = document.cookie.match(new RegExp(
    "(?:^|; )" + name.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') + "=([^;]*)"
  ));
  return matches ? decodeURIComponent(matches[1]) : undefined;
};

const getCookie0 = (cname) => {
  let name = cname + '=';
  let ca = document.cookie.split(';');

  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0)==' ') c = c.substring(1);

    if (c.indexOf(name) == 0) {
      return decodeURIComponent(c.substring(name.length,c.length));
    }
  }
  return '';
};

const setCookie = (name, value, attributes = {}) => {
  // hour  : t = 3600
  // day   : t = 86400
  // week  : t = 604800
  // month : t ~ 2592000
  // year  : t ~ 31536000, etc.
  // There is no 'forever', only very large numbers!
  // Example: setCookie('user', 'Mary', {secure: true, 'max-age': 3600});

  attributes = {
    path: '/',
    SameSite: 'Strict',
    secure: true,
    // add other defaults here
    ...attributes
  };

  if (attributes.expires instanceof Date) {
    attributes.expires = attributes.expires.toUTCString();
  }

  let updatedCookie = encodeURIComponent(name) + "=" + encodeURIComponent(value);

  for (let attributeKey in attributes) {
    updatedCookie += "; " + attributeKey;
    let attributeValue = attributes[attributeKey];
    if (attributeValue !== true) {
      updatedCookie += "=" + attributeValue;
    }
  }

  document.cookie = updatedCookie;
  console.log('Cookie: ' + document.cookie);
};

/*
const setCookie = (name, value, t, options = {}) => {
  if (!options.path) {
    options.path = '/';
  }

  if (!(options.SameSite == 'Strict' ||
      options.SameSite == 'Lax' ||
      options.SameSite == 'None')) {
    options.SameSite = 'Strict';
  }

  if (t > 0) {
    let date = new Date();
    date.setTime(date.getTime() + (t * 1000));
    // options.expires = date.toGMTString();
    options.expires = date.toUTCString();
  }

  console.log('options.SameSite: ' + options.SameSite);
  console.log('options.expires: ' + options.expires);
  console.log('options.path: ' + options.path);

  let updatedCookie = encodeURIComponent(name) + "=" + encodeURIComponent(value);

  for (let optionKey in options) {
    updatedCookie += "; " + optionKey;
    let optionValue = options[optionKey];
    if (optionValue !== true) {
      updatedCookie += "=" + optionValue;
    }
  }

  document.cookie = updatedCookie;
  console.log('Cookie: ' + document.cookie);
};

const setCookie1 = (name, value, options = {}) => {

  if (!(options.SameSite == 'Strict' ||
      options.SameSite == 'Lax' ||
      options.SameSite == 'None')) {
    options.SameSite = 'Strict';
  }

  if (options.expires && options.expires instanceof Date) {
    options.expires = options.expires.toUTCString();
  } else {
    let t = 3600;

    ma = 'max-age';
    if (options.ma) {
      t = options.ma;
    }

    let date = new Date();
    date.setTime(date.getTime() + (t * 1000));
    // options.expires = date.toGMTString();
    options.expires = date.toUTCString();

    console.log('t              : ' + t);
  }

  if (!options.path) {
    options.path = '/';
  }

  console.log('options.SameSite: ' + options.SameSite);
  console.log('options.expires: ' + options.expires);
  console.log('options.path: ' + options.path);

  let updatedCookie = encodeURIComponent(name) + "=" + encodeURIComponent(value);

  for (let optionKey in options) {
    updatedCookie += "; " + optionKey;
    let optionValue = options[optionKey];
    if (optionValue !== true) {
      updatedCookie += "=" + optionValue;
    }
  }

  document.cookie = updatedCookie;
};

const setCookie0 = (name,value,t) => {
  let date = new Date();
  date.setTime(date.getTime() + (t * 1000));
  let expire = '; expires=' + date.toGMTString();
  document.cookie = name + '=' + escape(value) + expire;
};
*/

const removeCookie = (name) => {
  setCookie(name, "", {
    'max-age': -1
  })
  console.log("Cookie " + name + " was removed.");
};

const deleteCookie = (name) => {
  setCookie(name, "", {
    'max-age': -1
  })
};

const checkMessageLen = (s) => {
  if (s) {
    // The following procedure allows to break a very
    // long line into pieces for better HTML rendering;

    s = s.replace(/":"/g, '": "');
    if (s.length <= 120) return s;
    return s.substr(0,117) + '...';
  }
  return '';
};

const getCurrFileName = () => {
  const currPage = window.location.pathname;
  return currPage.substring(currPage.lastIndexOf('/') + 1);
  // ... currPage.split('/').pop();
};

const randomStr = (len) => {
  let result = '';
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  // const charsLen = chars.length;
  const charsLen = 62;

  for (let i = 0; i < len; i++) {
    result += chars.charAt(Math.floor(Math.random() * charsLen));
  }
  return result;
};

const dec2hex = (dec) => {
  return dec.toString(16).padStart(2, "0")
};

const randomId = (len) => {
  let arr = new Uint8Array((len || 40) / 2)
  window.crypto.getRandomValues(arr)
  return Array.from(arr, dec2hex).join('')
};

const sleep = (ms) => {
  return new Promise(resolve => setTimeout(resolve, ms));
};
// console.log("Hello");
// sleep(2000).then(() => { console.log("World!"); });
//    or
// async function delayedGreeting() {
//   console.log("Hello");
//   await sleep(2000);
//   console.log("World!");
//   await sleep(2000);
//   console.log("Goodbye!");
// }
//
// delayedGreeting();

// -END-
