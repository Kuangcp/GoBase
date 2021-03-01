let basePath = '/api/v1.0/'
let curFileKey = 'cur-file'
let totalItem = []

function get(url, handle) {
    let httpRequest = new XMLHttpRequest();
    httpRequest.open('GET', url, true);
    httpRequest.send();
    /**
     * 获取数据后的处理程序
     */
    httpRequest.onreadystatechange = function () {
        if (httpRequest.readyState === 4 && httpRequest.status === 200) {
            handle(httpRequest)
        }
    };
}
