"use strict";
const request = require('request');
const Status_1 = require("./Status");
const BootBot = require('bootbot');
class DHLSolver {
    constructor(userChat, awb) {
        this.link = 'http://www.dhl.ro/shipmentTracking?AWB=';
        this.statuses = [];
        this.userChat = userChat;
        this.awb = awb;
    }
    checkIfStatusHasChanged() {
        return true;
    }
    getStatuses() {
        var statuses = [];
        return statuses;
    }
    solve() {
        // Make the request
        request(this.link.concat(this.awb.toString()), (err, res, body) => {
            if (err) {
                console.log(err);
                return;
            }
            // Pass the response to the parser
            this.parse(body);
        });
    }
    parse(json) {
        var obj = JSON.parse(json.toString());
        var statuses = obj.results[0].checkpoints;
        statuses.forEach((elem) => {
            this.statuses.push(new Status_1.Status(elem['description'], elem['date'], elem['location']));
        });
        console.log(this.statuses);
        this.userChat.say(this.statuses[-1].toString());
    }
}
module.exports = { DHLSolver: DHLSolver };
