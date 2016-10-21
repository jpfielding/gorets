export default class MetadataService {

  static get(connectionId) {
    return fetch(`${config.api}/rpc`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: 1,
        method: 'MetadataService.Get',
        params: [{
          id: connectionId,
        }],
      }),
    });
  }

  static search(params) {
    return fetch(`${config.api}/rpc`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: 1,
        method: 'SearchService.Run',
        params: [{
          ...params,
          format: 'COMPACT-DECODED',
          'query-type': 'DMQL2',
          'count-type': 1,
        }],
      }),
    });
  }

  static getObjectMetadata(params) {
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
    // {"id": 1,
    // "method": "ObjectService.Get",
    // "params": [
    // {
    //         "id": "aartx",
    //         "resource": "ActiveAgent",
    //         "type": "Photo",
    //         "objectid": "<KeyFieldValue>:*,<KeyFieldValue>:0",
    //         "location": 0
    // }
    // ]
    // }

}
