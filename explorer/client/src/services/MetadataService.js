export default class MetadataService {

  static get(connectionId) {
    return fetch(`${config.api}/api/metadata`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ id: connectionId }),
    });
  }

}
