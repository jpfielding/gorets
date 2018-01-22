const { assert } = require('chai');
const pageHandler = require('../support/support.js');

/* eslint-env node, mocha */
/* global browser, $ */

describe('goRets Explorer', () => {
  describe('Testing Connections Drop Downs', () => {
    beforeEach('Reset', () => {
      browser.url('/');
    });
    it('Should find endpoint', () => {
      assert.equal(browser.getTitle(), 'goRETS Explorer');
    });
    it('Source URL should be provided', () => {
      assert.isAbove(pageHandler.checkConfigMenu(), 0);
    });
    it('Populates connection list of connection from Config', () => {
      const init = pageHandler.checkConnectionMenu();
      pageHandler.selectConfig(0);
      assert.isAbove(pageHandler.checkConnectionMenu(), init);
    });
    it('Should notify user of Config name conflicts', () => {
      pageHandler.selectConfig(0);
      const init = pageHandler.checkConnectionMenu();
      pageHandler.selectConfig(0);
      assert.equal(pageHandler.checkConnectionMenu(), init);
      assert.include($('.error-out').getText(), 'already in use');
    });
    it('Should be able to add a new Succesful Config url', () => {
      const init = pageHandler.checkConnectionMenu();
      pageHandler.addConfig('TestNew', 'VersionB');
      assert.isTrue(pageHandler.hasConfig('TestNew'), 'The new source TestNew was added to source URL dropdown');
      assert.isAbove(pageHandler.checkConnectionMenu(), init);
    });
    it('Should not be able to add a new source url with duplicate name', () => {
      const init = pageHandler.checkConnectionMenu();
      pageHandler.addConfig('TestNew', 'VersionB');
      pageHandler.addConfig('TestNew', 'Extra');
      assert.include($('.error-out').getText(), 'already in use');
      assert.equal(pageHandler.checkConnectionMenu(), init);
    });
    it('Should open but fail to load data for bad source', () => {
      pageHandler.addManualSource('badsource');
      assert.isTrue($('#badsource-tab').isExisting(), 'The source tab badsource-tab exists');
      $('#badsource-tab').click();
      assert.include($('#badsource-error').getText(), 'Unknown Config');
      assert.isTrue($('#badsource-loading').isExisting(), 'The loading annimation sould still exist');
    });
    it('Should remove tab on close', () => {
      pageHandler.addManualSource('badsource');
      assert.isTrue($('#badsource-tab').isExisting(), 'The existance of the badsource tab');
      browser.moveToElement('#badsource-tab');
      $('#badsource-tab-close').click();
      assert.isFalse($('#badsource-tab').isExisting(), 'The existance of the badsource tab');
    });
    it('Should not open duplicate tabs', () => {
      pageHandler.selectConfigConnection(0, 0);
      pageHandler.selectConnection(0);
      assert.include($('.error-out').getText(), 'Failed to create tab, Does that tab already exist?');
    });
  });
});
