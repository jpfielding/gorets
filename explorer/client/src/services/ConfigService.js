export default class ConfigService {
  static getConfigList(url, args) {
    return fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: 1,
        method: 'ConfigService.List',
        params: [{
          ...args,
        }],
      }),
    });
  }

}
