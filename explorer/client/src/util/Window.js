import { createHistory } from 'history';

export const pullParamState = () => {
  const obj = {};
  const location = window.location.href;
  const queryString = location.split('?')[1];
  if (queryString) {
    const searchParams = queryString.split('&');
    searchParams.forEach(param => {
      const hasEquals = param.indexOf('=') > -1;
      const key = hasEquals ? param.split('=')[0] : param;
      const val = hasEquals ? decodeURIComponent(param.split('=')[1]) : 1;
      obj[key] = val;
    });
  }
  return obj;
};

const queryParams = (source) => {
  const array = [];
  Object.keys(source).forEach(key =>
    array.push(`${encodeURIComponent(key)}=${encodeURIComponent(source[key])}`)
  );
  return array.join('&');
};

export const pushParamState = (obj) => {
  const objClone = { ...obj };
  Object.keys(objClone).forEach(key => {
    if (objClone[key] === null || objClone[key] === undefined) {
      delete objClone[key];
    }
  });
  const history = createHistory();
  const location = history.getCurrentLocation();
  let path = location.hash || location.pathname;
  if (path.includes('?')) {
    path = path.split('?')[0];
  }
  history.push({
    search: `?${queryParams(objClone)}`,
    pathname: path,
  });
};
