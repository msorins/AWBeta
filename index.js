const BootBot = require('bootbot'); //https://github.com/Charca/bootbot
const express = require('express');
const options = require('./options');
const {Wit, log} = require('node-wit');

const bot = new BootBot({
    accessToken: options.accessToken,
    verifyToken: options.verifyToken,
    appSecret: options.appSecret
});

const wit = new Wit({
    accessToken: options.witToken,
    logger: new log.Logger(log.DEBUG) // optional
});


// Users
bot.on('message', (payload, chat) => {
    // Receive message from payload.sender.id
    // wit.message(payload.message.text, JSON.stringify(payload.message.nlp))
    //     .then((data) => {
    //         // write an answer to the user
    //         operationSelector(chat, data);
    //     })
    //     .catch(console.error);
    chat.say("hello");
});


bot.start();
