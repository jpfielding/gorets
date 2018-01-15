const dataA = JSON.parse(
  `{
    "result":{
      "wirelog":"U2FtcGxlIFdpcmVsb2c=",
      "Objects": [
        {
          "Blob":"Test",
          "ContentID":"77777",
          "ContentType":"image/jpeg",
          "ObjectID": 1
        }
      ]
    }
  }`
);

const dataB = JSON.parse(
  `{
    "result":{
      "wirelog":"U2FtcGxlIFdpcmVsb2c=",
      "Objects": [
        {
          "Blob":"Test",
          "ContentID":"111111",
          "ContentType":"image/jpeg",
          "ObjectID": 2
        }
      ]
    }
  }`
);

const dataC = JSON.parse(
  `{
    "result":{
      "wirelog":"U2FtcGxlIFdpcmVsb2c=",
      "Objects": [
        {
          "RetsError":"RetsError 1"
        },
        {
          "RetsError":"RetsError 2"
        }
      ]
    }
  }`
);

const dataD = JSON.parse(
  `{
    "result":{
      "wirelog":"U2FtcGxlIFdpcmVsb2c=",
      "Objects": [
        {
          "Blob":"Test",
          "ContentID":"111111",
          "ContentType":"image/jpeg",
          "ObjectID": 2
        },
        {
          "RetsError":"RetsError 1"
        },
        {
          "RetsError":"RetsError 2"
        },
        {
          "Blob":"Test",
          "ContentID":"44444",
          "ContentType":"image/jpeg",
          "ObjectID": 5
        }
      ]
    }
  }`
);

const dataE =
  {
    error: new Error('Invalid Request'),
  };

module.exports = {
  getData: (config, args) => {
    console.log(config);
    if (config.id === 'testp:johndoe' && args.ids === 'Example:0') {
      return { json: () => dataA };
    }
    if (config.id === 'testp:johndoe' && args.ids === 'Error:0') {
      return { json: () => dataC };
    }
    if (config.id === 'testp:johndoe' && args.ids === 'Mixed:0') {
      return { json: () => dataD };
    }
    if (config.id === 'zzzzp:janedoe' && args.ids === 'Example:0') {
      return { json: () => dataB };
    }
    return { json: () => dataE };
  },
};
