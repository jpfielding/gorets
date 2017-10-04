export default class ObjectsService {

  // parameters for the objects request
  static params = {
    id: null, // source
    resource: null,
    ids: null,  // <id-1>:*,<id-2>:1,<id-2>:2
    type: null, // Photo | HiRes, depends on the server metadata
    location: null, // the extraction method requested, 0 = inline, 1 = url
  };


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
