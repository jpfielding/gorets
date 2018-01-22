const { assert } = require('chai');

/* global browser, $ */

browser.addCommand('moveToElement', (selector) => {
  // Firefox actions api implimentation. Pull out because
  // Error: The requested resource could not be found, or
  // a request was received using an HTTP method that is
  // not supported by the mapped resource.
  /* if (browser.desiredCapabilities.browserName === 'firefox') {
    const location = browser.getLocation(selector);
    browser.actions(
      [{
        type: 'pointer',
        id: 'finger1',
        parameters: { pointerType: 'mouse' },
        actions: [
          { type: 'pointerMove', duration: 0, x: location.x + 2, y: location.y + 2 },
        ],
      }]);
    browser.actions();
    return;
  } */
  browser.moveToObject(selector);
});

class pageHandler {
  // Adds a new config source URL to the source URL dropdown
  static addConfig(name, url) {
    browser.click('#config-autocomplete');
    browser.moveToElement('#add-source-url');
    $('#add-source-url').click();
    $('#new-url-name').setValue(name);
    $('#new-url-value').setValue(url);
    $('#new-url-submit').click();
  }

  // Fills in the New Connection tab
  static addManualSource(id) {
    $('#newcon-id').setValue(id);
    $('#newcon-login').setValue('Test.com');
    $('#newcon-username').setValue('Test');
    $('#newcon-password').setValue('password');
    $('#newcon-useragent').setValue('Agent');
    $('#newcon-useragentpassword').setValue('agntpwd');
    $('#newcon-proxy').setValue('Test');
    $('#newcon-version').setValue('Test');
    $('#newcon-submit').click();
  }

  // Gets a count of elements in the Config Menu
  static checkConfigMenu() {
    browser.click('#config-autocomplete');
    return $('#config-menu').$$('.clickable').length;
  }

  // Gets a count of elements in the Connnection Menu
  static checkConnectionMenu() {
    browser.click('#connections-autocomplete');
    return $('#connection-menu').$$('.clickable').length;
  }

  // Checks if a Config Exists by the provided name
  static hasConfig(name) {
    browser.click('#config-autocomplete');
    return $('#config-menu').$(`#${name}`).isExisting();
  }

  // Selects/Opens a Config to populate the list of connecitons
  static selectConfig(id) {
    browser.click('#config-autocomplete');
    const elem = $('#config-menu').$$('.clickable')[id];
    // elem.moveToElement();
    elem.click();
  }

  // Selects/Opens a Connection to create a tab with its connection info
  static selectConnection(id) {
    browser.click('#connections-autocomplete');
    const elem = $('#connection-menu').$$('.clickable')[id];
    // elem.moveToElement();
    elem.click();
  }

  static selectConfigConnection(fig, nec) {
    pageHandler.selectConfig(fig);
    pageHandler.selectConnection(nec);
  }

  // Sets the Object Query data and run the query
  // Assumes you are on the objects tab of a source
  static runObjectQuery(pre, resource, ids) {
    $(`#${pre}-Objects-query-resource`).setValue(resource);
    $(`#${pre}-Objects-query-ids`).setValue(ids);
    $(`#${pre}-Objects-submit-Photo`).click();
  }


  // A set of shorthand ways of calling asserts
  static exists(id, extra) {
    assert.isTrue($(id).isExisting(), extra);
  }

  static equal(id, comp) {
    assert.equal($(id).getValue(), comp);
  }

  static equalT(id, comp) {
    assert.equal($(id).getText(), comp);
  }

  static equalMulti(pre, ids, comps) {
    ids.forEach((e, i) => {
      assert.equal($(`${pre}${e}`).getValue(), comps[i]);
    });
  }

  static include(id, inc) {
    assert.include($(id).getValue(), inc);
  }

  static includeT(id, inc) {
    assert.include($(id).getText(), inc);
  }

}

module.exports = pageHandler;
