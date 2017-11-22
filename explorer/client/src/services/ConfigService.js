  // request:
  // endpoint/rpc
  // {
  // "id": 0,
  // "method": "ConfigService.Get",
  // "params": {
  //    "connection": {} // see config object
  //    "name": "MRIS", 
  //    "source": "MRIS"  
  //    }
  // }

  // response:
  //  {
  //   "result": {},
  //   "error": nil,
  //   "id": 0
  //   }
  export default class ConfigService {
  static getConfigList(url, args) {
    return fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: 1,
        method: 'ConfigService.List',
        params: [{
          ...args,
        }],
      }),
    });
  }

}
