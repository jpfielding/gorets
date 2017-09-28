export default class SearchService {

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
}
