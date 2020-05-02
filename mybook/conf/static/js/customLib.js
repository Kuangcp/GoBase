function numPanel() {
    const numPad = new NumKeyBoard({
        precision: 2,      //精确度
        minVal: 1,         //允许输入的最小值
        maxVal: 10000000   //允许输入的最大值
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

function recordPanel () {
    layer.tab({
        area: ['800px', '600px'],
        tab: [{
            title: '账单',
            content: $("#record_tab").html()
        }, {
            title: '分类',
            content: $("#month_tab").html()
        }]
    });

    $("#queryRecordsBtn").on('click', loadRecordTables);
    $("#queryMonthRecordsBtn").on('click', loadMonthTables);

    initDateArea();

    handleGet('/account/list', function (data) {
        if (data.Success) {
            console.log('/account/list', data);
            for (i in data.Data) {
                let account = data.Data[i];
                $('#accountTypeList').append($("<option></option>").attr("value", account.ID).text(account.Name));
            }
        } else {
            layer.msg('创建失败');
            console.log(data)
        }
    });
}