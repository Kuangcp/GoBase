const app = '/api';


function tip(area, title, content) {
    layer.open({
        type: 1,
        title: title,
        area: area,
        closeBtn: 1,
        resize: false,
        shadeClose: true, //点击遮罩关闭
        content: content
    });
}

function handleGet(url, success, fail) {
    const request = $.get({
        url: app + '' + url,
    });
    request.done(success);
    if (fail === undefined) {
        request.fail(failed);
    } else {
        request.fail(fail);
    }
}

function failed(data) {
    layer.alert('请求发生异常');
    console.log(data)
}

function handleMainPageAccount(account) {
    $('#accountArea').append('<label> <input type="radio" name="accountId" value="'
        + account.ID + '" required> ' + account.Name + ' </label><br/>');
    $('#targetAccountArea').append('<label> <input type="radio" name="targetAccountId" value="'
        + account.ID + '"> ' + account.Name + ' </label><br/>');

    $('#accountTypeList').append($("<option></option>").attr("value", account.ID).text(account.Name));
}

function loadRecordType() {
    handleGet('/category/typeList', function (data) {
        if (!data.Success) {
            layer.msg('加载记录类型失败');
            console.log(data);
            return;
        }

        console.log('/category/typeList', data);
        for (i in data.Data) {
            let typeEnum = data.Data[i];
            $('#typeArea').append('<label> <input type="radio" name="typeId" value="' + typeEnum.Index + '" ' +
                'onclick="loadCategoryByRecordType()" required> ' + typeEnum.Name + ' </label>');
        }
    });
}

function loadCategoryByRecordType() {
    $('#categoryArea').html('');

    let typeId = $('input:radio[name="typeId"]:checked').val();
    if (typeId === '3') {
        $("#targetAccountBlock").css("display", "block")
    } else {
        $("#targetAccountBlock").css("display", "none")
    }

    let col = 16;
    handleGet('/category/list?recordType=' + typeId, function (data) {
        if (!data.Success) {
            layer.msg('加载帐单分类失败');
            console.log(data);
            return;
        }

        for (i in data.Data) {
            let typeEnum = data.Data[i];
            if (i % col === 0) {
                $('#categoryArea').append("<tr>");
            }
            $('#categoryArea').append('<td><label> <input type="radio" name="categoryId" value="'
                + typeEnum.ID + '" required> ' + typeEnum.Name + ' </label></td>');
            if (i % col === col - 1) {
                $('#categoryArea').append("</tr>");
            }
        }
    });
}

function createRecord() {
    let typeId = $('input:radio[name="typeId"]:checked').val();
    let targetAccountId = $('input:radio[name="targetAccountId"]:checked').val();
    if (typeId === 3) {
        if (targetAccountId === undefined || targetAccountId === null) {
            layer.msg('目标账户必填');
            return
        }
    }

    $.post(app + "/record/create", $("#recordForm").serialize(), function (data) {
        if (data.Success) {
            layer.msg('记账成功');
            $("#recordResult").append($("#recordForm").serialize()+"<br/>")
        } else {
            layer.msg('记账失败');
            console.log(data)
        }
    });
}