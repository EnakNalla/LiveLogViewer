export namespace backend {
	
	export class Response {
	    success: boolean;
	    data: any;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class Settings {
	    theme: string;
	    tailThreshold: number;
	    tailLines: number;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.tailThreshold = source["tailThreshold"];
	        this.tailLines = source["tailLines"];
	    }
	}

}

