const dataA = {
  result: {
    configs: [
      {
        id: 'testp:johndoe',
        loginURL: 'http://www.fake.rets/login',
        retsVersion: 'RETS/1.5',
        username: 'johndoe',
        password: 'password',
        userAgent: 'UserAgent',
      },
    ],
  },
};

const dataB = {
  result: {
    configs: [
      {
        id: 'zzzzp:janedoe',
        loginURL: 'http://www.fake.rets/login',
        retsVersion: 'RETS/1.5',
        username: 'janedoe',
        password: 'password',
        userAgent: 'UserAgent',
      },
    ],
  },
};

const dataE = {
  result: {},
  error: 'Test Error',
};

module.exports = {
  getData: (config) => {
    console.log(config);
    if (config === 'VersionA') {
      return { json: () => dataA };
    }
    if (config === 'VersionB') {
      return { json: () => dataB };
    }
    return { json: () => dataE };
  },
};
