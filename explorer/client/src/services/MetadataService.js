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
    id: null,  // the source to extract metadata for
    extraction: null,  // the extraction type to use (COMPACT-DECODED, COMPACT-INCREMENTAL, STANDARD-XML)
  }

  static get(connection) {
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
          connection: connection,
        }],
      }),
    });
  }
}
