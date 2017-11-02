export default class ConfigService {

  static login(loginParams) {
    return fetch(`${config.api}/rpc`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: 1,
        method: 'ConfigService.Add',
        params: [{
          connection: loginParams,
          test: true,
        }],
      }),
    });
  }

  static getConfigList() {
    return fetch(`${config.api}/rpc`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: 1,
        method: 'ConfigService.List',
        params: [{}],
      }),
    });
  }

}
