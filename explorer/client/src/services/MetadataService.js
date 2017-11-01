export default class MetadataService {

  // empty metadata
  static empty = {
    System: {
      'METADATA-RESOURCE': {
        Resource: [],
      },
      SystemDescription: 'Loading metadata...',
      SystemID: 'Loading...',
    },
  };

  // metadata params
  static params = {
    extraction: null,  // the extraction type to use (COMPACT-DECODED, COMPACT-INCREMENTAL, STANDARD-XML)
  }

  static get(conn, args) {
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
          ...args,
          connection: conn,
        }],
      }),
    });
  }
}
