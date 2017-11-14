export default class ConfigService {

  static getConfigList(url) {
    return fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: 1,
        method: 'ConfigService.List',
        params: [{
          active: true,
        }],
      }),
      mode: 'no-cors',
    });
  }

}
