const request = require('request');
import { Status } from "./Status";
import { ISolver } from './ISolver';
const BootBot = require('bootbot');

class DHLSolver implements  ISolver {
    awb: String;
    userChat: Object;

    link: String = 'http://www.dhl.ro/shipmentTracking?AWB=';
    statuses: Array<Status> = [];

    constructor(userChat: Object, awb: String) {
        this.userChat = userChat;
        this.awb = awb;
    }

    checkIfStatusHasChanged(): Boolean {
        return true;
    }

    getStatuses(): Status[] {
        var statuses: Array<Status> = [];

        return statuses;
    }

    solve() : void {
        // Make the request
        request( this.link.concat(this.awb.toString()), (err: any, res: any, body: any) => {
            if (err) { console.log(err); return; }

            // Pass the response to the parser
            this.parse(body);
        });
    }

    parse(json: String) {
        var obj: any = JSON.parse(json.toString());
        var statuses: Array<Object> = obj.results[0].checkpoints;

        statuses.forEach( (elem: Object) => {
            this.statuses.push(  new Status(elem['description'], elem['date'], elem['location']))
        });

        console.log(this.statuses);
        this.userChat.say( this.statuses[-1].toString() );

    }
}``

export = { DHLSolver: DHLSolver }