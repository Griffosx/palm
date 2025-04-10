export namespace controllers {

	export class AttachmentResponse {
	    id: number;
	    filename: string;
	    size: number;
	    mimeType: string;

	    static createFrom(source: any = {}) {
	        return new AttachmentResponse(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.filename = source["filename"];
	        this.size = source["size"];
	        this.mimeType = source["mimeType"];
	    }
	}
	export class RecipientResponse {
	    id: number;
	    email: string;
	    name: string;
	    type: string;

	    static createFrom(source: any = {}) {
	        return new RecipientResponse(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.email = source["email"];
	        this.name = source["name"];
	        this.type = source["type"];
	    }
	}
	export class EmailResponse {
	    id: number;
	    accountId: number;
	    subject: string;
	    body: string;
	    senderName: string;
	    senderEmail: string;
	    receivedAt: string;
	    isRead: boolean;
	    importance: string;
	    recipients: RecipientResponse[];
	    attachments?: AttachmentResponse[];

	    static createFrom(source: any = {}) {
	        return new EmailResponse(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.accountId = source["accountId"];
	        this.subject = source["subject"];
	        this.body = source["body"];
	        this.senderName = source["senderName"];
	        this.senderEmail = source["senderEmail"];
	        this.receivedAt = source["receivedAt"];
	        this.isRead = source["isRead"];
	        this.importance = source["importance"];
	        this.recipients = this.convertValues(source["recipients"], RecipientResponse);
	        this.attachments = this.convertValues(source["attachments"], AttachmentResponse);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListEmailsResponse {
	    emails: EmailResponse[];
	    totalCount: number;
	    page: number;
	    pageSize: number;
	    totalPages: number;

	    static createFrom(source: any = {}) {
	        return new ListEmailsResponse(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.emails = this.convertValues(source["emails"], EmailResponse);
	        this.totalCount = source["totalCount"];
	        this.page = source["page"];
	        this.pageSize = source["pageSize"];
	        this.totalPages = source["totalPages"];
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

