export namespace fs {
	
	export class AudioFile {
	    title: string;
	    file_path: string;
	    ext: string;
	
	    static createFrom(source: any = {}) {
	        return new AudioFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.file_path = source["file_path"];
	        this.ext = source["ext"];
	    }
	}

}

