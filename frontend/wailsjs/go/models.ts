export namespace main {
	
	export class ClientStatus {
	    id: string;
	    name: string;
	    path: string;
	    format: string;
	    detected: boolean;
	    supported: boolean;
	    enabled: boolean;
	    hasBackup: boolean;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new ClientStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.format = source["format"];
	        this.detected = source["detected"];
	        this.supported = source["supported"];
	        this.enabled = source["enabled"];
	        this.hasBackup = source["hasBackup"];
	        this.notes = source["notes"];
	    }
	}
	export class MCPServer {
	    id: string;
	    name: string;
	    type: string;
	    command?: string;
	    args?: string[];
	    url?: string;
	    env?: Record<string, string>;
	    workingDir?: string;
	    notes?: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new MCPServer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.command = source["command"];
	        this.args = source["args"];
	        this.url = source["url"];
	        this.env = source["env"];
	        this.workingDir = source["workingDir"];
	        this.notes = source["notes"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class AppState {
	    configPath: string;
	    servers: MCPServer[];
	    clients: ClientStatus[];
	
	    static createFrom(source: any = {}) {
	        return new AppState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.configPath = source["configPath"];
	        this.servers = this.convertValues(source["servers"], MCPServer);
	        this.clients = this.convertValues(source["clients"], ClientStatus);
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
	export class ApplyTargetResult {
	    clientId: string;
	    clientName: string;
	    path: string;
	    success: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ApplyTargetResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.clientId = source["clientId"];
	        this.clientName = source["clientName"];
	        this.path = source["path"];
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	}
	export class ApplyResult {
	    appliedAt: string;
	    appliedServers: number;
	    results: ApplyTargetResult[];
	
	    static createFrom(source: any = {}) {
	        return new ApplyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appliedAt = source["appliedAt"];
	        this.appliedServers = source["appliedServers"];
	        this.results = this.convertValues(source["results"], ApplyTargetResult);
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
	
	export class ClientConfigPreview {
	    clientId: string;
	    path: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new ClientConfigPreview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.clientId = source["clientId"];
	        this.path = source["path"];
	        this.content = source["content"];
	    }
	}
	
	
	export class ServerInput {
	    id: string;
	    name: string;
	    type: string;
	    command: string;
	    argsText: string;
	    url: string;
	    envText: string;
	    workingDir: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.command = source["command"];
	        this.argsText = source["argsText"];
	        this.url = source["url"];
	        this.envText = source["envText"];
	        this.workingDir = source["workingDir"];
	        this.notes = source["notes"];
	    }
	}

}

