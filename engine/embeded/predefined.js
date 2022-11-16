function get(data, path, defaults=null) {
    var paths = path.split('.');
    for (var i = 0; i < paths.length; i++) {
        if (data === null || data === undefined) return defaults;
        data = data[paths[i]];
    }
	if (data === null || data === undefined) return defaults;
    return data;
}

function __json_stringify(data) {
	try {
		return JSON.stringify(data);
	} catch (err) { return ''; }
}

function __json_parse(data) {
	try {
		return JSON.parse(data);
	} catch (err) { return {}; }
}

const safeStringify = __json_stringify;
const safeParse = __json_parse;
const println = print;
