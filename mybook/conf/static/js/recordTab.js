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
        if (!data.Success) {
            layer.msg('加载账单失败');
            console.log(data);
            return;
        }

        console.log('/category/typeList', data);
        $("#record_table_body tbody").find('tr').each(function () {
            $(this).remove();
        });

        appendRecordRow(data, 'record_table_body');
    });
}

function appendRecordRow(data, targetBlock) {
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
        $('#' + targetBlock + ' > tbody:last-child').append(line);
        total += record.Amount;
    }

    $("#total").html('￥' + total / 100.0)
}