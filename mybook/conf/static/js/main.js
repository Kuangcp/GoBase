const app = '/mybook';

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

function loadAccount() {
    handleGet('/account/list', function (data) {
        if (data.Success) {
            console.log(data);
            for (i in data.Data) {
                let accont = data.Data[i];
                $('#accountArea').append('<label> <input type="radio" name="accountId" value="'
                    + accont.ID + '" required> ' + accont.Name + ' </label><br/>');
                $('#targetAccountArea').append('<label> <input type="radio" name="targetAccountId" value="'
                    + accont.ID + '"> ' + accont.Name + ' </label><br/>');
            }
        } else {
            layer.msg('创建失败');
            console.log(data)
        }
    });
}

function loadRecordType() {
    handleGet('/category/typeList', function (data) {
        if (data.Success) {
            console.log(data);
            for (i in data.Data) {
                let typeEnum = data.Data[i];
                $('#typeArea').append('<label> <input type="radio" name="typeId" value="'
                    + typeEnum.Index + '" onclick="loadCategory()" required> ' + typeEnum.Name + ' </label>');
            }
        } else {
            layer.msg('创建失败');
            console.log(data)
        }
    });
}

function loadCategory() {
    $('#categoryArea').html('');

    let typeId = $('input:radio[name="typeId"]:checked').val();
    if (typeId == 3) {
        $("#targetAccountBlock").css("display", "block")
    } else {
        $("#targetAccountBlock").css("display", "none")
    }

    let col = 16;
    handleGet('/category/list?recordType=' + typeId, function (data) {
        if (data.Success) {
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
        } else {
            layer.msg('加载帐单分类失败');
            console.log(data)
        }
    });
}

function createRecord() {
    let typeId = $('input:radio[name="typeId"]:checked').val();
    let targetAccountId = $('input:radio[name="targetAccountId"]:checked').val();
    if (typeId == 3) {
        if (targetAccountId === undefined || targetAccountId === null) {
            layer.msg('目标账户必填');
            return
        }
    }

    $.post("/mybook/record/create", $("#recordForm").serialize(), function (data) {
        if (data.Success) {
            layer.msg('创建成功');
        } else {
            layer.msg('创建失败');
            console.log(data)
        }
    });
}

function numPanel(){
    var numPad= new NumKeyBoard({
        precision: 2,        //精确度
        minVal:1,           //允许输入的最小值
        maxVal:100000          //允许输入的最大值
    });

    //打开数字键盘弹框,参数为弹框确定按钮的回调函数，回调函数的参数是输入的值
    numPad.setNumVal($("#amount").val());

    numPad.open(function(val){
        $("#amount").val(val);
    });
}