const s = require('../support/support.js');

/* eslint-env node, mocha */
/* global browser, $ */

describe('goRets Explorer', () => {
  describe('Testing Source', () => {
    beforeEach('Reset', () => {
      // Static Link for testing a specific source
      browser.url(`/#/?s=eyJpZCI6InRlc3RwOmpvaG5kb
      2UiLCJsb2dpblVSTCI6Imh0dHA6Ly93d3cuZmFrZS5yZ
      XRzL2xvZ2luIiwicmV0c1ZlcnNpb24iOiJSRVRTLzEuN
      SIsInVzZXJuYW1lIjoiam9obmRvZSIsInBhc3N3b3JkI
      joicGFzc3dvcmQiLCJ1c2VyQWdlbnQiOiJVc2VyQWdlb
      nQifQ==&i=eyJpZCI6InRlc3RwOmpvaG5kb2UifQ==`);
      browser.refresh();
    });

    // Check the source is loaded correctly
    it('Should have a Connection Config with proper credentials', () => {
      s.exists('#testpjohndoe-config', 'The connection config tab testpjohndoe-config exists');
      browser.moveToElement('#testpjohndoe-config');
      s.equalMulti(
        '#testpjohndoe-config-',
        ['id', 'loginurl', 'username', 'password', 'useragent', 'useragentpassword', 'proxy', 'version'],
        ['testp:johndoe', 'http://www.fake.rets/login', 'johndoe', 'password', 'UserAgent', '', '', 'RETS/1.5'],
      );
    });

    // Checks the source can be modified
    it('Should pull new metadata on config submit', () => {
      // Faulty First
      $('#testpjohndoe-config').click();
      $('#testpjohndoe-config-id').setValue('garbage');
      $('#testpjohndoe-config-update').click();
      s.includeT('#testpjohndoe-error', 'Unknown Config');
      s.exists('#testpjohndoe-loading', 'The loading annimation sould still exist');

      // Then Good
      $('#testpjohndoe-config').click();
      $('#testpjohndoe-config-id').setValue('zzzzp:janedoe');
      $('#testpjohndoe-config-update').click();
      s.exists('#testpjohndoe-Metadata-tab', 'New metadata loaded');
      $('#testpjohndoe-Metadata-tab').click();
      s.equalT('#testpjohndoe-Metadata-systemtitle', 'ALPHA'); // Check data loaded correctly
    });

    // Metadata Tab
    describe('Metadata', () => {
      // Existance
      it('Should have a tab', () => {
        s.exists('#testpjohndoe-Metadata-tab', 'The METADATA tab exists');
      });

      // Basic Functionality
      it('Should display when a source is selected', () => {
        $('#testpjohndoe-Metadata-tab').click();
        s.equalT('#testpjohndoe-Metadata-systemtitle', 'TEST');
        $('#testpjohndoe-Metadata-A-B').click();
        s.exists('#testpjohndoe-Metadata-body', 'The metadata display body is open');
      });
      // TODO: Find a way to test the filter and order parts of the table
    });

    // Search Tab
    describe('Search', () => {
      // Existance
      it('Should have a tab', () => {
        s.exists('#testpjohndoe-Search-tab', 'The SEARCH tab exists');
      });

      // Imports
      it('Should have imported data from metadata', () => {
        $('#testpjohndoe-Metadata-tab').click();
        $('#testpjohndoe-Metadata-A-B').click();
        $('#testpjohndoe-Search-tab').click();
        s.equal('#testpjohndoe-Search-query-resource', 'A');
        s.equal('#testpjohndoe-Search-query-class', 'B');
        s.include('#testpjohndoe-Search-query-query', 'LastModifiedDateTime=');
      });

      // Success
      it('Should be able to search and display results.', () => {
        $('#testpjohndoe-Metadata-tab').click();
        $('#testpjohndoe-Metadata-A-B').click();
        $('#testpjohndoe-Search-tab').click();
        $('#testpjohndoe-Search-query-class').setValue('B');
        $('#testpjohndoe-Search-submit').click();
        s.exists('#testpjohndoe-Search-history-current-select', 'There is a current search');
        s.exists('#testpjohndoe-Search-history-0-select', 'There is a search in ordered memory');
      });

      // Failure
      it('Should be able to search and display results.', () => {
        $('#testpjohndoe-Metadata-tab').click();
        $('#testpjohndoe-Metadata-A-B').click();
        $('#testpjohndoe-Search-tab').click();
        $('#testpjohndoe-Search-query-class').setValue('C');
        $('#testpjohndoe-Search-submit').click();
      });
    });

    // Objects Tab
    describe('Objects', () => {
      // Existance
      it('Should have a tab', () => {
        s.exists('#testpjohndoe-Objects-tab', 'The OBJECTS tab exists');
      });

      // Imports
      it('Should have imported data from metadata', () => {
        $('#testpjohndoe-Metadata-tab').click();
        $('#testpjohndoe-Metadata-A-B').click();
        $('#testpjohndoe-Objects-tab').click();
        s.equal('#testpjohndoe-Objects-query-resource', 'A');
        s.exists('#testpjohndoe-Objects-submit-Photo', 'Imported the object types');
      });

      // Success
      it('Should Provide Results on valid query', () => {
        $('#testpjohndoe-Objects-tab').click();
        s.runObjectQuery('testpjohndoe', 'A', 'Example:0');
        s.exists('#testpjohndoe-Objects-history-current-select', 'There is a current search');
        s.exists('#testpjohndoe-Objects-history-0-select', 'There is a search in ordered memory');
        s.exists('#testpjohndoe-Objects-result-0-image', 'The resulting object is displayed');
      });

      // Failure
      it('Should Display objects with rets errors', () => {
        $('#testpjohndoe-Objects-tab').click();
        s.runObjectQuery('testpjohndoe', 'A', 'Error:0');
        s.exists('#testpjohndoe-Objects-result-0-error', 'Object 1 should report an error');
        s.exists('#testpjohndoe-Objects-result-1-error', 'Object 2 should report an error');
      });
    });

    // Explore Tab
    // TODO: Remove this tab when it is no longer needed
    describe('Explore', () => {
      // Existance
      it('Should have a tab', () => {
        s.exists('#testpjohndoe-Explore-tab', 'The EXPLORE tab exists');
      });
    });

    // Wirelog Tab
    describe('Wirelog', () => {
      // Existance
      it('Should have a tab', () => {
        s.exists('#testpjohndoe-Wirelog-tab', 'The WIRELOG tab exists');
      });
    });
  });
});
