export default class ConnectionService {

  static login(loginParams) {
    return fetch(`${config.api}/rpc`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: 1,
        method: 'ConnectionService.Add',
        params: [{
          connection: loginParams,
          test: true,
        }],
      }),
    });
  }

  static getConnectionList() {
    return fetch(`${config.api}/rpc`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: 1,
        method: 'ConnectionService.List',
        params: [{}],
      }),
    });
  }

}
