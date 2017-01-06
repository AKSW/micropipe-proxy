const request = require('request');
const {createService, reply} = require('./createService');

const foxUrl = process.env.FOX_URL || 'http://localhost:4444/api';

const handleMessage = (payload) => {
  const {body, config, replyTo, route, newRoute, res, reply} = payload;

  res.sendStatus(204);

  const json = {
    defaults: 0,
    foxlight: 'OFF',
    input: body.text,
    lang: 'en',
    type: 'text',
    task: 'ner',
    output: 'JSON-LD',
    nif: 0,
    returnHtml: false,
  };

  request({
    method: 'POST',
    url: foxUrl,
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(json),
  }, (err, res) => {
    if (err) {
      console.error(err);
      throw err;
    }

    if (res && res.statusCode !== 200) {
      console.error(`Error code: ${res.statusCode}, ${res.statusMessage}`);
      throw new Error(`Error code: ${res.statusCode}, ${res.statusMessage}`);
    }

    const result = JSON.parse(res.body);
    const output = JSON.parse(decodeURIComponent(result.output));
    const annotations = output['@graph'] ? output['@graph'] : [];
    console.log('got annotations:', annotations);

    // add annotations
    body.annotations = annotations;

    // send reply
    reply({data: body, route: newRoute});
  });
};

createService({
  message: handleMessage,
  connect() {
    console.log('Service connected!');
  },
});
