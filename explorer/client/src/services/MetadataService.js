export default class MetadataService {

  static get() {
    return fetch(`${config.api}/api/metadata`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
    });
  }

}
