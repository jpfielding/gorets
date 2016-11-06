export default class ObjectsService {

  static getObjects(params) {
    return fetch(`${config.api}/rpc`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: 1,
        method: 'ObjectService.Get',
        params: [{
          ...params,
          location: 0,
        }],
      }),
    });
  }
}
