const dataA = JSON.parse(
  `{
    "result":{
      "wirelog":"U2FtcGxlIFdpcmVsb2c=",
      "columns":[
        "AKey"
      ],
      "count":0,
      "maxRows":false,
      "rows":[
        [
          1
        ],
        [
          2
        ],
        [
          3
        ]
      ]
    }
  }`
);

const dataB = JSON.parse(
  `{
    "result":{
      "wirelog":"U2FtcGxlIFdpcmVsb2c=",
      "columns":[
        "AKey"
      ],
      "count":0,
      "maxRows":false,
      "rows":[
        [
          1
        ],
        [
          2
        ],
        [
          3
        ]
      ]
    }
  }`
);

const dataE =
  {
    error: 'Invalid Search Query',
  };

module.exports = {
  getData: (config, args) => {
    console.log(config);
    if (config.id === 'testp:johndoe' && args.class === 'B') {
      return { json: () => dataA };
    }
    if (config.id === 'zzzzp:janedoe' && args.class === 'B') {
      return { json: () => dataB };
    }
    return { json: () => dataE };
  },
};
