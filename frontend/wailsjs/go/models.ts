export namespace fs {
	
	export class AudioFile {
	    title: string;
	    artist: string;
	    album: string;
	    album_artist: string;
	    genre: string;
	    year: number;
	    track_num: number;
	    track_total: number;
	    disc_num: number;
	    disc_total: number;
	    has_art: boolean;
	    file_path: string;
	    ext: string;
	
	    static createFrom(source: any = {}) {
	        return new AudioFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.artist = source["artist"];
	        this.album = source["album"];
	        this.album_artist = source["album_artist"];
	        this.genre = source["genre"];
	        this.year = source["year"];
	        this.track_num = source["track_num"];
	        this.track_total = source["track_total"];
	        this.disc_num = source["disc_num"];
	        this.disc_total = source["disc_total"];
	        this.has_art = source["has_art"];
	        this.file_path = source["file_path"];
	        this.ext = source["ext"];
	    }
	}

}

