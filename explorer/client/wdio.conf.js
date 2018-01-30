exports.config = {
  specs: [
    './test/specs/**/*.js',
  ],

  suites: {
    base: [
      './test/specs/base.js',
    ],
    source: [
      './test/specs/source.js',
    ],
  },

  services: ['selenium-standalone'],

  maxInstances: 10,

  capabilities: [
    {
      maxInstances: 5,

      browserName: 'chrome',
      chromeOptions: {
        args: ['headless', 'disable-gpu', 'window-size=1920,1080'],
      },
    },
    /* { firefox actions api not working in test/support/support.js
      maxInstances: 5,

      browserName: 'firefox',
      'moz:firefoxOptions': {
        args: ['-headless'],
      },
    }, */
  ],

  options: {
    host: 'localhost',
    port: 4444,
  },

  sync: true,
  //
  // Level of logging verbosity: silent | verbose | command | data | result | error
  logLevel: 'silent',

  coloredLogs: true,
  deprecationWarnings: false,


  bail: 2,

  screenshotPath: './errorShots/',

  baseUrl: 'http://localhost:8000',

  waitforTimeout: 10000,
  connectionRetryTimeout: 90000,
  connectionRetryCount: 3,

  framework: 'mocha',

  mochaOpts: {
    ui: 'bdd',
  },

/* seleniumArgs: {
    seleniumArgs: ['-port', '4444'],
  }, */
};
