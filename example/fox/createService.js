const express = require('express');
const bodyParser = require('body-parser');
const request = require('request');

const noop = () => {};

// create reply function
const reply = ({route, replyTo, data, config}) => new Promise((resolve, reject) => {
  request({
    method: 'POST',
    url: 'http://localhost:8080',
    json: {route, replyTo, data, config},
  }, (error, response) => {
    if (!error && response.statusCode === 200) {
      resolve();
      return;
    }
    reject(error);
  });
});
// export send as a standalone function
exports.send = reply;

// export create service function
exports.createService = ({message, connect = noop, error = noop}) => {
  // init server
  const app = express();
  // add body parsing
  app.use(bodyParser.json()); // for parsing application/json
  app.use(bodyParser.urlencoded({extended: true})); // for parsing application/x-www-form-urlencoded
  // error handling inside of express
  app.use((err, req, res, next) => { // eslint-disable-line
    // send error to subject
    error(err);
    // dispatch status
    res.status(500).send(err);
  });
  // handle incoming messages
  app.post('/', (req, res) => {
    const {body, config, replyTo, newRoute, route} = req.body;
    const payload = {body, config, replyTo, route, newRoute, res, reply};
    message(payload);
  });
  // start server
  app.listen(3000, function() {
    const host = this.address().address;
    const port = this.address().port;
    connect({host, port});
  });
};
