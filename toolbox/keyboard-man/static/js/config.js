const app = '/api/v1.0';

function handleGet(url, success, fail) {
    const request = $.get({
        url: app + '' + url,
    });
    request.done(success);
    request.fail(fail);
}

function handlePost(url, data, success, fail) {
    const request = $.post({
        url: app + '' + url,
        contentType: "application/json",
        data: JSON.stringify(data)
    });
    request.done(success);
    request.fail(fail);
}

function tip(area, title, content){
    layer.open({
        type: 1,
        title: title,
        area: area,
        closeBtn:1,
        resize: false,
        shadeClose: true, //点击遮罩关闭
        content: content
      });
}