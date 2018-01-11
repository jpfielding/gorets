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

module.exports = {
  'testp:johndoe': {
    json: () => dataA,
  },
  'zzzzp:janedoe': {
    json: () => dataB,
  },
};
