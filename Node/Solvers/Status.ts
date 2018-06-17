export class Status {
    status: string;
    date: Date;
    location: string;

    public constructor(status: string, date: Date, location: string) {
        this.status = status;
        this.date = date;
        this.location = location;
    }
}

