import data from 'services/ObjectsServiceTest';

export default class ObjectsService {

  // parameters for the objects request
  static params = {
    resource: null,
    ids: null,  // <id-1>:*,<id-2>:1,<id-2>:2
    type: null, // Photo | HiRes, depends on the server metadata
    location: 0, // the extraction method requested, 0 = inline, 1 = url
  };


  static getObjects(conn, args) {
    if (config.api === 'test') {
      console.log('Test');
      return Promise.resolve(data.getData(conn, args));
    }
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
          ...args,
          connection: conn,
        }],
      }),
    });
  }
}
