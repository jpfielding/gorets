export default class LoginService {

  static login(loginParams) {
    return fetch(`${config.api}/api/login`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(loginParams),
    });
  }

}
