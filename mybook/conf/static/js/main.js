const app = '/mybook';

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

function loadAccount() {
    handleGet('/account/list', function (data) {
        if (data.Success) {
            console.log('/account/list', data);
            for (i in data.Data) {
                let accont = data.Data[i];
                $('#accountArea').append('<label> <input type="radio" name="accountId" value="'
                    + accont.ID + '" required> ' + accont.Name + ' </label><br/>');
                $('#targetAccountArea').append('<label> <input type="radio" name="targetAccountId" value="'
                    + accont.ID + '"> ' + accont.Name + ' </label><br/>');

                $('#accountTypeList').append($("<option></option>").attr("value", accont.ID).text(accont.Name));
                // .append("<option value='"+accont.ID+"'>"+accont.Name+"</option>");
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
            console.log('/category/typeList', data);
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

function loadRecordTables() {
    let start = $("#startDate").val();
    let end = $("#endDate").val();
    let typeId = $("#typeId option:selected").val();
    let accountType = $("#accountTypeList option:selected").val();

    let now = new Date();
    // 获取系统前一周的时间
    if (!start) {
        let date = new Date(now - 7 * 24 * 3600 * 1000);
        start = date.toISOString().slice(0, 10);
        $("#startDate").val(start);
    }
    if (!end) {
        end = now.toISOString().slice(0, 10);
        $("#endDate").val(end);
    }
    url = '/record/list?startDate=' + start + '&endDate=' + end + '&typeId=' + typeId + '&accountId=' + accountType;
    handleGet(url, function (data) {
        if (data.Success) {
            // console.log('/category/typeList', data);

            $("#record_table_body tbody").find('tr').each(function () {
                $(this).remove();
            });

            let total = 0;
            for (i in data.Data) {
                let record = data.Data[i];

                let line = "<tr>";
                line += '<td>' + record.ID + '</td>';
                line += '<td style="text-align: right"> ' + record.AccountName + '</td>';
                line += '<td>' + record.RecordTypeName + '</td>';
                line += '<td>' + record.CategoryName + '</td>';
                line += '<td style="text-align: right">' + record.Amount / 100.0 + ' </td>';
                line += '<td>' + record.Comment + '</td>';
                line += '<td style="width: 140px">' + record.RecordTime + '</td>';

                line += '</tr>';
                $('#record_table_body > tbody:last-child').append(line);
                total += record.Amount;
            }

            $("#total").html('￥' + total / 100.0)
        } else {
            layer.msg('加载账单失败');
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

function numPanel() {
    var numPad = new NumKeyBoard({
        precision: 2,        //精确度
        minVal: 1,           //允许输入的最小值
        maxVal: 10000000          //允许输入的最大值
    });

    //打开数字键盘弹框,参数为弹框确定按钮的回调函数，回调函数的参数是输入的值
    numPad.setNumVal($("#amount").val());

    numPad.open(function (val) {
        $("#amount").val(val);
    });
}

function initDateArea() {
    let start = $("#startDate").val();
    let end = $("#endDate").val();
    let now = new Date();
    // 获取系统前一周的时间
    if (!start) {
        let date = new Date(now - 7 * 24 * 3600 * 1000);
        start = date.toISOString().slice(0, 10);
        $("#startDate").val(start);
    }
    if (!end) {
        end = now.toISOString().slice(0, 10);
        $("#endDate").val(end);
    }
}