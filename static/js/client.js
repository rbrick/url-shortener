function makeGet(path, data, callback) {
    makeRequest('GET', path, data, callback);
}

function makePost(path, data, callback) {
    makeRequest('POST', path, data, callback);
}

function makeDelete(path, data, callback) {
    makeRequest('DELETE', path, data, callback);
}

function makeRequest(method, path, data, callback) {
    var request = new XMLHttpRequest();
    request.addEventListener('load', function() {
        callback(this);
    });

    var url = path + (method.toLowerCase() === 'post' ? '' : '?' + encode(data));
    request.open(method, url);

    if (method.toLowerCase() !== 'get') {
        request.send(encode(data));
    } else {
        request.send();
    }
}

function encode(data) {
    var formData = new FormData();

    for (var k in data) {
       formData.append(k, data[k]);
    }
    return formData;
}
