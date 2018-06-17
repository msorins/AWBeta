'use strict';
const BootBot = require('bootbot'); //  https://github.com/Charca/bootbot
const options = require('./options');
const {Wit, log} = require('node-wit');
const DHLSolver = require('./Solvers/DHLSolver').DHLSolver;

const bot = new BootBot({
  accessToken: options.accessToken,
  verifyToken: options.verifyToken,
  appSecret: options.appSecret,
});

const wit = new Wit({
  accessToken: options.witToken,
  logger: new log.Logger(log.DEBUG), // optional
});

// Select the operation
function operationSelector(chat, data) {
  // Get the operation
  if ('dhl' in data['entities']) {
    var solver = new DHLSolver(chat, data['entities']['dhl'][0].value);
    solver.solve();
  }
}

// Users
bot.on('message', (payload, chat) => {
  // Receive message from payload.sender.id
  wit.message(payload.message.text, JSON.stringify(payload.message.nlp))
    .then((data) => {
      operationSelector(chat, data);
    })
    .catch(console.error);

  chat.say('hello');
});


bot.start();
